package server

import (
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:8080"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// public routes
	r.GET("/", s.helloWorldHandler)

	// auth routes
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/login", s.loginHandler)
		authGroup.GET("/google", s.authGoogleHandler)
		authGroup.GET("/google/callback", s.authGoogleCallbackHandler)
		authGroup.GET("/logout", s.logoutHandler)
		authGroup.GET("/register", s.userCreateHandler)

		authGroup.POST("/register", s.userRegisterHandler)
	}

	// user routes
	userGroup := r.Group("/user")
	userGroup.Use(s.authCheck())
	{
		userGroup.GET("/mystories", s.userDashboardHandler)
	}

	// api routes
	publicApi := r.Group("/api")
	{
		publicApi.GET("/getstories/:username/public", s.getPublicUserStoriesHandler)
	}

	privateApi := r.Group("/api")
	privateApi.Use(s.authCheck())
	{
		privateApi.POST("/newstory", s.newStoryHandler)
		publicApi.POST("/editstory", s.editStoryHandler)

		publicApi.GET("/getstories/:username/private", s.getPublicUserStoriesHandler)
		publicApi.GET("/getstories/:username/all", s.getPublicUserStoriesHandler)
	}

	// static assets
	staticAssets := r.Group("/")
	{
		staticAssets.StaticFile("/favicon.ico", "./favicon.ico")
		staticAssets.Static("/css", "webpages/css")
		staticAssets.Static("/images", "webpages/images")
		staticAssets.Static("/scripts", "webpages/scripts")
	}

	return r
}

func execTemplate(c *gin.Context, data any, filenames ...string) error {
	t, err := template.ParseFiles(filenames...)
	if err != nil {
		return err
	}

	err = t.Execute(c.Writer, data)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}

func (s *Server) logoutHandler(c *gin.Context) {
	session, err := s.store.Get(c.Request, "sessionISWP")

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	session.Values = make(map[any]any)
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) authGoogleCallbackHandler(c *gin.Context) {
	code, stateString := c.Query("code"), c.Query("state")

	if code == "" {
		fmt.Println("Invalid callback")
		c.Writer.WriteHeader(400)
		return
	}

	session, err := s.store.Get(c.Request, "sessionISWP")

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if session.Values["state"] != stateString {
		c.Writer.WriteHeader(401)
		return
	}

	v, err := s.exchangeTokenForUser(code, c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Values["userInfo"] = v
	if err := session.Save(c.Request, c.Writer); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.getUser(&v)
	if (err != nil) && errors.Is(err, pgx.ErrNoRows) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
		return
	} else if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Values["userInfo"] = v
	if err := session.Save(c.Request, c.Writer); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) authGoogleHandler(c *gin.Context) {
	u, err := s.getUserInfoFromSession(c)

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	} else if u.getState(AuthDB) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	} else if u.getState(AuthGoogle) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
		return
	}

	oAuthReqUrl, err := s.setAuthSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, oAuthReqUrl, http.StatusTemporaryRedirect)
}

func (s *Server) helloWorldHandler(c *gin.Context) {
	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = execTemplate(c, u, "webpages/index.html"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) loginHandler(c *gin.Context) {
	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if u.getState(AuthGoogle) && u.getState(AuthDB) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	} else if u.getState(AuthGoogle) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
		return
	}

	if err = execTemplate(c, u, "webpages/login.html"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) userCreateHandler(c *gin.Context) {
	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if u.getState(AuthDB) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	} else if !u.getState(AuthGoogle) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
	}

	if err = execTemplate(c, u, "webpages/createuser.html"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) userRegisterHandler(c *gin.Context) {
	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if u.getState(AuthDB) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	} else if !u.getState(AuthGoogle) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
	}

	username := c.PostForm("username")

	if err = validateUsername(username); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		if _, err = c.Writer.WriteString(err.Error()); err != nil {
			fmt.Println(err)
		}
		return
	}

	u.Username = username

	err = s.registerUser(&u)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			fmt.Println(pgErr.Message)
			fmt.Println(pgErr.Code)
		} else {
			fmt.Println(err)
		}

		c.Writer.WriteHeader(http.StatusInternalServerError)
		if _, err = c.Writer.WriteString(err.Error()); err != nil {
			fmt.Println(err)
		}

		return
	}

	if err = s.setUserInfoToSession(c, u); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

func (s *Server) userDashboardHandler(c *gin.Context) {
	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = execTemplate(c, u, "webpages/dashboard.html")
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) getPublicUserStoriesHandler(c *gin.Context) {
	stories, err := s.getStories(c.Param("username"))
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	if _, err := c.Writer.Write([]byte(fmt.Sprintf("%v", stories))); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) newStoryHandler(c *gin.Context) {
	st := story{
		StoryID:     rand.Text()[:15],
		StoryName:   c.PostForm("story_name"),
		Description: c.PostForm("description"),
		Visibility:  -1,
	}

	if visStr := c.PostForm("visibility"); visStr != "" {
		st.Visibility, _ = strconv.Atoi(visStr[:1])
	}

	storyID, err := s.newStory(c, &st)

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		if _, err := c.Writer.Write([]byte(err.Error())); err != nil {
			fmt.Println(err)
		}
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	if _, err := c.Writer.Write([]byte(storyID)); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) editStoryHandler(c *gin.Context) {
	st := story{
		StoryID:     c.PostForm("story_id"),
		StoryName:   c.PostForm("story_name"),
		Description: c.PostForm("description"),
		Visibility:  -1,
	}

	if visStr := c.PostForm("visibility"); visStr != "" {
		st.Visibility, _ = strconv.Atoi(visStr[:1])
	}

	if err := s.editStory(c, &st); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString(err.Error())
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.WriteString("Story edited!")
}

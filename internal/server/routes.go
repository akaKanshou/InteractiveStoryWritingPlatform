package server

import (
	"crypto/rand"
	"errors"
	"fmt"
	"forgeverse/internal/auth"
	"forgeverse/internal/db"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//TODO: change path params to postform params in getstories/user/

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
		publicApi.GET("/getstories/:username", s.getUserStoriesHandler)
		publicApi.GET("/getchapter/:chapterID", s.getChapterHandler)

		publicApi.POST("/getstorychapters", s.getStoryChaptersHandler)

		publicApi.POST("/getedges", s.getEdgesHandler)
	}

	privateApi := r.Group("/api")
	privateApi.Use(s.authCheck())
	{
		privateApi.POST("/newstory", s.newStoryHandler)
		privateApi.POST("/editstory", s.editStoryHandler)

		privateApi.POST("/newchapter", s.newChapterHandler)
		privateApi.POST("/editchapter", s.editChapterHandler)

		privateApi.POST("/addedge", s.addEdgeHandler)
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
	err := auth.Logout(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) authGoogleCallbackHandler(c *gin.Context) {
	var err fverrors.Error

	err = auth.GoogleOAuthCallback(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	user, err := auth.GetUserFromSession(c)

	if err != nil && errors.Is(err, fverrors.NoLoginErr) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		return
	}

	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	if auth.CheckAuth(user, auth.Db) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	user, err = db.GetUserByID(user.UserID)
	if (err != nil) && (errors.Is(err, fverrors.UserNotFoundError)) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
		return
	} else if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	if err = auth.SaveUserToSession(c, user); err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) authGoogleHandler(c *gin.Context) {
	var err fverrors.Error

	user, err := auth.GetUserFromSession(c)
	if err == nil {
		if auth.CheckAuth(user, auth.Db) {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else if auth.CheckAuth(user, auth.Google) {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
			return
		}
	}

	if !errors.Is(err, fverrors.NoLoginErr) {
		fverrors.SendErrorResponse(c, err)
		return
	}

	oAuthReqUrl, err := auth.InitGoogleOAuth(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	http.Redirect(c.Writer, c.Request, oAuthReqUrl, http.StatusTemporaryRedirect)
}

func (s *Server) helloWorldHandler(c *gin.Context) {
	u, err := auth.GetUserFromSession(c)
	if err != nil && !errors.Is(err, fverrors.NoLoginErr) {
		fverrors.SendErrorResponse(c, err)
		return
	} else if errors.Is(err, fverrors.NoLoginErr) {
		u = new(models.User{
			Username: "stranger",
		})
	}

	if err := execTemplate(c, u, "webpages/index.html"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) loginHandler(c *gin.Context) {
	var err fverrors.Error

	user, err := auth.GetUserFromSession(c)
	if err != nil && !errors.Is(err, fverrors.NoLoginErr) {
		fverrors.SendErrorResponse(c, err)
		return
	}

	if err == nil {
		if auth.CheckAuth(user, AuthDB) {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
		return
	}

	if err := execTemplate(c, user, "webpages/login.html"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) userCreateHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if (err != nil) && (errors.Is(err, fverrors.NoLoginErr)) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		return
	} else if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	if auth.CheckAuth(user, AuthDB) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	if err := execTemplate(c, user, "webpages/createuser.html"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) userRegisterHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if (err != nil) && (errors.Is(err, fverrors.NoLoginErr)) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		return
	} else if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	user.Username = c.PostForm("username")

	user, err = db.RegisterUser(user)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	err = auth.SaveUserToSession(c, user)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

func (s *Server) userDashboardHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if (err != nil) && (!errors.Is(err, fverrors.NoLoginErr)) {
		fverrors.SendErrorResponse(c, err)
		return
	}

	if !auth.CheckAuth(user, auth.Guest|auth.Db) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
		return
	}

	if err := execTemplate(c, user, "webpages/dashboard.html"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) getUserStoriesHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if (err != nil) && (!errors.Is(err, fverrors.NoLoginErr)) {
		fverrors.SendErrorResponse(c, err)
		return
	}

	username, visibility := c.Param("username"), c.Query("visibility")
	stories, err := db.GetStoriesByUser(user, username, visibility)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, stories)
}

func (s *Server) newStoryHandler(c *gin.Context) {
	story := Story{
		StoryID:     rand.Text()[:15],
		StoryName:   c.PostForm("story_name"),
		Description: c.PostForm("description"),
		Visibility:  -1,
	}

	if visStr := c.PostForm("visibility"); visStr != "" {
		vis, err := strconv.Atoi(visStr[:1])
		if err != nil {
			fmt.Println(err)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		story.Visibility = vis
	}

	storyID, err := s.newStory(c, &story)

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
	story := Story{
		StoryID:     c.PostForm("story_id"),
		StoryName:   c.PostForm("story_name"),
		Description: c.PostForm("description"),
		Visibility:  -1,
	}

	if visStr := c.PostForm("visibility"); visStr != "" {
		vis, err := strconv.Atoi(visStr[:1])
		if err != nil {
			fmt.Println(err)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		story.Visibility = vis
	}

	if err := s.editStory(c, &story); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		if _, err = c.Writer.WriteString(err.Error()); err != nil {
			fmt.Println(err)
		}
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	if _, err := c.Writer.WriteString("Story edited!"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) newChapterHandler(c *gin.Context) {
	chapter := Chapter{
		ChapterID:   rand.Text()[:15],
		ChapterName: c.PostForm("chapter_name"),
		FileID:      rand.Text()[:15],
		StoryID:     c.PostForm("story_id"),
	}

	chapterID, err := s.newChapter(c, &chapter, c.PostForm("content"))

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		if _, err := c.Writer.Write([]byte(err.Error())); err != nil {
			fmt.Println(err)
		}
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	if _, err := c.Writer.Write([]byte(chapterID)); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) editChapterHandler(c *gin.Context) {
	chapter := &Chapter{
		ChapterID:   c.PostForm("chapter_id"),
		ChapterName: c.PostForm("chapter_name"),
	}

	if err := s.editChapter(c, chapter, c.PostForm("content")); err != nil {
		fmt.Println(err)
	}

	c.Writer.WriteHeader(http.StatusOK)
	if _, err := c.Writer.WriteString("Chapter edited!"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) getStoryChaptersHandler(c *gin.Context) {
	chapters, err := s.getStoryChapters(c.PostForm("story_id"), c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	if _, err := c.Writer.WriteString(fmt.Sprintf("%d chapters: %v", len(chapters), chapters)); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) getChapterHandler(c *gin.Context) {
	chapter, err := s.getChapterWithContent(c.Param("chapterID"), c)

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, chapter)
}

func (s *Server) addEdgeHandler(c *gin.Context) {
	from, to := c.PostForm("from_chapter"), c.PostForm("to_chapter")

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.checkEdge(from, to, &u)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err := c.Writer.WriteString(err.Error())
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err = s.addEdge(from, to)

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err := c.Writer.WriteString(err.Error())
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	if _, err := c.Writer.WriteString("Edge added!"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) getEdgesHandler(c *gin.Context) {
	chapterID := c.PostForm("chapter_id")

	edges, err := s.getAllEdges(chapterID)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	if _, err := c.Writer.WriteString(fmt.Sprintf("%d edges: %v", len(edges), edges)); err != nil {
		fmt.Println(err)
	}
}

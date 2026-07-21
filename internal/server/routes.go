package server

import (
	"errors"
	"fmt"
	"forgeverse/internal/api"
	"forgeverse/internal/auth"
	"forgeverse/internal/db"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"
	"html/template"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//TODO: implement efficient content response with fileFromFS(potentially)
//TODO: add location header to userRegisterHandler

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
	{
		userGroup.GET("/mystories", s.userMyStoriesHandler)
	}

	// api routes
	publicApi := r.Group("/api")
	{
		publicApi.GET("/getstories/:username", s.getUserStoriesHandler)
		publicApi.GET("/getmystories", s.getMyStoriesHandler)
		publicApi.GET("/getchapter/:chapter_id", s.getChapterHandler)

		publicApi.GET("/getstorychapters/:story_id", s.getStoryChaptersHandler)

		publicApi.GET("/getedges/:chapter_id", s.getEdgesHandler)
	}

	privateApi := r.Group("/api")
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

	user.AuthState = auth.Google | auth.Db
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
		if auth.CheckAuth(user, auth.Db) {
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

	if auth.CheckAuth(user, auth.Db) {
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

func (s *Server) userMyStoriesHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if (err != nil) && (!errors.Is(err, fverrors.NoLoginErr)) {
		fverrors.SendErrorResponse(c, err)
		return
	}

	if !auth.IsPrivateAuthenticated(user) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/register")
		return
	}

	if err := execTemplate(c, user, "webpages/mystories.html"); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) getUserStoriesHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if (err != nil) && (!errors.Is(err, fverrors.NoLoginErr)) {
		fverrors.SendErrorResponse(c, err)
		return
	}

	stories, err := api.GetStoriesByUser(user, c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stories": stories,
	})
}

func (s *Server) getMyStoriesHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if (err != nil) && (!errors.Is(err, fverrors.NoLoginErr)) {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPrivateAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
		return
	}

	c.AddParam("username", user.Username)
	stories, err := api.GetStoriesByUser(user, c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stories": stories,
	})
}

func (s *Server) newStoryHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPrivateAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
		return
	}

	storyID, err := api.CreateNewStory(c, user)

	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.Header("Location", fmt.Sprintf("%s/story/%s", HomeURL, storyID))
	c.JSON(http.StatusCreated, gin.H{"story_id": storyID})
}

func (s *Server) editStoryHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPrivateAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
		return
	}

	err = api.EditStory(c, user)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (s *Server) newChapterHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPrivateAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
		return
	}

	chapterID, err := api.CreateNewChapter(c, user)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.Header("Location", fmt.Sprintf("%s/chapter/%s", HomeURL, chapterID))
	c.JSON(http.StatusCreated, gin.H{"chapter_id": chapterID})
}

func (s *Server) editChapterHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPrivateAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
		return
	}

	err = api.EditChapter(c, user)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (s *Server) getStoryChaptersHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil && !errors.Is(err, fverrors.NoLoginErr) {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPublicAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
	}

	chapters, err := api.GetChaptersByStory(user, c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"chapters": chapters})
}

func (s *Server) getChapterHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil && !errors.Is(err, fverrors.NoLoginErr) {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPublicAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
	}

	chapter, err := api.GetChapterByID(user, c)

	c.JSON(http.StatusOK, gin.H{"chapter": chapter})
}

func (s *Server) addEdgeHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPrivateAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
		return
	}

	if err := api.CreateNewEdge(c, user); err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

func (s *Server) getEdgesHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil && !errors.Is(err, fverrors.NoLoginErr) {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPublicAuthenticated(user) {
		fverrors.SendErrorResponse(c, fverrors.UnAuthorizedErr)
	}

	edges, err := api.GetEdgesByChapter(user, c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"edges": edges})
}

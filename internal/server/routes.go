package server

import (
	"errors"
	"fmt"
	"forgeverse/internal/api"
	"forgeverse/internal/auth"
	"forgeverse/internal/db"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//TODO: implement efficient content response with fileFromFS(potentially)
//TODO: add location header to userRegisterHandler
//TODO: check error handling in templates

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
	r.GET("/home", s.homeHandler)

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

	dummyApi := r.Group("/dummy")
	{
		dummyApi.GET("/make", s.makeDummyUser)

		dummyApi.POST("/newstory", s.newDummyStoryHandler)

		dummyApi.POST("/newchapter", s.newDummyChapterHandler)

		dummyApi.POST("/addedge", s.addDummyEdgeHandler)
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

/*
	GET /auth/logout

Description: Logout endpoint
*/
func (s *Server) logoutHandler(c *gin.Context) {
	err := auth.Logout(c)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

/*
	GET /auth/google/callback

Description: Google OAuth callback endpoint
*/
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

/*
	GET /auth/google

Description: Google OAuth init endpoint
*/
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

/*
	GET /

Description: Homepage endpoint
*/
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

	servePage(c, "index", u)
}

/*
	GET /auth/login

Description: Login page endpoint
*/
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

	servePage(c, "login", user)
}

/*
	GET /user/register

Description: User registration endpoint
*/
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

	servePage(c, "createuser", user)
}

/*
	POST /user/register

Description: Post username on this endpoint to try to register the user

Form Data:

	username = Username to register. Valid regex: [a-z0-9_]{3,18}

Returns:

	If successful: Status code 201.
*/
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

/*
	GET /user/mystories

Description: Show currently logged-in user's stories
*/
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

	stories, err := api.GetStoriesByUser(user, c, user.Username, "all")
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	servePage(c, "mystories", models.NewMyStoriesData(user, stories))
}

/*
	GET /api/getstories/:username

Description: Get stories by the specified user of the specified visibility

Queries:

	visibility = public | private | all

Returns:

	If successful:
	{
		stories: [
			{
				story_name:	string
				story_id:	string
				description:	string
				visibility:	int
				username:	string
			}
		]
	}
*/
func (s *Server) getUserStoriesHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if (err != nil) && (!errors.Is(err, fverrors.NoLoginErr)) {
		fverrors.SendErrorResponse(c, err)
		return
	}

	stories, err := api.GetStoriesByUser(user, c, c.Param("username"), c.Query("visibility"))
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stories": stories,
	})
}

/*
	GET /api/getmystories

Description: Shortcut for /api/getstories/<logged-in user>. Fails if user is not logged in.
*/
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
	stories, err := api.GetStoriesByUser(user, c, c.Param("username"), c.Query("visibility"))
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stories": stories,
	})
}

/*
	POST /api/newstory

Description: Create new empty story.

Form Data:

	story_name:	string
	description:	string
	visibility:	string

Returns:

	If successful:
	{
		story_id:	string
	}
*/
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

/*
	POST /api/editstory

Description: Edit existing story.

Form Data:

	story_id:	string
	story_name:	string
	description:	string
	visibility:	string

Returns:

	If successful, Status code 200
*/
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

/*
	POST /api/newchapter

Description: Create new chapter with specified contents.

Form Data:

	chapter_name:	string
	story_id:	string
	content:	string

Returns:

	If successful:
	{
		chapter_id:	string
	}
*/
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

/*
	POST /api/editchapter

Description: Create new chapter with specified contents.

Form Data:

	chapter_id:	string
	chapter_name:	string
	content:	string

Returns:

	If successful, Status code 200
*/
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

	c.JSON(http.StatusOK, gin.H{})
}

/*
	GET /api/getstorychapters/:story_id

Description: Get metadata chapters of a given story. Content is not included in response.

Returns:

	If successful:

		{
			chapters: [
				{
					chapter_name:	string
					chapter_id:	string
					story_id:	string
				}
			]
		}
*/
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

/*
	GET /api/getchapter/:chapter_id

Description: Get specifier chapter, along with its contents.

Returns:

	If successful:

		{
			chapters: [
				{
					chapter_name:	string
					chapter_id:	string
					story_id:	string
					content:	string
				}
			]
		}
*/
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

/*
	POST /api/addedge

Description: Add a new edge between chapters.

Returns:

	If successful, Status 201
*/
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

/*
	GET /api/getedge/:chapter_id

Description: Gets all incoming and outgoing edges of specified chapter

Returns:

	If successful:

	{
		edges: [
			{
				from_chapter:	string
				to_chapter:	string
			}
		]
	}
*/
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

func (s *Server) homeHandler(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil && !errors.Is(err, fverrors.NoLoginErr) {
		fverrors.SendErrorResponse(c, err)
		return
	} else if !auth.IsPublicAuthenticated(user) {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		return
	}

	servePage(c, "home", user)
}

package server

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

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

	r.GET("/", s.helloWorldHandler)

	r.GET("/login", s.loginHandler)
	r.GET("/auth/google", s.authGoogleHandler)
	r.GET("/auth/google/callback", s.authGoogleCallbackHandler)
	r.GET("/logout", s.logoutHandler)

	r.GET("/user/create", s.userCreateHandler)
	r.GET("/user/dashboard", s.userDashboardHandler)

	r.POST("/user/create", s.userRegisterHandler)

	r.Static("/css", "webpages/css")
	r.Static("/images", "webpages/images")
	r.Static("/scripts", "webpages/scripts")

	return r
}

func (s *Server) logoutHandler(c *gin.Context) {
	session, errGettingState := s.store.Get(c.Request, "sessionISWP")

	if errGettingState != nil {
		fmt.Println(errGettingState)
		c.Writer.WriteHeader(500)
		return
	}

	session.Options.MaxAge = -1
	session.Values = make(map[any]any)
	errSaving := session.Save(c.Request, c.Writer)
	if errSaving != nil {
		fmt.Println(errSaving)
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
		c.Writer.WriteHeader(500)
		return
	}

	if session.Values["state"] != stateString {
		c.Writer.WriteHeader(401)
		return
	}

	v, err := s.exchangeTokenForUser(code, c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	session.Values["userInfo"] = v
	if err := session.Save(c.Request, c.Writer); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	err = s.getUser(&v)
	if (err != nil) && errors.Is(err, pgx.ErrNoRows) {
		c.Redirect(http.StatusTemporaryRedirect, "/user/create")
		return
	} else if err != nil {
		c.Writer.WriteHeader(500)
		return
	}

	session.Values["userInfo"] = v
	if err := session.Save(c.Request, c.Writer); err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) authGoogleHandler(c *gin.Context) {
	userInfo, err := s.getUserInfoFromSession(c)

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	} else if userInfo.getState(AuthDB) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	} else if userInfo.getState(AuthGoogle) {
		c.Redirect(http.StatusTemporaryRedirect, "/user/create")
		return
	}

	oAuthReqUrl, err := s.setAuthSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	http.Redirect(c.Writer, c.Request, oAuthReqUrl, http.StatusTemporaryRedirect)
}

func (s *Server) helloWorldHandler(c *gin.Context) {
	t, err := template.ParseFiles("./webpages/index.html")

	userInfo, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	if err = t.Execute(c.Writer, userInfo); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) loginHandler(c *gin.Context) {
	userInfo, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	if userInfo.getState(AuthDB) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	t, err := template.ParseFiles("./webpages/login.html")
	if err = t.Execute(c.Writer, nil); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) userCreateHandler(c *gin.Context) {
	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	if err := s.getUser(&u); (err != nil) && (err != pgx.ErrNoRows) {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	} else if err == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	t, err := template.ParseFiles("./webpages/createuser.html")
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	if err = t.Execute(c.Writer, u); err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Server) userRegisterHandler(c *gin.Context) {
	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
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
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

func (s *Server) userDashboardHandler(c *gin.Context) {

}

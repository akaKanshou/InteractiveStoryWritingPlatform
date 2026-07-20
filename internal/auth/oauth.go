package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"forgeverse/internal/models"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	fverrors "forgeverse/internal/errors"
)

var googleOAuthConfig *oauth2.Config

func init() {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENTID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENTSECRET"),
		RedirectURL:  "https://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func InitGoogleOAuth(c *gin.Context) (string, fverrors.Error) {
	session, err := store.Get(c.Request, "forgeVerseSession")
	if err != nil {
		return "", fverrors.NewServerError(err)
	}

	stateString := url.QueryEscape(rand.Text())
	session.Values["state"] = stateString

	if c.Query("remember") == "true" {
		session.Values["remember"] = true
	} else {
		session.Values["remember"] = false
	}

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		return "", fverrors.NewServerError(err)
	}

	oAuthReqUrl := googleOAuthConfig.AuthCodeURL(stateString, oauth2.AccessTypeOffline)

	return oAuthReqUrl, nil
}

func GoogleOAuthCallback(c *gin.Context) fverrors.Error {
	session, err := store.Get(c.Request, "forgeVerseSession")
	if err != nil {
		return fverrors.NewServerError(err)
	}

	code, stateString := c.Query("code"), c.Query("state")

	if session.Values["state"] != stateString {
		return fverrors.NewInvalidRequestError("state of callback does not match", errors.New("invalid state"))
	}

	token, err := googleOAuthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		return fverrors.NewServerError(err)
	}

	authClient := googleOAuthConfig.Client(context.Background(), token)

	res, err := authClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return fverrors.NewServerError(err)
	}
	defer res.Body.Close()

	user := new(models.User{})
	err = json.NewDecoder(res.Body).Decode(user)
	if err != nil {
		return fverrors.NewServerError(err)
	}

	user.AuthState = Google
	user.Remember = session.Values["remember"].(bool)

	if err := SaveUserToSession(c, user); err != nil {
		return err
	}

	return nil
}

func Logout(c *gin.Context) fverrors.Error {
	session, err := store.Get(c.Request, "forgeVerseSession")
	if err != nil {
		return fverrors.NewServerError(err)
	}

	session.Options.MaxAge = -1
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		return fverrors.NewServerError(err)
	}

	return nil
}

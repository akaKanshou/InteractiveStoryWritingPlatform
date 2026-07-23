package auth

import (
	"crypto/rand"
	"encoding/gob"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"github.com/gorilla/sessions"

	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"
)

var store *sessions.CookieStore

func init() {
	gob.Register(models.User{})

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(7 * 86400)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false

}

type State = uint8

const (
	NoAccess State = 0x1
	Guest    State = 0x2
	Google   State = 0x04
	Db       State = 0x08
	Admin    State = 0x10
)

func SaveUserToSession(c *gin.Context, user *models.User) fverrors.Error {
	session, err := store.Get(c.Request, "forgeVerseSession")
	if err != nil {
		return fverrors.NewServerError(err)
	}

	session.Values["user"] = *user

	if user.Remember {
		session.Options.MaxAge = 86400 * 7
	} else {
		session.Options.MaxAge = 86400
	}

	if err := session.Save(c.Request, c.Writer); err != nil {
		return fverrors.NewServerError(err)
	}

	return nil
}

func GetUserFromSession(c *gin.Context) (*models.User, fverrors.Error) {
	session, err := store.Get(c.Request, "forgeVerseSession")
	if err != nil {
		return nil, fverrors.NewServerError(err)
	}

	userAnyCheck := session.Values["user"]
	if userAnyCheck != nil {
		if user, ok := userAnyCheck.(models.User); ok && !CheckAuth(&user, Guest) {
			return &user, nil
		} else {
			return &user, fverrors.NoLoginErr
		}
	}

	user := new(models.User{
		Username:  "Guest" + rand.Text()[:12],
		AuthState: Guest,
		Remember:  true,
	})

	if err := SaveUserToSession(c, user); err != nil {
		return nil, fverrors.NewServerError(err)
	}

	return user, fverrors.NoLoginErr
}

func CheckAuth(user *models.User, authState State) bool {
	return user.AuthState&authState > 0
}

func IsPublicAuthenticated(user *models.User) bool {
	if CheckAuth(user, NoAccess) {
		return false
	}

	return CheckAuth(user, Admin|Db|Guest)
}

func IsPrivateAuthenticated(user *models.User) bool {
	if CheckAuth(user, NoAccess) {
		return false
	}

	return CheckAuth(user, Admin|Db)
}

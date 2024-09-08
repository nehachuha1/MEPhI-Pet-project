package session

import (
	"errors"
	"github.com/labstack/echo/v4"
	"mephiMainProject/pkg/services/server/config"
)

var (
	currentConfig = config.NewConfig()
	ErrorNoAuth   = errors.New("No session found")
	sessionKey    = currentConfig.Sess.SessionKey
	jwtSecretKey  = []byte(currentConfig.Sess.JwtKey)
)

func NewSession(username string) *config.Session {
	return &config.Session{
		Username: username,
	}
}

func SessionFromContext(ctx echo.Context) (*config.Session, error) {
	sess, ok := ctx.Get(sessionKey).(*config.Session)
	if !ok || sess == nil {
		return nil, ErrorNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx echo.Context, sess *config.Session) {
	ctx.Set(sessionKey, sess)
}

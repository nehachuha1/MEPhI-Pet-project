package session

import (
	"context"
	"errors"
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

func SessionFromContext(ctx context.Context) (*config.Session, error) {
	sess, ok := ctx.Value(sessionKey).(*config.Session)
	if !ok || sess == nil {
		return nil, ErrorNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *config.Session) context.Context {
	return context.WithValue(ctx, sessionKey, sess)
}

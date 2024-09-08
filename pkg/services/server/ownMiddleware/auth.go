package ownMiddleware

import (
	"github.com/labstack/echo/v4"
	"mephiMainProject/pkg/services/server/session"
)

func Auth(sm *session.SessionManager) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			currentSession, err := sm.Check(c)
			if err == nil || currentSession != nil {
				session.ContextWithSession(c, currentSession)
				_ = next(c)
			} else {
				_ = next(c)
			}
			return nil
		}
	}
}

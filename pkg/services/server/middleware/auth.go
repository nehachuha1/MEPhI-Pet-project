package middleware

import (
	"mephiMainProject/pkg/services/server/session"
	"net/http"
)

func Auth(sm *session.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentSession, err := sm.Check(w, r)
			if err == nil || currentSession != nil {
				ctx := session.ContextWithSession(r.Context(), currentSession)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

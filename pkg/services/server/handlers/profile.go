package handlers

import (
	"go.uber.org/zap"
	"mephiMainProject/pkg/services/server/profile"
	"mephiMainProject/pkg/services/server/session"
)

type ProfileHandler struct {
	Logger      *zap.SugaredLogger
	Sessions    *session.SessionManager
	ProfileRepo profile.ProfileRepo
}

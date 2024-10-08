package session

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"log"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/database"
	"time"
)

type SessionManager struct {
	dbORM         *database.DatabaseORM
	currentConfig *config.Config
}

func NewSessionManager(db *database.DatabaseORM, cfg *config.Config) *SessionManager {
	return &SessionManager{
		dbORM:         db,
		currentConfig: cfg,
	}
}

func CreateNewToken(user config.User, sessionId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username":   user.Login,
			"session_id": sessionId,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * 60 * 60 * 24 * 3).Unix(),
	})
	return token.SignedString(jwtSecretKey)
}

func (sm *SessionManager) Check(c echo.Context) (*config.Session, error) {
	tokenWithCookie, err := c.Cookie("session")
	if err != nil {
		log.Printf("Check error: there's no auth cookie")
		return nil, ErrorNoAuth
	}
	tokenString := tokenWithCookie.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrorNoAuth
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, ErrorNoAuth
	}
	user, isOk := (claims["user"]).(map[string]interface{})
	if !isOk {
		return nil, ErrorNoAuth
	}
	sessionIdFromToken, success := (user["session_id"]).(string)
	if !success {
		return nil, ErrorNoAuth
	}
	sessID := &config.SessionID{ID: sessionIdFromToken}
	currentSession, err := sm.dbORM.CheckSession(sessID)
	if err != nil {
		log.Printf("Current session check err - %v", err)
		return nil, ErrorNoAuth
	}
	return currentSession, nil
}

func (sm *SessionManager) Create(login string) (*config.Session, error) {
	newSession := NewSession(login)
	sessionWithID, err := sm.dbORM.CreateSession(newSession)
	if err != nil {
		log.Printf("SessionManager.Create error - %v", err)
		return &config.Session{}, nil
	}
	checkedSession, err := sm.dbORM.CheckSession(&sessionWithID.SessID)
	if err != nil {
		log.Printf("Check session err - %v", err)
		return &config.Session{}, nil
	}
	return checkedSession, nil
}

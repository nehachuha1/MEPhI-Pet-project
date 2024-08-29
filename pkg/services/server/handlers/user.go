package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/session"
	"mephiMainProject/pkg/services/server/user"
	"net/http"
	"time"
)

type UserHandler struct {
	Logger   *zap.SugaredLogger
	Sessions *session.SessionManager
	UserRepo user.UserRepo
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant read request body")
	}
	authData := &config.UserAuthData{}
	err = json.Unmarshal(body, authData)

	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	usrAuthData, err := u.UserRepo.Authorize(authData.Login, authData.Password)
	if err != nil { // формируем ошибку при регистрации
		authErrResp(w, "username", authData.Login, err)
		return
	}

	sess, err := u.Sessions.Create(w, usrAuthData.Login)

	if err != nil {
		http.Error(w, `Session isn't create`+err.Error(), http.StatusInternalServerError)
		return
	}
	u.Logger.Infof("Successfully created session for username %v", sess.Username)
	usr := config.User{Login: sess.Username}
	token, err := session.CreateNewToken(usr, sess.SessID.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp, err := json.Marshal(map[string]interface{}{
		"token":       token,
		"Status code": http.StatusFound,
	})
	newCookie := http.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(time.Second * 60 * 60 * 24 * 3),
	}
	http.SetCookie(w, &newCookie)

	CheckMarshalError(w, err, resp)
	u.Logger.Infof("Send token on client for user with username: %v ", sess.Username)

}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant read request body")
	}
	authData := &config.UserAuthData{}
	err = json.Unmarshal(body, authData)

	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	_, err = u.UserRepo.Register(authData.Login, authData.Password)
	if err != nil { // формируем ошибку при регистрации
		authErrResp(w, "username", authData.Login, err)
		return
	}
	sess, err := u.Sessions.Create(w, authData.Login)
	if err != nil {
		http.Error(w, `Session isn't create`+err.Error(), http.StatusInternalServerError)
		return
	}
	u.Logger.Infof("Successfully created session for user with ID %v. Full session: %#v", sess.Username, sess)

	usr := config.User{Login: sess.Username}
	token, err := session.CreateNewToken(usr, sess.SessID.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp, err := json.Marshal(map[string]interface{}{
		"token":       token,
		"Status code": http.StatusFound,
	})
	newCookie := http.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(time.Second * 60 * 60 * 24 * 3),
	}
	http.SetCookie(w, &newCookie)

	CheckMarshalError(w, err, resp)
	u.Logger.Infof("Send token on client for user with username: %v ", sess.Username)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	resp, err := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	CheckMarshalError(w, err, resp)
}

func CheckMarshalError(w http.ResponseWriter, err error, resp []byte) {
	if err != nil {
		http.Error(w, "Marshaling error", http.StatusBadRequest)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "Writing response err", http.StatusInternalServerError)
		return
	}
}

func authErrResp(w http.ResponseWriter, param string, value string, err error) {
	var (
		resp  []byte
		error error
	)
	w.WriteHeader(http.StatusUnprocessableEntity)
	errors := make([]map[string]string, 0)
	errors = append(errors, map[string]string{
		"location": "body",
		"param":    param,
		"value":    value,
		"msg":      err.Error(),
	})
	resp, error = json.Marshal(map[string][]map[string]string{
		"errors": errors,
	})

	CheckMarshalError(w, error, resp)
}

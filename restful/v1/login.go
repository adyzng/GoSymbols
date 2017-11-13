package v1

import (
	"fmt"
	"net/http"

	"github.com/adyzng/GoSymbols/restful"
	"github.com/adyzng/GoSymbols/restful/auth"
	"github.com/adyzng/GoSymbols/restful/session"
	log "gopkg.in/clog.v1"
)

// check if login required
//
func loginRequired(r *http.Request) (string, *auth.GraphToken) {
	ssid := ""
	if s, _ := r.Cookie(session.CookieSessID); s != nil {
		ssid = s.Value
	}
	if data := session.GetManager().Get(ssid); data != nil {
		return ssid, data.(*auth.GraphToken)
	}
	return ssid, nil
}

// AuthLogin login by oauth to Arcserve domain
//	[:]/auth/login
//
func AuthLogin(w http.ResponseWriter, r *http.Request) {
	redirect := "/"
	if _, token := loginRequired(r); token == nil {
		redirect = auth.AuthURL()
	} else {
		log.Info("[Login] User %s already logined.", token.UserName)
	}

	log.Trace("[Login] Redirect URL: %s.", redirect)
	w.Header().Set("Location", redirect)
	w.WriteHeader(http.StatusFound)
}

// AuthLogout user logout
//   [:]/auth/logout
//
func AuthLogout(w http.ResponseWriter, r *http.Request) {
	sess, _ := r.Cookie(session.CookieSessID)
	if sess == nil {
		log.Trace("[Logout] Session ID is empty.")
		w.WriteHeader(http.StatusOK)
		return
	}
	if data := session.GetManager().Delete(sess.Value); data != nil {
		if token, ok := data.(*auth.GraphToken); ok {
			log.Info("[Logout] user %s.", token.UserName)
		}
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

// Authorize handle response from the microsoft online authorize
//  [:]/auth/authorize
//
func Authorize(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code, state := r.FormValue("code"), r.FormValue("state")

	errType, errDesc := r.FormValue("error"), r.FormValue("error_description")
	if errType != "" {
		log.Error(2, "[Login] Authorize error: %s, desc: %s.", errType, errDesc)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if code == "" || state == "" {
		log.Warn("[Login] Empty auth code.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := auth.QueryToken(code, state)
	if err != nil {
		res := restful.RestResponse{ErrCodeMsg: restful.ErrLoginFailed}
		res.Message = fmt.Sprintf("%s", err)

		res.WriteJSON(w)
		log.Error(2, "[Login] Query token failed: %v.", err)
		return
	}

	sessID := session.GetManager().Create(token)
	if user, _ := auth.GetUserProfile("", token); user != nil {
		token.UserName = user.DisplayName
		log.Info("[Login] User (%s) login succeed.", user.DisplayName)
	}

	// set http cookie
	http.SetCookie(w, &http.Cookie{
		Name:     session.CookieSessID,
		MaxAge:   int(session.CookieMaxAge),
		HttpOnly: false,
		Value:    sessID,
		Path:     "/",
	})

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

// GetUserProfile get user information
//  [:]/user/profile
//
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ssid, token := loginRequired(r)
	if token == nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Warn("[User] Login required.")
		return
	}

	resp := restful.RestResponse{}
	if user, err := auth.GetUserProfile(ssid, token); err == nil {
		resp.Data = user
		log.Trace("[User] Get user profile for %s.", user.DisplayName)
	} else {
		resp.ErrCodeMsg = restful.ErrLoginNeeded
		resp.Message = fmt.Sprintf("%s", err)
	}

	resp.WriteJSON(w)
}

// GetUserPhoto get user profile photo
//   [:]/user/photo
//
func GetUserPhoto(w http.ResponseWriter, r *http.Request) {
	ssid, token := loginRequired(r)
	if token == nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Warn("[User] Login required.")
		return
	}
	if err := auth.GetUserPhoto(ssid, token, w); err != nil {
		log.Error(2, "[User] Get user photo failed: %v.", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}

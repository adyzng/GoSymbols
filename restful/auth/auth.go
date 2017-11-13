package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/adyzng/GoSymbols/restful/session"

	"github.com/adyzng/GoSymbols/config"
	"github.com/adyzng/GoSymbols/restful"
	log "gopkg.in/clog.v1"
)

const (
	// MicrosoftAuthURI restful api
	adAuthURI  = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	adTokenURI = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	graphURL   = "https://graph.microsoft.com/v1.0"
)

func getURL(typ string) string {
	switch typ {
	case "me":
		return graphURL + "/me"
	case "photo":
		return graphURL + "/me/photo/$value"
	default:
		return graphURL
	}
}

// AuthURL combine the auth url
//
func AuthURL() string {
	if location, err := url.Parse(adAuthURI); err == nil {
		params := location.Query()
		params.Add("client_id", config.ClientID)
		params.Add("redirect_uri", config.RedirectURI)
		params.Add("response_type", "code")
		params.Add("response_mode", "form_post")
		params.Add("scope", config.GraphScope)
		params.Add("state", fmt.Sprintf("%d", time.Now().Unix()))
		location.RawQuery = params.Encode()
		return location.String()
	}
	return ""
}

// QueryToken get token from provider
//
/***Example
POST /{tenant}/oauth2/v2.0/token HTTP/1.1
Host: https://login.microsoftonline.com
Content-Type: application/x-www-form-urlencoded

client_id=6731de76-14a6-49ae-97bc-6eba6914391e
&scope=https%3A%2F%2Fgraph.microsoft.com%2Fmail.read
&code=OAAABAAAAiL9Kn2Z27UubvWFPbm0gLWQJVzCTE9UkP3pSx1aXxUjq3n8b2JRLk4OxVXr...
&redirect_uri=http%3A%2F%2Flocalhost%2Fmyapp%2F
&grant_type=authorization_code
&client_secret=JqQX2PNo9bpM0uEihUPzyrh
*/
func QueryToken(code, state string) (*GraphToken, error) {
	payload := url.Values{
		"client_id":     {config.ClientID},
		"redirect_uri":  {config.RedirectURI},
		"client_secret": {config.ClientKey},
		"scope":         {config.GraphScope},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	}

	buff, err := restful.HttpPost(adTokenURI, strings.NewReader(payload.Encode()), nil)
	if err != nil {
		if buff != nil {
			gErr := GraphError{}
			if json.NewDecoder(buff).Decode(&gErr); gErr.Error != "" {
				log.Warn("[Auth] Request Access Token error: %+v.", gErr)
				err = errors.New(gErr.Description)
			}
		}
		log.Error(2, "[Auth] Request Access Token failed: %v.", err)
		return nil, err
	}
	//log.Info("[Auth] Response %s.", string(buff.Bytes()))
	token := GraphToken{}
	if err := json.NewDecoder(buff).Decode(&token); err != nil {
		log.Error(2, "[Auth] Decode token failed: %v.", err)
		return nil, err
	}
	if len(token.AccessToken) == 0 {
		log.Warn("[Auth] Get invalid token.")
		return nil, fmt.Errorf("invalid access token")
	}

	log.Trace("[Auth] Token: %s %s.", token.Type, token.AccessToken[:10])
	token.State = state
	token.ExpireAt = time.Now().Unix() + token.ExpireAt
	return &token, nil
}

// RefreshToken refresh token if access token is expired
//
/***Refresh Example:
POST /{tenant}/oauth2/v2.0/token HTTP/1.1
Host: https://login.microsoftonline.com
Content-Type: application/x-www-form-urlencoded

client_id=6731de76-14a6-49ae-97bc-6eba6914391e
&scope=https%3A%2F%2Fgraph.microsoft.com%2Fmail.read
&refresh_token=OAAABAAAAiL9Kn2Z27UubvWFPbm0gLWQJVzCTE9UkP3pSx1aXxUjq...
&redirect_uri=http%3A%2F%2Flocalhost%2Fmyapp%2F
&grant_type=refresh_token
&client_secret=JqQX2PNo9bpM0uEihUPzyrh
*/
func RefreshToken(token *GraphToken) (*GraphToken, error) {
	payload := url.Values{
		"client_id":     {config.ClientID},
		"redirect_uri":  {config.RedirectURI},
		"client_secret": {config.ClientKey},
		"scope":         {config.GraphScope},
		"refresh_token": {token.RefreshToken},
		"grant_type":    {"refresh_token"},
	}

	buff, err := restful.HttpPost(adTokenURI, strings.NewReader(payload.Encode()), nil)
	if err != nil {
		if buff != nil {
			gErr := GraphError{}
			if json.NewDecoder(buff).Decode(&gErr); gErr.Error != "" {
				log.Warn("[Auth] Refresh Token error: %+v.", gErr)
				err = errors.New(gErr.Description)
			}
		}
		log.Error(2, "[Auth] Refresh Token failed: %v.", err)
		return nil, err
	}

	var tokenNew GraphToken
	if err := json.NewDecoder(buff).Decode(&tokenNew); err != nil {
		log.Error(2, "[Auth] Decode token failed: %v.", err)
		return nil, err
	}
	log.Info("[Auth] Refresh token (%+v).", tokenNew)

	tokenNew.State = token.State
	tokenNew.ExpireAt = time.Now().Unix() + tokenNew.ExpireAt
	return &tokenNew, nil
}

func refreshToken(sessID string, token *GraphToken) *GraphToken {
	//
	// if token expired, add 10s extra cost.
	//
	if token.ExpireAt != time.Now().Unix()+10 {
		return token
	}
	log.Info("[User] Refresh token for %s.", token.UserName)
	if tokenNew, _ := RefreshToken(token); tokenNew != nil {
		if sessID != "" {
			session.GetManager().Set(sessID, tokenNew)
		}
		return tokenNew
	}
	return token
}

// GetUserProfile request user profile by access token
//
func GetUserProfile(sessID string, token *GraphToken) (*GraphUser, error) {
	token = refreshToken(sessID, token)

	buff, err := restful.HttpGet(getURL("me"), func(req *http.Request) {
		str := fmt.Sprintf("%s %s", token.Type, token.AccessToken)
		req.Header.Set("Authorization", str)
	})

	if err != nil {
		if buff != nil {
			var gErr *GraphError
			if json.NewDecoder(buff).Decode(&gErr); gErr != nil {
				log.Warn("[Auth] Request profile error: %+v.", gErr)
				err = errors.New(gErr.Description)
			}
		}
		log.Error(2, "[Auth] Request profile failed: %v.", err)
		return nil, err
	}

	user := GraphUser{}
	if err := json.NewDecoder(buff).Decode(&user); err != nil {
		log.Error(2, "[Auth] Decode user info failed: %v.", err)
		return nil, err
	}
	return &user, nil
}

// GetUserPhoto get user photo with access token
//
// refer: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/profilephoto_get
//
func GetUserPhoto(sessID string, token *GraphToken, w io.Writer) error {
	token = refreshToken(sessID, token)
	buff, err := restful.HttpGet(getURL("photo"), func(req *http.Request) {
		str := fmt.Sprintf("%s %s", token.Type, token.AccessToken)
		req.Header.Set("Authorization", str)
	})

	if err != nil {
		if buff != nil {
			var gErr *GraphError
			if json.NewDecoder(buff).Decode(&gErr); gErr != nil {
				log.Warn("[Auth] Request user photo error: %+v.", gErr)
				err = errors.New(gErr.Description)
			}
		}
		log.Error(2, "[Auth] Request user photo failed: %v.", err)
		return err
	}

	_, err = io.Copy(w, buff)
	return err
}

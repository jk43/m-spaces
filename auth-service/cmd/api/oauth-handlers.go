package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

type OrgID = string
type OAuthProvider = string

var oauthHandlers = map[OrgID]map[OAuthProvider]OAuthHandler{}

type OAuthHandler interface {
	PreOAuth(w http.ResponseWriter, r *http.Request) string
	PostOAuth(w http.ResponseWriter, r *http.Request) (*service.User, error)
}

func getOAuthHandler(r *http.Request) (OAuthHandler, error) {
	org, err := service.GetOrg(r, "")
	provider := chi.URLParam(r, "provider")
	if err != nil {
		return nil, err
	}

	handler, ok := oauthHandlers[org.ID][provider]
	if ok {
		return handler, nil
	}

	if _, ok := oauthHandlers[org.ID]; !ok {
		oauthHandlers[org.ID] = make(map[OAuthProvider]OAuthHandler)
	}

	if provider == "google" {
		handler = &GoogleOAuth{
			Config: &oauth2.Config{
				ClientID:     org.Settings.Auth.OAuth[provider].ID,
				ClientSecret: org.Settings.Auth.OAuth[provider].Secret,
				RedirectURL:  org.Settings.Auth.OAuth[provider].Redirect,
				Scopes:       org.Settings.Auth.OAuth[provider].Scopes,
				Endpoint:     google.Endpoint,
			},
		}
		oauthHandlers[org.ID][provider] = handler
	}
	if provider == "microsoft" {
		handler = &MicrosoftOAuth{
			Config: &oauth2.Config{
				ClientID:     org.Settings.Auth.OAuth[provider].ID,
				ClientSecret: org.Settings.Auth.OAuth[provider].Secret,
				RedirectURL:  org.Settings.Auth.OAuth[provider].Redirect,
				Scopes:       org.Settings.Auth.OAuth[provider].Scopes,
				Endpoint:     microsoft.AzureADEndpoint(org.Settings.Auth.OAuth[provider].Options["tenant"]),
			},
		}
		oauthHandlers[org.ID][provider] = handler
	}
	if provider == "facebook" {
		handler = &FacebookOAuth{
			Config: &oauth2.Config{
				ClientID:     org.Settings.Auth.OAuth[provider].ID,
				ClientSecret: org.Settings.Auth.OAuth[provider].Secret,
				RedirectURL:  org.Settings.Auth.OAuth[provider].Redirect,
				Scopes:       org.Settings.Auth.OAuth[provider].Scopes,
				Endpoint:     facebook.Endpoint,
			},
		}
		oauthHandlers[org.ID][provider] = handler
	}
	if provider == "x" {
		handler = &XOAuth{
			Config: &oauth2.Config{
				ClientID:     org.Settings.Auth.OAuth[provider].ID,
				ClientSecret: org.Settings.Auth.OAuth[provider].Secret,
				RedirectURL:  org.Settings.Auth.OAuth[provider].Redirect,
				Scopes:       org.Settings.Auth.OAuth[provider].Scopes,
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://twitter.com/i/oauth2/authorize",
					TokenURL: "https://api.twitter.com/2/oauth2/token",
				},
			},
		}
		oauthHandlers[org.ID][provider] = handler
	}
	return handler, nil
}

// https://console.cloud.google.com/auth/clients
type GoogleOAuth struct {
	Config *oauth2.Config
}

func (h *GoogleOAuth) PreOAuth(w http.ResponseWriter, r *http.Request) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c := &http.Cookie{Name: "oauthstate", Value: state}
	http.SetCookie(w, c)
	url := h.Config.AuthCodeURL(state)
	return url
}

func (h *GoogleOAuth) PostOAuth(w http.ResponseWriter, r *http.Request) (*service.User, error) {
	token, err := h.Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return nil, err
	}
	client := h.Config.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userData := utils.MapStringAny{}
	err = json.Unmarshal(data, &userData)
	user := service.User{
		Email:          userData["email"].(string),
		FirstName:      userData["given_name"].(string),
		LastName:       userData["family_name"].(string),
		ProfileImage:   userData["picture"].(string),
		RegisterMethod: service.RegisterMethodOAuthGoogle,
	}
	return &user, nil
}

type MicrosoftOAuth struct {
	Config *oauth2.Config
}

func (h *MicrosoftOAuth) PreOAuth(w http.ResponseWriter, r *http.Request) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c := &http.Cookie{Name: "oauthstate", Value: state}
	http.SetCookie(w, c)
	url := h.Config.AuthCodeURL(state)
	return url
}

func (h *MicrosoftOAuth) PostOAuth(w http.ResponseWriter, r *http.Request) (*service.User, error) {
	token, err := h.Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return nil, err
	}
	client := h.Config.Client(context.Background(), token)
	response, err := client.Get("https://graph.microsoft.com/v1.0/me")
	utils.TermDebugging(`response`, response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userData := utils.MapStringAny{}
	err = json.Unmarshal(data, &userData)
	utils.TermDebugging(`data`, string(data))
	user := service.User{
		Email:     userData["mail"].(string),
		FirstName: userData["surname"].(string),
		LastName:  userData["givenName"].(string),
		//ProfileImage:   userData["picture"].(string),
		RegisterMethod: service.RegisterMethodOAuthGoogle,
	}
	return &user, nil
}

// https://developers.facebook.com/
type FacebookOAuth struct {
	Config *oauth2.Config
}

func (h *FacebookOAuth) PreOAuth(w http.ResponseWriter, r *http.Request) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c := &http.Cookie{Name: "oauthstate", Value: state}
	http.SetCookie(w, c)
	url := h.Config.AuthCodeURL(state)
	return url
}

func (h *FacebookOAuth) PostOAuth(w http.ResponseWriter, r *http.Request) (*service.User, error) {
	token, err := h.Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return nil, err
	}
	client := h.Config.Client(context.Background(), token)
	response, err := client.Get("https://graph.facebook.com/v19.0/me?fields=id,name,email,picture,first_name,last_name")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userData := utils.MapStringAny{}
	err = json.Unmarshal(data, &userData)
	if err != nil {
		return nil, err
	}
	user := service.User{
		Email:        userData["email"].(string),
		FirstName:    userData["first_name"].(string),
		LastName:     userData["last_name"].(string),
		ProfileImage: userData["picture"].(map[string]any)["data"].(map[string]any)["url"].(string),
	}
	return &user, nil
}

type XOAuth struct {
	Config *oauth2.Config
}

func (h *XOAuth) PreOAuth(w http.ResponseWriter, r *http.Request) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c := &http.Cookie{Name: "oauthstate", Value: state}
	http.SetCookie(w, c)
	url := h.Config.AuthCodeURL(state)
	return url
}

func (h *XOAuth) PostOAuth(w http.ResponseWriter, r *http.Request) (*service.User, error) {
	token, err := h.Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return nil, err
	}
	client := h.Config.Client(context.Background(), token)
	response, err := client.Get("https://api.x.com/2/users/me")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	userData := utils.MapStringAny{}
	err = json.Unmarshal(data, &userData)
	utils.TermDebugging(`userData`, userData)
	user := service.User{
		Email:     userData["email"].(string),
		FirstName: userData["first_name"].(string),
		LastName:  userData["last_name"].(string),
	}
	return &user, nil
}

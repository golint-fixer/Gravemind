package twitch

import (
	"net/http"
	"net/url"
)

type api struct {
	*apiSegment

	// API subpaths
	OAuth *oauth

	// Variables
	Client       *http.Client
	ClientId     string
	ClientSecret string
	CallbackURL  string
}

var API = NewAPI()

func NewAPI() *api {
	a := &api{
		path:   "https://api.twitch.tv/kraken",
		Client: http.DefaultClient,
	}
	a.OAuth = newOAuth(a)

	return a
}

type baseResponse struct {
	Token *tokenResponse `json:"token,omitempty"`
}
type tokenResponse struct {
	Valid         bool                   `json:"valid,omitempty"`
	UserName      string                 `json:"user_name,omitempty"`
	Authorization *authorizationResponse `json:"authorization,omitempty"`
}
type authorizationResponse struct {
	CreatedAt timestamp `json:"created_at,omitempty"`
	UpdatedAt timestamp `json:"updated_at,omitempty"`
	Scopes    []string  `json:"scopes,omitempty"`
}

func (a *api) Base(token string) (*baseResponse, error) {
	v := &baseResponse{}
	err := a.fetch(v, "GET", "/", url.Values{
		"oauth_token": []string{token},
	})
	return v, err
}

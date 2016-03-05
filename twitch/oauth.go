package twitch

import (
	"crypto/rand"
	"encoding/base64"
	"net/url"
)

type oauth struct {
	*apiSegment
}

func newOAuth(parent *apiSegment) *oauth {
	o := &oauth{
		parent: parent,
		path:   "/oauth",
	}

	return o
}

func (o *oauth) Authorize(scopes []string) (url string, state string, err error) {
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	state = base64.URLEncoding.EncodeToString(b)

	path := "/authorize"
	for s := o; s != nil; s = s.parent {
		path = s.path + path
	}

	url = path + "?" + url.Values{
		"response_type": []string{"code"},
		"client_id":     []string{o.root().ClientId},
		"redirect_uri":  []string{o.root().CallbackURL},
		"scope":         scopes,
		"state":         []string{state},
	}.Encode()
}

type tokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
}

func (o *oauth) Token(code string, state string) (token string, err error) {
	r := &tokenResponse{}
	err = o.fetch(r, "POST", "/token", url.Values{
		"grant_type":    []string{"authorization_code"},
		"client_id":     []string{o.root().ClientId},
		"client_secret": []string{o.root().ClientSecret},
		"redirect_uri":  []string{o.root().CallbackURL},
		"code":          []string{code},
		"state":         []string{state},
	})
	token = r.AccessToken
}

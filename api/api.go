package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"

	"github.com/fugiman/tyrantbot/twitch"
)

func init() {
	twitch.API.CallbackURL = "https://api.svipul.net/callback"
}

func home(c web.C, w http.ResponseWriter, r *http.Request) {
	var userId string
	session := c.Env["session"].(*sessions.Session)
	if v, ok := session.Values["userid"]; ok {
		userId = v.(string)
	} else {
		url, state, err := twitch.API.OAuth.Authorize([]string{})
		if err != nil {
			fmt.Fprintf(w, "Error building redirect url: %s", err)
			return
		}

		session.Values["state"] = state
		err = session.Save(r, w)
		if err != nil {
			fmt.Fprintf(w, "Error saving session: %s", err)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "大丈夫 %s", userId)
}

func login(c web.C, w http.ResponseWriter, r *http.Request) {
	// Ensure this isn't a hijacked auth flow
	session := c.Env["session"].(*sessions.Session)
	if session.Values["state"].(string) != c.URLParams["state"] {
		fmt.Fprintf(w, "Error: State mismatch")
		return
	}

	// Exchange the code for a token
    token, err := twitch.API.OAuth.Token(c.URLParams["code"], c.URLParams["state"])
    if err != nil {
        fmt.Fprintf(w, "Error retrieving token: %s", err)
        return
    }

	// Look up who this token belongs to
	base, err := twitch.API.Base(token)
    if err != nil {
        fmt.Fprintf(w, "Error fetching username: %s", err)
        return
    }
    if base.

	// Save the userid & token in the session
	err = session.Save(r, w)
	if err != nil {
		fmt.Fprintf(w, "Error saving session: %s", err)
		return
	}

	fmt.Fprintf(w, "大丈夫")
}

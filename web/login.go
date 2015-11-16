package web

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/session"
)

func loginForm(w http.ResponseWriter, r *http.Request) {
	// Initialise the session service
	sessionService := session.NewService(
		theService.cnf,
		r,
		w,
		theService.oauthService,
	)
	if err := sessionService.InitUserSession(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	renderTemplate(w, "login.tmpl", map[string]interface{}{
		"error": sessionService.GetFlashMessage(),
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	// Initialise the session service
	sessionService := session.NewService(
		theService.cnf,
		r,
		w,
		theService.oauthService,
	)
	if err := sessionService.InitUserSession(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the client
	client, err := theService.oauthService.FindClientByClientID(
		r.Form.Get("client_id"),
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Authenticate the user
	user, err := theService.oauthService.AuthUser(
		r.Form.Get("email"),
		r.Form.Get("password"),
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Default scope
	scope := "read_write"

	// Grant an access token
	accessToken, err := theService.oauthService.GrantAccessToken(
		client,
		user,
		scope,
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Get a refresh token
	refreshToken, err := theService.oauthService.GetOrCreateRefreshToken(
		client,
		user,
		scope,
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Log in the user and store the user session in a cookie
	if err := sessionService.LogIn(&session.UserSession{
		Client:       client,
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}); err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Redirect to the authorize page
	http.Redirect(w, r, "/web/authorize?"+r.URL.Query().Encode(), http.StatusFound)
}

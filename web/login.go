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
	if err := sessionService.InitSession("user_session"); err != nil {
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
	if err := sessionService.InitSession("user_session"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the submitted form data
	r.ParseForm()
	username := r.Form["email"][0]
	password := r.Form["password"][0]

	// Fetch the trusted client
	client, err := theService.oauthService.FindClientByClientID(
		theService.cnf.TrustedClient.ClientID,
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Authenticate the user
	user, err := theService.oauthService.AuthUser(username, password)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, "/web/login", http.StatusFound)
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
		http.Redirect(w, r, "/web/login", http.StatusFound)
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
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Log in the user and store the user session in a cookie
	if err := sessionService.LogIn(&session.UserSession{
		UserID:       user.ID,
		Username:     user.Username,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}); err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Redirect to the authorize page
	http.Redirect(w, r, "/web/authorize", http.StatusFound)
}

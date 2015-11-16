package web

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/session"
)

func registerForm(w http.ResponseWriter, r *http.Request) {
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
	renderTemplate(w, "register.tmpl", map[string]interface{}{
		"error": sessionService.GetFlashMessage(),
	})
}

func register(w http.ResponseWriter, r *http.Request) {
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

	// Check that the submitted email hasn't been registered already
	if theService.oauthService.UserExists(r.Form.Get("email")) {
		sessionService.SetFlashMessage("Email already taken")
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Create a user
	_, err := theService.oauthService.CreateUser(
		r.Form.Get("email"),
		r.Form.Get("password"),
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Redirect to the login page
	http.Redirect(w, r, "/web/login?"+r.URL.Query().Encode(), http.StatusFound)
}

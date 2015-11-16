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
	if err := sessionService.InitSession("user_session"); err != nil {
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
	if err := sessionService.InitSession("user_session"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the submitted form data
	r.ParseForm()
	username := r.Form["email"][0]
	password := r.Form["password"][0]

	// Check that the submitted email hasn't been registered already
	if theService.oauthService.UserExists(username) {
		sessionService.SetFlashMessage("Email already taken")
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	// Create a user
	_, err := theService.oauthService.CreateUser(username, password)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	// Redirect to the login page
	http.Redirect(w, r, "/web/login", http.StatusFound)
}

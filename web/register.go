package web

import "net/http"

func registerForm(w http.ResponseWriter, r *http.Request) {
	// Initialise a new session service
	sessionService := newSessionService(s.cnf)
	if err := sessionService.initSession(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	renderTemplate(w, "register.tmpl", map[string]interface{}{
		"error": sessionService.getFlashMessage(r, w),
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	// Initialise a new session service
	sessionService := newSessionService(s.cnf)
	if err := sessionService.initSession(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the submitted form data
	r.ParseForm()
	username := r.Form["email"][0]
	password := r.Form["password"][0]

	// Check that the submitted email hasn't been registered already
	if s.oauthService.UserExists(username) {
		sessionService.setFlashMessage("Email already taken", r, w)
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	// Create a user
	_, err := s.oauthService.CreateUser(username, password)
	if err != nil {
		sessionService.setFlashMessage(err.Error(), r, w)
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	// Redirect to the login page
	http.Redirect(w, r, "/web/login", http.StatusFound)
}

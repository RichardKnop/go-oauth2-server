package web

import "net/http"

func registerForm(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	renderTemplate(w, "register.html", map[string]interface{}{
		"error":       sessionService.GetFlashMessage(),
		"queryString": getQueryString(r.URL.Query()),
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
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
	_, err = theService.oauthService.CreateUser(
		r.Form.Get("email"),    // username
		r.Form.Get("password"), // password
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Redirect to the login page
	redirectWithQueryString("/web/login", r.URL.Query(), w, r)
}

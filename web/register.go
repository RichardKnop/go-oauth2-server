package web

import "net/http"

func registerForm(w http.ResponseWriter, r *http.Request) {
	sessionService := noLoginRequired(w, r)
	if sessionService == nil {
		return
	}

	// Render the template
	renderTemplate(w, "register.tmpl", map[string]interface{}{
		"error": sessionService.GetFlashMessage(),
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	sessionService := noLoginRequired(w, r)
	if sessionService == nil {
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
	redirectAndKeepQueryString("/web/login", w, r)
}

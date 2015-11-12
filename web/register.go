package web

import "net/http"

func registerForm(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "register.tmpl", map[string]interface{}{
		"error": getLastFlashMessage(session, r, w),
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]

	if oauthService.UserExists(username) {
		addFlashMessage(session, r, w, "Username already taken")
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	_, err = oauthService.CreateUser(username, password)

	if err != nil {
		addFlashMessage(session, r, w, err.Error())
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/web/login", http.StatusFound)
}

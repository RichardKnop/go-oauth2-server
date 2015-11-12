package web

import (
	"log"
	"net/http"
)

func loginForm(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "login.tmpl", map[string]interface{}{
		"error": getLastFlashMessage(session, r, w),
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	username := r.Form["email"][0]
	password := r.Form["password"][0]

	_, err = oauthService.AuthUser(username, password)
	log.Print(err)

	if err != nil {
		session.AddFlash(err.Error())
		session.Save(r, w)
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	// client := &oauth.Client{} // TODO
	// return s.oauthService.GrantAccessToken(client, user, "read_write")
	// accessToken, refreshToken, err := accounts.GetService().Login(
	// 	r.Form["username"][0],
	// 	r.Form["password"][0],
	// )
	//
	// log.Print(accessToken)
	// log.Print(refreshToken)
	// log.Print(err)
	// // TODO
}

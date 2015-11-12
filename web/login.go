package web

import (
	"html/template"
	"log"
	"net/http"
)

func loginForm(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{}
	if flashes := session.Flashes(); len(flashes) > 0 {
		data["error"] = flashes[0]
		session.Save(r, w)
	}

	tmpl, _ := template.ParseFiles("web/templates/login.html.tmpl")
	tmpl.Execute(w, data)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	username := r.Form["username"][0]
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

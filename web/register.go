package web

import (
	"html/template"
	"log"
	"net/http"
)

func registerForm(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "areatech")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{}
	if flashes := session.Flashes(); len(flashes) > 0 {
		log.Print(flashes)
		data["error"] = flashes[0]
	}

	tmpl, _ := template.ParseFiles("web/templates/register.html.tmpl")
	tmpl.Execute(w, data)
}

func register(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "areatech")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]

	if oauthService.UserExists(username) {
		log.Print("Username already taken")
		session.AddFlash("Username already taken")
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	_, err = oauthService.CreateUser(username, password)
	log.Print(err)

	if err != nil {
		session.AddFlash(err.Error())
		http.Redirect(w, r, "/web/register", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/web/login", http.StatusFound)
}

package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/accounts"
)

func login(w http.ResponseWriter, r *http.Request) {
	accountsService := accounts.GetService()
	data := map[string]interface{}{}

	if r.Method == "POST" {
		r.ParseForm()
		user, err := accountsService.Login(r.Form["username"][0], r.Form["password"][0])
		if err != nil {
			data["error"] = err.Error()
			renderLogin(w, r, data)
			return
		}

		log.Print(user)
		// probably redirect to a page specified in the query string
	}

	renderLogin(w, r, data)
}

func renderLogin(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	tmpl, _ := template.ParseFiles("web/templates/login.html.tmpl")
	tmpl.Execute(w, data)
}

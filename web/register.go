package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/accounts"
)

func register(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	accountsService := accounts.GetService()

	if r.Method == "POST" {
		r.ParseForm()
		user, err := accountsService.Register(r.Form["username"][0], r.Form["password"][0])

		if err != nil {
			data["error"] = err.Error()
			renderRegister(w, r, data)
			return
		}

		log.Print(user)
		// TODO - probably redirect to the login page
	}

	renderRegister(w, r, data)
}

func renderRegister(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	tmpl, _ := template.ParseFiles("web/templates/register.html.tmpl")
	tmpl.Execute(w, data)
}

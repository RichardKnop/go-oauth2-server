package web

import (
	"html/template"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == "POST" {
		r.ParseForm()
		// TODO
	}

	renderRegister(w, r, data)
}

func renderRegister(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	tmpl, _ := template.ParseFiles("web/templates/register.html.tmpl")
	tmpl.Execute(w, data)
}

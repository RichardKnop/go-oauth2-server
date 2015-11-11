package web

import (
	"net/http"
	"text/template"
)

func authorize(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == "POST" {
		r.ParseForm()
		// TODO
	}

	renderAuthorize(w, r, data)
}

func renderAuthorize(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	tmpl, _ := template.ParseFiles("web/templates/authorize.html.tmpl")
	tmpl.Execute(w, data)
}

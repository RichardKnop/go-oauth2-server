package web

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func getSession(r *http.Request) (*sessions.Session, error) {
	s, err := sessionStore.Get(r, "areatech")
	if err != nil {
		return nil, err
	}
	s.Options = sessionOptions
	return s, nil
}

func addFlashMessage(s *sessions.Session, r *http.Request, w http.ResponseWriter, msg string) {
	s.AddFlash(msg)
	s.Save(r, w)
}

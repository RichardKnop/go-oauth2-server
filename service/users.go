package service

import (
	"net/http"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
	"github.com/ant0ine/go-json-rest/rest"
	"golang.org/x/crypto/bcrypt"
)

// UsersHandler ...
func UsersHandler(w rest.ResponseWriter, r *rest.Request) {
	clientID, clientPassword, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	cnf := config.NewConfig()

	db, err := database.NewDatabase(cnf)
	if err != nil {
		rest.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}

	client := Client{}
	if db.Where("client_id = ?", clientID).First(&client).RecordNotFound() {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(clientPassword)); err != nil {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	// TODO
}

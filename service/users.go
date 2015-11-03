package service

import (
	"net/http"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
	"github.com/ant0ine/go-json-rest/rest"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser - registers a new user
func RegisterUser(w rest.ResponseWriter, r *rest.Request) {
	user := User{}
	if err := r.DecodeJsonPayload(&user); err != nil {
		rest.Error(w, "Decode JSON error", http.StatusBadRequest)
		return
	}

	cnf := config.NewConfig()

	db, err := database.NewDatabase(cnf)
	if err != nil {
		rest.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}

	if db.Where("username = ?", user.Username).First(&User{}).RecordNotFound() {
		rest.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 3)
	if err != nil {
		rest.Error(w, "Bcrypt error", http.StatusInternalServerError)
		return
	}

	user.Password = string(hash)
	if err := db.Create(&user).Error; err != nil {
		rest.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.WriteJson(&user)
}

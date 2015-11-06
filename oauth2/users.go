package oauth2

import (
	"fmt"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Registers a new user
func register(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	user := User{}
	if err := r.DecodeJsonPayload(&user); err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Username are case insensitive
	if db.Where("LOWER(username) = LOWER(?)", user.Username).First(&User{}).RowsAffected > 0 {
		api.Error(w, fmt.Sprintf("%s already taken", user.Username), http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 3)
	if err != nil {
		api.Error(w, "Bcrypt error", http.StatusInternalServerError)
		return
	}

	user.Password = string(passwordHash)
	if err := db.Create(&user).Error; err != nil {
		api.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteJson(map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
}

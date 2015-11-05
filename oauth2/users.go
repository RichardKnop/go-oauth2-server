package oauth2

import (
	"fmt"
	"net/http"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Registers a new user
func register(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	requiredFields := []string{"username", "password", "first_name", "last_name"}
	for _, requiredField := range requiredFields {
		if r.FormValue(requiredField) == "" {
			rest.Error(w, fmt.Sprintf("%s required", requiredField), http.StatusBadRequest)
			return
		}
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")

	// Case insensitive search, usernames will probably be emails and
	// foo@bar.com is identical to FOO@BAR.com
	if db.Where("LOWER(username) = LOWER(?)", username).First(&User{}).RowsAffected > 0 {
		rest.Error(w, fmt.Sprintf("%s already taken", username), http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 3)
	if err != nil {
		rest.Error(w, "Bcrypt error", http.StatusInternalServerError)
		return
	}

	user := User{
		Username:  username,
		Password:  string(passwordHash),
		FirstName: firstName,
		LastName:  lastName,
	}
	if err := db.Create(&user).Error; err != nil {
		rest.Error(w, "Error saving user", http.StatusInternalServerError)
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

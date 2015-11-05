package oauth2

import (
	"net/http"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Registers a new user
func registerUser(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	user := User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	// Case insensitive search, usernames will probably be emails and
	// foo@bar.com is identical to FOO@BAR.com
	if db.Where("LOWER(username) = LOWER(?)", user.Username).First(&User{}).RowsAffected > 0 {
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteJson(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
}

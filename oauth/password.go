package oauth

import (
	"golang.org/x/crypto/bcrypt"
)

// verifyPassword compares password and the hashed password
func verifyPassword(passwordHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

// hashPassword creates a bcrypt password hash
func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 3)
}

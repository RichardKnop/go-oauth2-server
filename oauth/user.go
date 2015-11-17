package oauth

import (
	"errors"

	"github.com/RichardKnop/go-oauth2-server/password"
)

// UserExists returns true if user exists
func (s *Service) UserExists(username string) bool {
	_, err := s.FindUserByUsername(username)
	return err == nil
}

// AuthUser authenticates user
func (s *Service) AuthUser(username, thePassword string) (*User, error) {
	// Fetch the user
	user, err := s.FindUserByUsername(username)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// Verify the password
	if password.VerifyPassword(user.Password, thePassword) != nil {
		return nil, errors.New("Invalid password")
	}

	return user, nil
}

// CreateUser saves a new user to database
func (s *Service) CreateUser(username, thePassword string) (*User, error) {
	passwordHash, err := password.HashPassword(thePassword)
	if err != nil {
		return nil, errors.New("Bcrypt error")
	}
	user := User{
		Username: username,
		Password: string(passwordHash),
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("Error saving user to database")
	}
	return &user, nil
}

// FindUserByUsername looks up a user by username
func (s *Service) FindUserByUsername(username string) (*User, error) {
	// Usernames are case insensitive
	user := new(User)
	if s.db.Where("LOWER(username) = LOWER(?)", username).First(user).RecordNotFound() {
		return nil, errors.New("User not found")
	}
	return user, nil
}

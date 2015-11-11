package oauth

import (
	"errors"
)

// UserExists returns true if user exists
func (s *Service) UserExists(username string) bool {
	_, err := s.findUserByUsername(username)
	return err == nil
}

// AuthUser authenticates user
func (s *Service) AuthUser(username, password string) (*User, error) {
	// Fetch the user
	user, err := s.findUserByUsername(username)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// Verify the password
	if verifyPassword(user.Password, password) != nil {
		return nil, errors.New("Invalid password")
	}

	return user, nil
}

// CreateUser saves a new user to database
func (s *Service) CreateUser(username, password string) (*User, error) {
	passwordHash, err := hashPassword(password)
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

func (s *Service) findUserByUsername(username string) (*User, error) {
	// Usernames are case insensitive
	user := new(User)
	if s.db.Where("LOWER(username) = LOWER(?)", username).First(user).RecordNotFound() {
		return nil, errors.New("User not found")
	}
	return user, nil
}

package oauth

import (
	"errors"
	"fmt"
	"time"

	pass "github.com/RichardKnop/go-oauth2-server/password"
	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/jinzhu/gorm"
)

var (
	// MinPasswordLength defines minimum password length
	MinPasswordLength = 6

	// ErrPasswordTooShort ...
	ErrPasswordTooShort = fmt.Errorf(
		"Password must be at least %d characters long",
		MinPasswordLength,
	)
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("User not found")
	// ErrInvalidUserPassword ...
	ErrInvalidUserPassword = errors.New("Invalid user password")
	// ErrUserPasswordNotSet ...
	ErrUserPasswordNotSet = errors.New("User password not set")
	// ErrUsernameTaken ...
	ErrUsernameTaken = errors.New("Username taken")
)

// UserExists returns true if user exists
func (s *Service) UserExists(username string) bool {
	_, err := s.FindUserByUsername(username)
	return err == nil
}

// FindUserByUsername looks up a user by username
func (s *Service) FindUserByUsername(username string) (*User, error) {
	// Usernames are case insensitive
	user := new(User)
	notFound := s.db.Where("LOWER(username) = LOWER(?)", username).
		First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// CreateUser saves a new user to database
func (s *Service) CreateUser(username, password string) (*User, error) {
	return s.createUserCommon(s.db, username, password)
}

// CreateUserTx saves a new user to database using injected db object
func (s *Service) CreateUserTx(tx *gorm.DB, username, password string) (*User, error) {
	return s.createUserCommon(tx, username, password)
}

// SetPassword sets a user password
func (s *Service) SetPassword(user *User, password string) error {
	return s.setPasswordCommon(s.db, user, password)
}

// SetPasswordTx sets a user password in a transaction
func (s *Service) SetPasswordTx(tx *gorm.DB, user *User, password string) error {
	return s.setPasswordCommon(tx, user, password)
}

// AuthUser authenticates user
func (s *Service) AuthUser(username, password string) (*User, error) {
	// Fetch the user
	user, err := s.FindUserByUsername(username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check that the password is set
	if !user.Password.Valid {
		return nil, ErrUserPasswordNotSet
	}

	// Verify the password
	if pass.VerifyPassword(user.Password.String, password) != nil {
		return nil, ErrInvalidUserPassword
	}

	return user, nil
}

func (s *Service) createUserCommon(db *gorm.DB, username, password string) (*User, error) {
	// Start with a user without a password
	user := &User{
		Username: username,
		Password: util.StringOrNull(""),
	}

	// If the password is being set already, create a bcrypt hash
	if password != "" {
		if len(password) < MinPasswordLength {
			return nil, ErrPasswordTooShort
		}
		passwordHash, err := pass.HashPassword(password)
		if err != nil {
			return nil, err
		}
		user.Password = util.StringOrNull(string(passwordHash))
	}

	// Check the username is available
	if s.UserExists(user.Username) {
		return nil, ErrUsernameTaken
	}

	// Create the user
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) setPasswordCommon(db *gorm.DB, user *User, password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	// Create a bcrypt hash
	passwordHash, err := pass.HashPassword(password)
	if err != nil {
		return err
	}

	// Set the password on the user object
	return db.Model(user).UpdateColumns(User{
		Password: util.StringOrNull(string(passwordHash)),
		Model:    gorm.Model{UpdatedAt: time.Now()},
	}).Error
}

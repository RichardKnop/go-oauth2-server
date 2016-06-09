package oauth

import (
	"errors"
	"time"

	pass "github.com/RichardKnop/go-oauth2-server/password"
	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("User not found")
	// ErrInvalidUserPassword ...
	ErrInvalidUserPassword = errors.New("Invalid user password")
	// ErrCannotSetEmptyUserPassword ...
	ErrCannotSetEmptyUserPassword = errors.New("Cannot set empty user password")
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

// SetPassword saves a new user to database
func (s *Service) SetPassword(user *User, password string) error {
	// Cannot set password to empty
	if password == "" {
		return ErrCannotSetEmptyUserPassword
	}

	// Create a bcrypt hash
	passwordHash, err := pass.HashPassword(password)
	if err != nil {
		return err
	}

	// Set the password on the user object
	if err := s.db.Model(user).UpdateColumns(User{
		Password: util.StringOrNull(string(passwordHash)),
		Model:    gorm.Model{UpdatedAt: time.Now()},
	}).Error; err != nil {
		return err
	}

	return nil
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

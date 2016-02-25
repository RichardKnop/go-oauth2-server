package oauth

import (
	"errors"

	pass "github.com/RichardKnop/go-oauth2-server/password"
	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/jinzhu/gorm"
)

var (
	errUserNotFound               = errors.New("User not found")
	errInvalidUserPassword        = errors.New("Invalid user password")
	errCannotSetEmptyUserPassword = errors.New("Cannot set empty user password")
	errUserPasswordNotSet         = errors.New("User password not set")
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
		return nil, errUserNotFound
	}

	return user, nil
}

// CreateUser saves a new user to database
func (s *Service) CreateUser(username, password string) (*User, error) {
	return createUserCommon(s.db, username, password)
}

// CreateUserTx saves a new user to database using injected db object
func (s *Service) CreateUserTx(tx *gorm.DB, username, password string) (*User, error) {
	return createUserCommon(tx, username, password)
}

// SetPassword saves a new user to database
func (s *Service) SetPassword(user *User, password string) error {
	// Cannot set password to empty
	if password == "" {
		return errCannotSetEmptyUserPassword
	}

	// Create a bcrypt hash
	passwordHash, err := pass.HashPassword(password)
	if err != nil {
		return err
	}

	// Set the password on the user object
	if err := s.db.Model(user).UpdateColumn(
		"password",
		string(passwordHash),
	).Error; err != nil {
		return err
	}

	return nil
}

// AuthUser authenticates user
func (s *Service) AuthUser(username, password string) (*User, error) {
	// Fetch the user
	user, err := s.FindUserByUsername(username)
	if err != nil {
		return nil, errUserNotFound
	}

	// Check that the password is set
	if !user.Password.Valid {
		return nil, errUserPasswordNotSet
	}

	// Verify the password
	if pass.VerifyPassword(user.Password.String, password) != nil {
		return nil, errInvalidUserPassword
	}

	return user, nil
}

func createUserCommon(db *gorm.DB, username, password string) (*User, error) {
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

	// Create the user
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

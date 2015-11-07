package oauth2

import "errors"

func validateUserData(user *User) error {
	if user.Username == "" {
		return errors.New("username required")
	}

	if user.Password == "" {
		return errors.New("password required")
	}

	return nil
}

package oauth2

import "errors"

func validateUserData(user *User) error {
	if user.Username == "" {
		return errors.New("username required")
	}

	if user.Password == "" {
		return errors.New("password required")
	}

	if user.FirstName == "" {
		return errors.New("first_name required")
	}

	if user.LastName == "" {
		return errors.New("last_name required")
	}

	return nil
}

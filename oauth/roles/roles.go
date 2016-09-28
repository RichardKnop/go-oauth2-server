package roles

import (
	"errors"
)

const (
	// Superuser ...
	Superuser = "superuser"
	// User ...
	User = "user"
)

var roleWeight = map[string]int{
	Superuser: 100,
	User:      1,
}

// IsGreaterThan returns true if role1 is greater than role2
func IsGreaterThan(role1, role2 string) (bool, error) {
	// Get weight of the first role
	weight1, ok := roleWeight[role1]
	if !ok {
		return false, errors.New("Role weight not found")
	}

	// Get weight of the second role
	weight2, ok := roleWeight[role2]
	if !ok {
		return false, errors.New("Role weight not found")
	}

	return weight1 > weight2, nil
}

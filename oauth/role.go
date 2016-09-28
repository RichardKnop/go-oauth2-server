package oauth

import (
	"errors"
)

var (
	// ErrRoleNotFound ...
	ErrRoleNotFound = errors.New("Role not found")
)

// FindRoleByID looks up a role by ID and returns it
func (s *Service) FindRoleByID(id string) (*Role, error) {
	role := new(Role)
	if s.db.Where("id = ?", id).First(role).RecordNotFound() {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

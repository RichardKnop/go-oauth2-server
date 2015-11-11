package oauth

import (
	"errors"
	"strings"
)

func (s *Service) getScope(requestedScope string) (string, error) {
	// Return the default scope if the requested scope is empty
	if requestedScope == "" {
		return s.getDefaultScope(), nil
	}

	// If the requested scope exists in the database, return it
	if s.scopeExists(requestedScope) {
		return requestedScope, nil
	}

	// Otherwise return error
	return "", errors.New("Invalid scope")
}

func (s *Service) getDefaultScope() string {
	// Fetch default scopes
	var scopes []string
	s.db.Model(new(Scope)).Where(Scope{IsDefault: true}).Pluck("scope", &scopes)

	// Return space delimited scope string
	return strings.Join(scopes, " ")
}

func (s *Service) scopeExists(requestedScope string) bool {
	// Split the requested scope string
	scopes := strings.Split(requestedScope, " ")

	// Count how many of requested scopes exist in the database
	var count int
	s.db.Model(new(Scope)).Where("scope in (?)", scopes).Count(&count)

	// Return true only if all requested scopes found
	return count == len(scopes)
}

func (s *Service) scopeNotGreater(newScope, oldScope string) bool {
	// Empty scope is never greater
	if newScope == "" {
		return true
	}

	// Split the old scope string
	oldScopes := strings.Split(oldScope, " ")

	// Iterate over new scopes
	for _, newScope := range strings.Split(newScope, " ") {
		// If the new scope was not part of the old scope string, return false
		if !stringInSlice(newScope, oldScopes) {
			return false
		}
	}

	// The new scope is the same or more restrictive than the old one, return true
	return true
}

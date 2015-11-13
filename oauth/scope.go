package oauth

import (
	"errors"
	"sort"
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
	s.db.Model(new(Scope)).Where(Scope{
		IsDefault: true,
	}).Pluck("scope", &scopes)

	// Sort the scopes alphabetically
	sort.Strings(scopes)

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

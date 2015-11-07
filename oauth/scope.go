package oauth

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

func getScope(db *gorm.DB, requestedScope string) (string, error) {
	// Return the default scope if the requested scope is empty
	if requestedScope == "" {
		return getDefaultScope(db), nil
	}

	// If the requested scope exists in the database, return it
	if scopeExists(db, requestedScope) {
		return requestedScope, nil
	}

	// Otherwise return error
	return "", errors.New("Invalid scope")
}

func getDefaultScope(db *gorm.DB) string {
	// Fetch default scopes
	var scopes []string
	db.Model(&Scope{}).Where(&Scope{IsDefault: true}).Pluck("scope", &scopes)

	// Return space delimited scope string
	return strings.Join(scopes, " ")
}

func scopeExists(db *gorm.DB, requestedScope string) bool {
	// Split the requested scope string
	scopes := strings.Split(requestedScope, " ")
	// Count how many of requested scopes exist in the database
	var count int
	db.Model(&Scope{}).Where("scope in (?)", scopes).Count(&count)

	// Return true only if all requested scopes found
	return count == len(scopes)
}

func scopeNotGreater(newScope, oldScope string) bool {
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

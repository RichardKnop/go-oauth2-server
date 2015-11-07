package oauth

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

// If requested scope is empty, returns default scope
// Otherwise validates the scope and returns error if invalid
func getScope(db *gorm.DB, requestedScope string) (string, error) {
	// Requested scope emtpy, so just get default scope
	if requestedScope == "" {
		var scopes []string
		db.Model(&Scope{}).Where(&Scope{IsDefault: true}).Pluck("scope", &scopes)
		return strings.Join(scopes, " "), nil
	}

	// All seems fine, return requested scope
	if scopeExists(db, requestedScope) {
		return requestedScope, nil
	}

	return "", errors.New("Invalid scope")
}

// Returns error if scope not fully defined in our database
func scopeExists(db *gorm.DB, scope string) bool {
	scopes := strings.Split(scope, " ")
	var count int
	db.Model(&Scope{}).Where("scope in (?)", scopes).Count(&count)
	return count == len(scopes)
}

// Returns true if a new scope DOES NOT include anything not in the old one
func scopeNotGreater(newScope, oldScope string) bool {
	newScopes := strings.Split(newScope, " ")
	oldScopes := strings.Split(oldScope, " ")

	for _, newScope := range newScopes {
		if !stringInSlice(newScope, oldScopes) {
			return false
		}
	}

	return true
}

// Helpful function similar to "x in y" Python construct
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

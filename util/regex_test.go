package util_test

import (
	"testing"

	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/stretchr/testify/assert"
)

func TestRegexExtractMatch(t *testing.T) {
	match, err := util.RegexExtractMatch(
		"...FOO...",
		"^...(?P<the_name>[A-Z]{3})...$",
		"the_name",
	)

	assert.Nil(t, err)
	assert.Equal(t, "FOO", match)
}

func TestRegexExtractMatches(t *testing.T) {
	matches, err := util.RegexExtractMatches(
		"HKDJPY",
		"^(?P<from_currency>[A-Z]{3})(?P<to_currency>[A-Z]{3})$",
		"from_currency",
		"to_currency",
	)

	assert.Nil(t, err)
	assert.Equal(t, "HKD", matches["from_currency"])
	assert.Equal(t, "JPY", matches["to_currency"])
}

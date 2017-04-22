package util_test

import (
	"testing"

	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, util.StringInSlice("a", []string{"a", "b", "c"}))

	assert.False(t, util.StringInSlice("d", []string{"a", "b", "c"}))
}

func TestSpaceDelimitedStringNotGreater(t *testing.T) {
	assert.True(t, util.SpaceDelimitedStringNotGreater("", "bar foo qux"))

	assert.True(t, util.SpaceDelimitedStringNotGreater("foo", "bar foo qux"))

	assert.True(t, util.SpaceDelimitedStringNotGreater("bar foo qux", "foo bar qux"))

	assert.False(t, util.SpaceDelimitedStringNotGreater("foo bar qux bogus", "bar foo qux"))
}

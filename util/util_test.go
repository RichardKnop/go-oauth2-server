package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice("a", []string{"a", "b", "c"}))

	assert.False(t, StringInSlice("d", []string{"a", "b", "c"}))
}

func TestSpaceDelimitedStringNotGreater(t *testing.T) {
	assert.True(t, SpaceDelimitedStringNotGreater("", "bar foo qux"))

	assert.True(t, SpaceDelimitedStringNotGreater("foo", "bar foo qux"))

	assert.True(t, SpaceDelimitedStringNotGreater("bar foo qux", "foo bar qux"))

	assert.False(t, SpaceDelimitedStringNotGreater("foo bar qux bogus", "bar foo qux"))
}

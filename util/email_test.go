package util_test

import (
	"testing"

	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	assert.False(t, util.ValidateEmail("test@user"))
	assert.True(t, util.ValidateEmail("test@user.com"))
}

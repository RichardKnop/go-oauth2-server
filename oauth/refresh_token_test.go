package oauth

// import (
// 	"log"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// func (suite *OauthTestSuite) TestGetOrCreateRefreshToken() {
// 	var refreshToken *RefreshToken
// 	var err error
//
// 	// Without user
// 	refreshToken, err = suite.service.getOrCreateRefreshToken(
// 		suite.client,
// 		nil,
// 		"foo bar",
// 	)
// 	assert.Nil(suite.T(), err)
//
// 	// With User
// 	refreshToken, err = suite.service.getOrCreateRefreshToken(
// 		suite.client,
// 		suite.user,
// 		"foo bar",
// 	)
// 	assert.Nil(suite.T(), err)
//
// 	log.Print(refreshToken)
// }

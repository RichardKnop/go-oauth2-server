package accounts

import (
	"github.com/RichardKnop/go-oauth2-server/oauth"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetOauthService() oauth.ServiceInterface
}

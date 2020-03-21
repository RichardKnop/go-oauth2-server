package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/RichardKnop/go-oauth2-server/models"
	"strings"
	"time"
)

var (
	// ErrAuthorizationCodeNotFound ...
	ErrAuthorizationCodeNotFound = errors.New("Authorization code not found")
	// ErrAuthorizationCodeExpired ...
	ErrAuthorizationCodeExpired = errors.New("Authorization code expired")
	// ErrPKCEMissingChallenge ...
	ErrPKCEMissingChallenge = errors.New("PKCE: code_challenge missing")
	// ErrPKCEMissingVerifier ...
	ErrPKCEMissingVerifier = errors.New("PKCE: code_verifier missing")
	// ErrPKCENoMatch ...
	ErrPKCENoMatch = errors.New("PKCE: failed challenge")
	// ErrPKCENotAllowed ...
	ErrPKCENotAllowed = errors.New("PKCE: not enabled for this client")
)

// CodeChallengeMethodS256 verifies an S256 challenge against a verifier
func (s *Service) CodeChallengeMethodS256(verifier string, challenge string) bool {
	hasher := sha256.New()
	hasher.Write([]byte(verifier))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return strings.TrimRight(sha, "=") == challenge
}

// CodeChallengeMethodPlain verifies an PLAIN challenge against a verifier
func (s *Service) CodeChallengeMethodPlain(verifier string, challenge string) bool {
	return verifier == challenge
}

// GrantAuthorizationCode grants a new authorization code
func (s *Service) GrantAuthorizationCode(client *models.OauthClient, user *models.OauthUser, expiresIn int, redirectURI, scope string, codeChallenge string) (*models.OauthAuthorizationCode, error) {
	if client.Public && codeChallenge == "" {
		return nil, ErrPKCEMissingChallenge
	}
	if !client.Public && len(codeChallenge) > 0 {
		return nil, ErrPKCENotAllowed
	}

	// Create a new authorization code
	authorizationCode := models.NewOauthAuthorizationCode(client, user, expiresIn, redirectURI, scope, codeChallenge)
	if err := s.db.Create(authorizationCode).Error; err != nil {
		return nil, err
	}
	authorizationCode.Client = client
	authorizationCode.User = user

	return authorizationCode, nil
}

// getValidAuthorizationCode returns a valid non expired authorization code
func (s *Service) getValidAuthorizationCode(code, redirectURI string, client *models.OauthClient, codeVerifier string, codeChallengeMethod string) (*models.OauthAuthorizationCode, error) {
	// Map of grant types against handler functions
	verifiers := map[string]func(verifier string, challenge string) bool{
		"plain": s.CodeChallengeMethodPlain,
		"S256":  s.CodeChallengeMethodS256,
	}

	// Fetch the auth code from the database
	authorizationCode := new(models.OauthAuthorizationCode)
	notFound := models.OauthAuthorizationCodePreload(s.db).Where("client_id = ?", client.ID).
		Where("code = ?", code).First(authorizationCode).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAuthorizationCodeNotFound
	}

	if client.Public {
		if codeVerifier == "" {
			return nil, ErrPKCEMissingVerifier
		}

		if codeChallengeMethod == "" {
			codeChallengeMethod = "S256"
		}

		if !verifiers[codeChallengeMethod](codeVerifier, authorizationCode.CodeChallenge) {
			return nil, ErrPKCENoMatch
		}

	}

	// Redirect URI must match if it was used to obtain the authorization code
	if redirectURI != authorizationCode.RedirectURI.String {
		return nil, ErrInvalidRedirectURI
	}

	// Check the authorization code hasn't expired
	if time.Now().After(authorizationCode.ExpiresAt) {
		return nil, ErrAuthorizationCodeExpired
	}

	return authorizationCode, nil
}

package authx

import "time"

type TokenClaims struct {
	UserId          string
	Email           string
	Roles           []string
	PasswordVersion int
}

type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type TokenGenerator interface {
	GeneratePair(claims TokenClaims) (*TokenPair, error)
	Verify(tokenString string) (*TokenClaims, error)
}

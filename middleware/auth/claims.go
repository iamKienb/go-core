package authx

import "time"

type Claims struct {
	UserId          string
	Email           string
	Roles           []string
	PasswordVersion int
}

type Pair struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type Generator interface {
	GeneratePair(claims Claims) (*Pair, error)
	Verify(tokenString string) (*Claims, error)
}

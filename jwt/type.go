package jwtx

import (
	"time"
)

type Claims struct {
	UserID          string
	Email           string
	FullName        string
	Roles           []string
	PasswordVersion int
}

type Pair struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

type JWTXService interface {
	GeneratePair(claims Claims) (*Pair, error)
	Verify(tokenString string) (*Claims, error)
}

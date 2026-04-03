package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId string
	Role   string
	Perms  []string
	jwt.RegisteredClaims
}

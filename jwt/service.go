package jwtx

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type jwtClaims struct {
	UserID          string   `json:"uid"`
	Email           string   `json:"email"`
	Roles           []string `json:"roles"`
	PasswordVersion int      `json:"pwd_v"`
	jwt.RegisteredClaims
}

func (x *JWTX) GeneratePair(claims Claims) (*Pair, error) {
	now := time.Now().UTC()
	accessExpAt := now.Add(x.cfg.AccessExpiry)
	refreshExpAt := now.Add(x.cfg.RefreshExpiry)

	accessToken, err := x.Sign(claims, accessExpAt)
	if err != nil {
		return nil, fmt.Errorf("jwt: sign access token: %w", err)
	}

	refreshToken, err := x.Sign(claims, refreshExpAt)
	if err != nil {
		return nil, fmt.Errorf("jwt: sign refresh token: %w", err)
	}

	return &Pair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpAt,
		RefreshTokenExpiresAt: refreshExpAt,
	}, nil
}

func (x *JWTX) Sign(claims Claims, expiryAt time.Time) (string, error) {
	if x == nil || x.privateKey == nil {
		return "", fmt.Errorf("jwt: private key is not configured")
	}

	claim := jwtClaims{
		UserID:          claims.UserID,
		Email:           claims.Email,
		Roles:           claims.Roles,
		PasswordVersion: claims.PasswordVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "user-command",
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claim).SignedString(x.privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (x *JWTX) Verify(tokenString string) (*Claims, error) {
	if x == nil || x.publicKey == nil {
		return nil, fmt.Errorf("jwt: public key is not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return x.publicKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("jwt: token is expired")
		}
		return nil, fmt.Errorf("jwt: parse failed: %w", err)
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &Claims{
			UserID:          claims.UserID,
			Email:           claims.Email,
			Roles:           claims.Roles,
			PasswordVersion: claims.PasswordVersion,
		}, nil
	}

	return nil, fmt.Errorf("jwt: invalid token")
}

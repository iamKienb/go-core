package auth

import (
	"context"
	"time"
)

type contextKey string

const userHeaderKey contextKey = "user_info"

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

func GetUserInfoFromCtx(ctx context.Context) *TokenClaims {
	if claims, ok := ctx.Value(userHeaderKey).(*TokenClaims); ok {
		return claims
	}

	return nil
}

func SetUserInfoToCtx(ctx context.Context, claims *TokenClaims) context.Context {
	return context.WithValue(ctx, userHeaderKey, claims)
}

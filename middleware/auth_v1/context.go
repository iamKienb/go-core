package auth_v1

import (
	"context"

	jwtx "github.com/iamKienb/go-core/jwt"
)

const (
	HeaderRequestID = "X-Request-ID"
	HeaderTraceID   = "X-Trace-ID"
	HeaderUserID    = "X-User-ID"
	HeaderUserRole  = "X-User-Role"
	HeaderUserEmail = "X-User-Email"
	HeaderUserName  = "X-User-Name"
)

type contextKey string

const (
	userHeaderKey contextKey = "user_info"
)

func GetUserInfoFromCtx(ctx context.Context) *jwtx.Claims {
	if claims, ok := ctx.Value(userHeaderKey).(*jwtx.Claims); ok {
		return claims
	}

	return nil
}

func SetUserInfoToCtx(ctx context.Context, claims *jwtx.Claims) context.Context {
	return context.WithValue(ctx, userHeaderKey, claims)
}

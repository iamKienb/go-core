package authx

import (
	"context"

	jwtx "github.com/iamKienb/go-core/jwt"
)

const (
	HeaderRequestID = "X-Request-ID"
	HeaderUserID    = "X-User-ID"
	HeaderUserRole  = "X-User-Role"
	HeaderUserEmail = "X-User-Email"
	HeaderUserName  = "X-User-Name"
)

type contextKey string

const (
	userHeaderKey      contextKey = "user_info"
	requestIDHeaderKey contextKey = "request_id"
)

func GetUserInfoFromCtx(ctx context.Context) *jwtx.Claims {
	if claims, ok := ctx.Value(userHeaderKey).(*jwtx.Claims); ok {
		return claims
	}

	return nil
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDHeaderKey).(string); ok {
		return id
	}

	return ""
}

func SetUserInfoToCtx(ctx context.Context, claims *jwtx.Claims) context.Context {
	return context.WithValue(ctx, userHeaderKey, claims)
}

func SetRequestIDToCtx(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDHeaderKey, reqID)
}

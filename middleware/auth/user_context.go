package authx

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	jwtx "github.com/iamKienb/go-core/jwt"
)

const (
	HeaderUserID    = "X-User-ID"
	HeaderUserRole  = "X-User-Role"
	HeaderUserEmail = "X-User-Email"
	HeaderUserName  = "X-User-Name"
)

type contextKey string

const userContextKey contextKey = "user_info_ctx_key"

func ExtractUserContextInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			userID := req.Header().Get(HeaderUserID)
			email := req.Header().Get(HeaderUserEmail)
			fullName := req.Header().Get(HeaderUserName)

			var roles []string
			if rawRoles := req.Header().Get(HeaderUserRole); rawRoles != "" {
				roles = strings.Split(rawRoles, ",")
			}

			if userID != "" || email != "" || len(roles) > 0 {
				claims := &jwtx.Claims{
					UserID:   userID,
					Email:    email,
					FullName: fullName,
					Roles:    roles,
				}
				ctx = context.WithValue(ctx, userContextKey, claims)
			}

			return next(ctx, req)
		}
	}
}

func GetUserInfoFromCtx(ctx context.Context) *jwtx.Claims {
	if claims, ok := ctx.Value(userContextKey).(*jwtx.Claims); ok {
		return claims
	}
	return nil
}

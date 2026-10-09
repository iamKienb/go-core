package auth_v1

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	jwtx "github.com/iamKienb/go-core/jwt"
)

func AuthInternalInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			var roles []string
			if rawRoles := req.Header().Get(HeaderUserRole); rawRoles != "" {
				roles = strings.Split(rawRoles, ",")
			}

			userID := req.Header().Get(HeaderUserID)
			email := req.Header().Get(HeaderUserEmail)
			fullName := req.Header().Get(HeaderUserName)
			if userID != "" || email != "" || len(roles) > 0 {
				claims := &jwtx.Claims{
					UserID:   userID,
					Email:    email,
					FullName: fullName,
					Roles:    roles,
				}
				ctx = SetUserInfoToCtx(ctx, claims)
			}

			return next(ctx, req)
		}
	}
}

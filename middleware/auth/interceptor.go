package auth

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

func AuthInternalInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			reqID := req.Header().Get(HeaderRequestID)
			if reqID == "" {
				reqID = "internal-" + uuid.NewString()
			}
			ctx = SetRequestIDToCtx(ctx, reqID)

			claims := &TokenClaims{
				UserId: req.Header().Get(HeaderUserID),
				Email:  req.Header().Get(HeaderUserEmail),
				Roles:  strings.Split(req.Header().Get(HeaderUserRole), ","),
			}
			ctx = SetUserInfoToCtx(ctx, claims)

			return next(ctx, req)
		}
	}
}

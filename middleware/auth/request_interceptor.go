package authx

import (
	"context"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

func RequestContextInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			reqID := req.Header().Get(HeaderRequestID)
			if reqID == "" {
				reqID = "req-" + uuid.NewString()
			}

			req.Header().Set(HeaderRequestID, reqID)
			ctx = SetRequestIDToCtx(ctx, reqID)

			return next(ctx, req)
		}
	}
}

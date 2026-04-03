package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	app_error "github.com/iamKienb/shopify-go-platform/error"
	"github.com/iamKienb/shopify-go-platform/response"
	"github.com/iamKienb/shopify-go-platform/utils"
	"go.uber.org/zap"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetHeader(utils.HeaderUserId)
		role := c.GetHeader(utils.HeaderUserRole)

		if uid == "" {
			ctx := c.Request.Context()
			reqID := utils.GetRequestID(ctx)

			utils.GetLogger().Warn("unauthorized request",
				zap.String("req_id", reqID),
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)

			response.Fail(c, app_error.ErrUnauthorized)
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), utils.ContextUserId, uid)
		ctx = context.WithValue(ctx, utils.ContextUserRole, role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

package observability

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamKienb/shopify-go-platform/utils"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		ctx := c.Request.Context()

		reqID := utils.GetRequestID(ctx)
		userID := utils.GetUserID(ctx)

		span := trace.SpanFromContext(ctx)
		traceID := span.SpanContext().TraceID().String()

		status := c.Writer.Status()
		latency := time.Since(start)

		fields := []zap.Field{
			zap.String("req_id", reqID),
			zap.String("trace_id", traceID),
			zap.String("user_id", userID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
			zap.String("ua", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			utils.GetLogger().Error("request failed",
				append(fields, zap.String("error", c.Errors.String()))...,
			)
			return
		}

		utils.GetLogger().Info("request completed", fields...)

	}
}

package observability

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamKienb/shopify-go-platform/utils"
	"go.opentelemetry.io/otel"
)

func Tracing() gin.HandlerFunc {
	tracer := otel.Tracer("gin-server")

	return func(c *gin.Context) {
		requestID := uuid.New().String()

		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		c.Writer.Header().Set(utils.HeaderRequestId, requestID)

		ctx = utils.SetRequestID(ctx, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

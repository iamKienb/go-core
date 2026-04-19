package observability

// import (
// 	"runtime/debug"

// 	"github.com/gin-gonic/gin"
// 	app_error "github.com/iamKienb/shopify-go-platform/error"
// 	"github.com/iamKienb/shopify-go-platform/response"
// 	"github.com/iamKienb/shopify-go-platform/utils"
// 	"go.opentelemetry.io/otel/trace"
// 	"go.uber.org/zap"
// )

// func Recovery() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				ctx := c.Request.Context()

// 				reqID := c.Writer.Header().Get(utils.HeaderRequestId)

// 				span := trace.SpanFromContext(ctx)

// 				traceID := span.SpanContext().TraceID().String()

// 				utils.GetLogger().Error("panic recovered",
// 					zap.String("req_id", reqID),
// 					zap.String("trace_id", traceID),
// 					zap.Any("error", err),
// 					zap.ByteString("stack", debug.Stack()),
// 				)

// 				response.Fail(c, *app_error.Internal())
// 				c.Abort()
// 			}
// 		}()
// 		c.Next()
// 	}
// }

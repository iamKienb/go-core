package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	app_error "github.com/iamKienb/shopify-go-platform/error"
	"github.com/iamKienb/shopify-go-platform/response"
	"github.com/iamKienb/shopify-go-platform/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var rateLimitScript = redis.NewScript(`
	local key = KEYS[1]
	local window = tonumber(ARGV[1])

	local current = redis.call("INCR", key)

	if current == 1 then
		redis.call("EXPIRE", key, window)
	else
		if redis.call("TTL", key) < 0 then
			redis.call("EXPIRE", key, window)
		end
	end
	return current
`)

func RateLimit(rdb *redis.Client, limit int, windowSec int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userID := utils.GetUserID(ctx)
		reqID := utils.GetRequestID(ctx)

		identifier := c.ClientIP()
		if userID != "" {
			identifier = userID
		}

		key := fmt.Sprintf("rate:%s:%d", identifier, time.Now().Unix()/int64(windowSec))

		count, err := rateLimitScript.Run(ctx, rdb, []string{key}, windowSec).Int()

		if err != nil {
			utils.GetLogger().Error("rate limit redis error",
				zap.String("req_id", reqID),
				zap.Error(err),
			)

			c.Next()
			return
		}

		if count > limit {

			utils.GetLogger().Warn("rate limit exceeded",
				zap.String("req_id", reqID),
				zap.String("user_id", userID),
				zap.String("ip", c.ClientIP()),
				zap.Int("count", count),
			)

			meta := map[string]any{
				"request_id": utils.GetRequestID(c),
				"user_id":    userID,
			}

			response.Fail(c, app_error.ErrTooManyRequests, meta)
			c.Abort()
			return
		}

		c.Next()
	}
}

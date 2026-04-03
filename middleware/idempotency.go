package middleware

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	app_error "github.com/iamKienb/shopify-go-platform/error"
	"github.com/iamKienb/shopify-go-platform/response"
	"github.com/iamKienb/shopify-go-platform/utils"
	"github.com/redis/go-redis/v9"
)

type cachedResponse struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func Idempotency(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		key := c.GetHeader(utils.HeaderIdemKey)
		if key == "" {
			c.Next()
			return
		}

		userID := utils.GetUserID(ctx)
		reqID := utils.GetRequestID(ctx)

		redisKey := "idem:" + userID + ":" + key

		success, err := rdb.SetNX(c, redisKey, "processing", 2*time.Minute).Result()
		if err != nil {
			c.Next()
			return
		}

		if !success {
			val, _ := rdb.Get(c, redisKey).Result()

			if val == "processing" {
				meta := map[string]any{
					"request_id": reqID,
					"user_id":    userID,
				}
				response.Fail(c, app_error.ErrTooManyRequests, meta)
				return
			}

			var cached cachedResponse
			_ = json.Unmarshal([]byte(val), &cached)

			c.Data(
				cached.Status,
				"application/json",
				[]byte(cached.Body),
			)
			return
		}

		rw := utils.NewResponseWriter(c.Writer)
		c.Writer = rw

		c.Next()

		if c.Writer.Status() < 300 {

			resp := cachedResponse{
				Status: c.Writer.Status(),
				Body:   rw.Body(),
			}

			data, _ := json.Marshal(resp)

			rdb.Set(context.Background(), redisKey, data, 24*time.Hour)

		} else {
			rdb.Del(context.Background(), redisKey)
		}
	}
}

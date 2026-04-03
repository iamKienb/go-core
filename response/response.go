package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	app_error "github.com/iamKienb/shopify-go-platform/error"
)

type ApiResponse[T any] struct {
	Code     int            `json:"code"`
	Message  string         `json:"message"`
	Data     T              `json:"data,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func Ok[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, ApiResponse[T]{
		Code:    0,
		Message: "Success",
		Data:    data,
	})
}

func Fail(c *gin.Context, err app_error.AppError, meta ...map[string]any) {
	var finalMeta map[string]any

	if len(meta) > 0 {
		finalMeta = meta[0]
	}

	c.JSON(err.HTTPStatus, ApiResponse[any]{
		Code:     err.Code,
		Message:  err.Message,
		Metadata: finalMeta,
	})
}

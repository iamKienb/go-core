package utils

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func NewResponseWriter(w gin.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		body:           bytes.NewBuffer(nil),
	}
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (r *ResponseWriter) Body() string {
	return r.body.String()
}

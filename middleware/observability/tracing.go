package observability

import (
	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
)

func NewTracingInterceptor() (connect.Interceptor, error) {
	return otelconnect.NewInterceptor(
		otelconnect.WithTrustRemote(),
	)
}

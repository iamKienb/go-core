package observability

import (
	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
)

func TracingInterceptor() (connect.Interceptor, error) {
	return otelconnect.NewInterceptor(
		otelconnect.WithTrustRemote(),
	)
}

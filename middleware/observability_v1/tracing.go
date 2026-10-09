package observability_v1

import (
	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
)

func TracingInterceptor() (connect.Interceptor, error) {
	return otelconnect.NewInterceptor(
		otelconnect.WithTrustRemote(),
	)
}

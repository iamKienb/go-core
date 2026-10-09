package shared

import "context"

const (
	requestIDHeaderKey string = "request_id"
	traceIDHeaderKey   string = "trace_id"
)

func SetRequestIDToCtx(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDHeaderKey, reqID)
}

func SetTraceIDToCtx(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDHeaderKey, traceID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDHeaderKey).(string); ok {
		return id
	}

	return ""
}

func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDHeaderKey).(string); ok {
		return id
	}

	return ""
}

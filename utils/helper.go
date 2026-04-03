package utils

import (
	"context"
)

func GetUserID(ctx context.Context) string {
	val, ok := ctx.Value(ContextUserId).(string)
	if !ok {
		return ""
	}

	return val
}

func GetRequestID(ctx context.Context) string {
	val, ok := ctx.Value(ContextRequestId).(string)
	if !ok {
		return ""
	}

	return val
}

func GetUserRole(ctx context.Context) string {
	val, ok := ctx.Value(ContextUserRole).(string)
	if !ok {
		return ""
	}

	return val
}

func SetRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, ContextRequestId, reqID)
}

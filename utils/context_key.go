package utils

type contextKey string

const (
	ContextRequestId contextKey = "request_id"
	ContextUserId    contextKey = "user_id"
	ContextIdemKey   contextKey = "idempotency_key"
	ContextUserRole  contextKey = "user_role"
)

const (
	HeaderRequestId = "X-Request-ID"
	HeaderUserId    = "X-User-ID"
	HeaderUserRole  = "X-User-Role"
	HeaderIdemKey   = "X-Idempotency-Key"
)

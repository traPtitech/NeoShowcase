package slogutil

type contextKey string

const (
	TraceIDKey contextKey = "ns_trace_id"
	UserIDKey  contextKey = "ns_user_id"
)

package config

const (
	// AppName is Application name
	AppName = "rtm-api"
	// APIVersion is API version
	APIVersion = "0"
	// BuildVersion is API build version
	BuildVersion = "0.3.0"

	CtxSubscription ctxKey = iota
	CtxTracerTransaction
	CtxTracerSpan
)

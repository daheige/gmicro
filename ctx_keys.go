package gmicro

// CtxKey ctx key type.
type CtxKey string

const (
	// XRequestID request_id
	XRequestID CtxKey = "x-request-id"

	// GRPCClientIP grpc client_ip
	GRPCClientIP CtxKey = "grpc_client_ip"

	// RequestMethod request method
	RequestMethod CtxKey = "request_method"

	// RequestURI request uri
	RequestURI CtxKey = "request_uri"
)

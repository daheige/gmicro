package gmicro

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
)

// Route represents the route for mux
type Route struct {
	Method  string
	Path    string
	Handler runtime.HandlerFunc
}

// AllPattern returns a pattern which matches any url
func AllPattern() runtime.Pattern {
	return runtime.MustPattern(runtime.NewPattern(1, []int{int(utilities.OpPush), 0}, []string{""}, ""))
}

package gmicro

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
)

// Route represents the route for mux
type Route struct {
	Method  string
	Pattern runtime.Pattern
	Handler runtime.HandlerFunc
}

// PathPattern returns a pattern which matches exactly with the path
func PathPattern(path string) runtime.Pattern {
	path = strings.TrimPrefix(path, "/")
	return runtime.MustPattern(runtime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{path}, ""))
}

// AllPattern returns a pattern which matches any url
func AllPattern() runtime.Pattern {
	return runtime.MustPattern(runtime.NewPattern(1, []int{int(utilities.OpPush), 0}, []string{""}, ""))
}

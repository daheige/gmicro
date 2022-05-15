package gmicro

import (
	"strings"

	gRuntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
)

// Route represents the route for mux
type Route struct {
	Method  string
	Pattern gRuntime.Pattern
	Handler gRuntime.HandlerFunc
}

// PathPattern returns a pattern which matches exactly with the path
func PathPattern(path string) gRuntime.Pattern {
	path = strings.TrimPrefix(path, "/")
	return gRuntime.MustPattern(gRuntime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{path}, ""))
}

// AllPattern returns a pattern which matches any url
func AllPattern() gRuntime.Pattern {
	return gRuntime.MustPattern(gRuntime.NewPattern(1, []int{int(utilities.OpPush), 0}, []string{""}, ""))
}

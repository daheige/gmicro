package gmicro

import (
	gRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
)

// Route represents the route for mux
type Route struct {
	Method  string
	Path    string
	Handler gRuntime.HandlerFunc
}

// AllPattern returns a pattern which matches any url
func AllPattern() gRuntime.Pattern {
	return gRuntime.MustPattern(gRuntime.NewPattern(1, []int{int(utilities.OpPush), 0}, []string{""}, ""))
}

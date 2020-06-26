// Package gmicro Grpc Microservices components.
package gmicro

import (
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

const (
	// the default timeout before the server shutdown abruptly
	defaultShutdownTimeout = 10 * time.Second
)

// Service represents the microservice.
type Service struct {
	GRPCServer       *grpc.Server // grpc server
	HTTPServer       *http.Server // if you need grpc gw,please use it
	httpHandler      http.Handler
	recovery         func()
	shutdownFunc     func() error // shutdown func
	shutdownTimeout  time.Duration
	interruptSignals []os.Signal // interrupt signal

	errorHandler       runtime.ProtoErrorHandlerFunc  // grpc error handler
	gMux               *runtime.ServeMux              // grpc gw runtime serverMux
	gMuxOptions        []runtime.ServeMuxOption       // grpc mux options
	streamInterceptors []grpc.StreamServerInterceptor // grpc steam interceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor  // grpc server interceptor
	grpcServerOptions  []grpc.ServerOption
	grpcDialOptions    []grpc.DialOption
	logger             Logger
}

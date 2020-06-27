// Package gmicro Grpc Microservices components.
package gmicro

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// the default timeout before the server shutdown abruptly
	defaultShutdownTimeout = 10 * time.Second

	// the default time waiting for running goroutines to finish their jobs before the shutdown starts
	defaultPreShutdownDelay = 1 * time.Second
)

// refer: https://github.com/grpc-ecosystem/grpc-gateway/blob/master/docs/_docs/customizingyourgateway.md
var defaultMuxOption = runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true})

// AnnotatorFunc is the annotator function is for injecting meta data from http request into gRPC context
type AnnotatorFunc func(context.Context, *http.Request) metadata.MD

// ReverseProxyFunc is the callback that the caller should implement
// to steps to reverse-proxy the HTTP/1 requests to gRPC
// handlerFromEndpoint http gw endPoint
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
type ReverseProxyFunc func(ctx context.Context, mux *runtime.ServeMux, grpcAddressAndPort string, opts []grpc.DialOption) error

// HTTPHandlerFunc is the http middleware handler function.
type HTTPHandlerFunc func(*runtime.ServeMux) http.Handler

// Service represents the microservice.
type Service struct {
	GRPCServer         *grpc.Server    // grpc server
	HTTPServer         *http.Server    // if you need grpc gw,please use it
	httpHandler        HTTPHandlerFunc // http.Handler
	grpcAddress        string          // grpc host eg: ip:port
	httpServerAddress  string          // http server host eg: ip:port
	recovery           func()
	shutdownFunc       func() // shutdown func
	shutdownTimeout    time.Duration
	preShutdownDelay   time.Duration
	interruptSignals   []os.Signal // interrupt signal
	annotators         []AnnotatorFunc
	staticDir          string                         // static dir
	errorHandler       runtime.ProtoErrorHandlerFunc  // grpc error handler
	mux                *runtime.ServeMux              // grpc gw runtime serverMux
	muxOptions         []runtime.ServeMuxOption       // grpc mux options
	routes             []Route                        // grpc http router
	streamInterceptors []grpc.StreamServerInterceptor // grpc steam interceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor  // grpc server interceptor
	grpcServerOptions  []grpc.ServerOption
	grpcDialOptions    []grpc.DialOption
	logger             Logger
	reverseProxyFuncs  []ReverseProxyFunc // http gw endpoint
	enablePrometheus   bool               // enable prometheus monitor
}

// DefaultHTTPHandler is the default http handler which does nothing
func DefaultHTTPHandler(mux *runtime.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})
}

// GrpcHandlerFunc uses the standard library h2c to convert http requests to http2
// In this way, you can co-exist with go grpc and http services, and use one port
// to provide both grpc services and http services.
// In June 2018, the golang.org/x/net/http2/h2c standard library representing the "h2c"
// logo was officially merged in, and since then we can use the official standard library (h2c)
// This standard library implements the unencrypted mode of HTTP/2,
// so we can use the standard library to provide both HTTP/1.1 and HTTP/2 functions on the same port
// The h2c.NewHandler method has been specially processed, and h2c.NewHandler will return an http.handler
// The main internal logic is to intercept all h2c traffic, then hijack and redirect it
// to the corresponding Hander according to different request traffic types to process
func GrpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func defaultService() *Service {
	s := Service{}
	s.httpHandler = DefaultHTTPHandler
	s.errorHandler = runtime.DefaultHTTPError
	s.shutdownFunc = func() {}
	s.shutdownTimeout = defaultShutdownTimeout
	s.preShutdownDelay = defaultPreShutdownDelay
	s.logger = dummyLogger

	// default interrupt signals to catch, you can use InterruptSignal option to append more
	s.interruptSignals = InterruptSignals
	s.streamInterceptors = []grpc.StreamServerInterceptor{}
	s.unaryInterceptors = []grpc.UnaryServerInterceptor{}

	// install validator interceptor
	s.streamInterceptors = append(s.streamInterceptors, grpc_validator.StreamServerInterceptor())
	s.unaryInterceptors = append(s.unaryInterceptors, grpc_validator.UnaryServerInterceptor())

	// install panic handler which will turn panics into gRPC errors
	s.streamInterceptors = append(s.streamInterceptors, grpc_recovery.StreamServerInterceptor())
	s.unaryInterceptors = append(s.unaryInterceptors, grpc_recovery.UnaryServerInterceptor())

	// default dial option is using insecure connection
	if len(s.grpcDialOptions) == 0 {
		s.grpcDialOptions = append(s.grpcDialOptions, grpc.WithInsecure())
	}

	// apply default marshaler option for mux, can be replaced by using MuxOption
	s.muxOptions = append(s.muxOptions, defaultMuxOption)

	return &s
}

// NewService creates a new microservice
func NewService(opts ...Option) *Service {
	s := defaultService()

	// app option functions.
	s.apply(opts...)

	// install prometheus interceptor
	if s.enablePrometheus {
		s.streamInterceptors = append(s.streamInterceptors, grpc_prometheus.StreamServerInterceptor)
		s.unaryInterceptors = append(s.unaryInterceptors, grpc_prometheus.UnaryServerInterceptor)

		// add /metrics HTTP/1 endpoint
		routeMetrics := Route{
			Method:  "GET",
			Pattern: PathPattern("metrics"),
			Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
				promhttp.Handler().ServeHTTP(w, r)
			},
		}

		s.routes = append(s.routes, routeMetrics)
	}

	// init gateway mux
	s.muxOptions = append(s.muxOptions, runtime.WithProtoErrorHandler(s.errorHandler))

	// init annotators
	for _, annotator := range s.annotators {
		s.muxOptions = append(s.muxOptions, runtime.WithMetadata(annotator))
	}

	s.mux = runtime.NewServeMux(s.muxOptions...)

	s.grpcServerOptions = append(s.grpcServerOptions, grpc_middleware.WithStreamServerChain(s.streamInterceptors...))
	s.grpcServerOptions = append(s.grpcServerOptions, grpc_middleware.WithUnaryServerChain(s.unaryInterceptors...))

	s.GRPCServer = grpc.NewServer(
		s.grpcServerOptions...,
	)

	// default http server config
	// http server addr is specified in the startGRPCGateway method below
	if s.HTTPServer == nil {
		s.HTTPServer = &http.Server{
			ReadHeaderTimeout: 5 * time.Second,  //read header timeout
			ReadTimeout:       5 * time.Second,  //read request timeout
			WriteTimeout:      10 * time.Second, //write timeout
			IdleTimeout:       20 * time.Second, //tcp idle time
		}
	}

	return s
}

// Getpid gets the process id of server
func (s *Service) Getpid() int {
	return os.Getpid()
}

// Start starts the microservice with listening on the ports
// start grpc gateway and http server on different port
func (s *Service) Start(httpPort int, grpcPort int, reverseProxyFunc ...ReverseProxyFunc) error {
	// http gw host and grpc host
	s.httpServerAddress = fmt.Sprintf("0.0.0.0:%d", httpPort)
	s.grpcAddress = fmt.Sprintf("0.0.0.0:%d", grpcPort)

	if len(reverseProxyFunc) > 0 {
		s.reverseProxyFuncs = append(s.reverseProxyFuncs, reverseProxyFunc...)
	}

	// intercept interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, s.interruptSignals...)

	// channels to receive error
	errChan1 := make(chan error, 1)
	errChan2 := make(chan error, 1)

	// start gRPC server
	go func() {
		s.logger.Printf("Starting gPRC server listening on %d", grpcPort)
		errChan1 <- s.startGRPCServer()
	}()

	// start HTTP/1.0 gateway server
	go func() {
		s.logger.Printf("Starting http server listening on %d", httpPort)
		errChan2 <- s.startGRPCGateway()
	}()

	// wait for context cancellation or shutdown signal
	select {
	// if gRPC server fail to start
	case err := <-errChan1:
		return err

	// if http server fail to start
	case err := <-errChan2:
		return err

	// if we received an interrupt signal
	case sig := <-sigChan:
		s.logger.Printf("Interrupt signal received: %v", sig)
		s.Stop()
		return nil
	}
}

// startGRPCServer start grpc server.
func (s *Service) startGRPCServer() error {
	// register reflection service on gRPC server.
	// reflection.Register(s.GRPCServer)

	lis, err := net.Listen("tcp", s.grpcAddress)
	if err != nil {
		return err
	}

	return s.GRPCServer.Serve(lis)
}

func (s *Service) startGRPCGateway() error {
	// apply routes
	for _, route := range s.routes {
		s.mux.Handle(route.Method, route.Pattern, route.Handler)
	}

	// Register http gw handlerFromEndpoint
	ctx := context.Background()
	var err error
	for _, h := range s.reverseProxyFuncs {
		err = h(ctx, s.mux, s.grpcAddress, s.grpcDialOptions)
		if err != nil {
			s.logger.Printf("register handler from endPoint error: %s", err.Error())
			return err
		}
	}

	// this is the fallback handler that will serve static files,
	// if file does not exist, then a 404 error will be returned.
	s.mux.Handle("GET", AllPattern(), func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		dir := s.staticDir
		if s.staticDir == "" {
			dir, _ = os.Getwd()
		}

		// check if the file exists and fobid showing directory
		path := filepath.Join(dir, r.URL.Path)
		if fileInfo, err := os.Stat(path); os.IsNotExist(err) || fileInfo.IsDir() {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, path)
	})

	// http server
	s.HTTPServer.Addr = s.httpServerAddress
	s.HTTPServer.Handler = s.httpHandler(s.mux)
	s.HTTPServer.RegisterOnShutdown(s.shutdownFunc)

	return s.HTTPServer.ListenAndServe()
}

// Stop stops the microservice gracefully
func (s *Service) Stop() {
	// disable keep-alives on existing connections
	s.HTTPServer.SetKeepAlivesEnabled(false)

	// we wait for a duration of preShutdownDelay for running goroutines to finish their jobs
	if s.preShutdownDelay > 0 {
		s.logger.Printf("Waiting for %v before shutdown starts", s.preShutdownDelay)
		time.Sleep(s.preShutdownDelay)
	}

	// gracefully stop gRPC server first
	s.GRPCServer.GracefulStop()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer cancel()

	// gracefully stop http server
	go s.HTTPServer.Shutdown(ctx)
	<-ctx.Done()
}

/**
* The following method is mainly for grpc server and http gw server to start on one port
 */

// StartGRPCAndHTTPServer grpc server and grpc gateway port share a port
func (s *Service) StartGRPCAndHTTPServer(port int) error {
	// http gw host and grpc host
	s.httpServerAddress = fmt.Sprintf("0.0.0.0:%d", port)
	s.grpcAddress = s.httpServerAddress

	// intercept interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, s.interruptSignals...)

	// channels to receive error
	errChan := make(chan error, 1)

	// start HTTP/1.0 gateway server and  gRPC server.
	go func() {
		s.logger.Printf("Starting http server and grp server listening on %d", port)
		errChan <- s.startWithSharePort()
	}()

	// wait for context cancellation or shutdown signal
	select {
	// if http server and gRPC server fail to start
	case err := <-errChan:
		return err
	// if we received an interrupt signal
	case sig := <-sigChan:
		s.logger.Printf("Interrupt signal received: %v", sig)
		s.stopGRPCAndHTTPServer()
		return nil
	}
}

func (s *Service) startWithSharePort() error {
	// apply routes
	for _, route := range s.routes {
		s.mux.Handle(route.Method, route.Pattern, route.Handler)
	}

	ctx := context.Background()
	var err error
	for _, h := range s.reverseProxyFuncs {
		err = h(ctx, s.mux, s.grpcAddress, s.grpcDialOptions)
		if err != nil {
			s.logger.Printf("register handler from endPoint error: %s", err.Error())
		}
	}

	// http server and h2c handler
	// create a http mux
	httpMux := http.NewServeMux()
	httpMux.Handle("/", s.mux)

	s.HTTPServer.Addr = s.httpServerAddress
	s.HTTPServer.Handler = GrpcHandlerFunc(s.GRPCServer, httpMux)
	s.HTTPServer.RegisterOnShutdown(s.shutdownFunc)

	return s.HTTPServer.ListenAndServe()
}

func (s *Service) stopGRPCAndHTTPServer() {
	// disable keep-alives on existing connections
	s.HTTPServer.SetKeepAlivesEnabled(false)

	// we wait for a duration of preShutdownDelay for running goroutines to finish their jobs
	if s.preShutdownDelay > 0 {
		s.logger.Printf("Waiting for %v before shutdown starts", s.preShutdownDelay)
		time.Sleep(s.preShutdownDelay)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer cancel()

	// gracefully stop http server
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go s.HTTPServer.Shutdown(ctx)
	<-ctx.Done()
}

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
	"runtime/debug"
	"strings"
	"time"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	gRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	gValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	gPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	gRuntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	// the default timeout before the server shutdown abruptly
	defaultShutdownTimeout = 5 * time.Second

	// the default time waiting for running goroutines to finish their jobs before the shutdown start.
	defaultPreShutdownDelay = 2 * time.Second
)

// refer: https://github.com/golang/protobuf/blob/v1.4.3/jsonpb/encode.go#L30
var defaultMuxOption = gRuntime.WithMarshalerOption(gRuntime.MIMEWildcard,
	&gRuntime.JSONPb{EmitDefaults: true, OrigName: false})

// AnnotatorFunc is the annotator function is for injecting meta data from http request into gRPC context
type AnnotatorFunc func(context.Context, *http.Request) metadata.MD

// HandlerFromEndpoint is the callback that the caller should implement
// to steps to reverse-proxy the HTTP/1 requests to gRPC
// handlerFromEndpoint http gw endPoint
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
type HandlerFromEndpoint func(ctx context.Context, mux *gRuntime.ServeMux,
	grpcAddressAndPort string, opts []grpc.DialOption) error

// HTTPHandlerFunc is the http middleware handler function.
type HTTPHandlerFunc func(*gRuntime.ServeMux) http.Handler

// Service represents the microservice.
type Service struct {
	GRPCServer           *grpc.Server    // gRPC server
	HTTPServer           *http.Server    // if you need gRPC gw,please use it
	httpHandler          HTTPHandlerFunc // http.Handler
	gRPCAddress          string          // gRPC host eg: ip:port
	httpServerAddress    string          // http server host eg: ip:port
	gRPCNetwork          string          // the gRPC network must be "tcp", "tcp4", "tcp6"
	recovery             func()          // goroutine exec recover catch stack
	shutdownFunc         func()          // shutdown func
	shutdownTimeout      time.Duration   // shutdown wait time
	preShutdownDelay     time.Duration
	interruptSignals     []os.Signal // interrupt signal
	annotators           []AnnotatorFunc
	staticDir            string                         // static dir
	enableStaticAccess   bool                           // enable static file access
	errorHandler         gRuntime.ProtoErrorHandlerFunc // gRPC error handler
	mux                  *gRuntime.ServeMux             // gRPC gw runtime serverMux
	muxOptions           []gRuntime.ServeMuxOption      // gRPC mux options
	routes               []Route                        // gRPC http router
	streamInterceptors   []grpc.StreamServerInterceptor // gRPC steam interceptor
	unaryInterceptors    []grpc.UnaryServerInterceptor  // gRPC server interceptor
	enableRequestAccess  bool                           // gRPC request log config
	gRPCServerOptions    []grpc.ServerOption
	gRPCDialOptions      []grpc.DialOption
	logger               Logger                // logger interface entry
	handlerFromEndpoints []HandlerFromEndpoint // http gw endpoint
	enablePrometheus     bool                  // enable prometheus monitor
}

// DefaultHTTPHandler is the default http handler which does nothing.
func DefaultHTTPHandler(mux *gRuntime.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})
}

// GRPCHandlerFunc uses the standard library h2c to convert http requests to http2
// In this way, you can co-exist with go grpc and http services, and use one port
// to provide both grpc services and http services.
// In June 2018, the golang.org/x/net/http2/h2c standard library representing the "h2c"
// logo was officially merged in, and since then we can use the official standard library (h2c)
// This standard library implements the unencrypted mode of HTTP/2,
// so we can use the standard library to provide both HTTP/1.1 and HTTP/2 functions on the same port
// The h2c.NewHandler method has been specially processed, and h2c.NewHandler will return an http.handler
// The main internal logic is to intercept all h2c traffic, then hijack and redirect it
// to the corresponding Hander according to different request traffic types to process
func GRPCHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
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
	s.errorHandler = gRuntime.DefaultHTTPError
	s.shutdownFunc = func() {}
	s.shutdownTimeout = defaultShutdownTimeout
	s.preShutdownDelay = defaultPreShutdownDelay
	s.logger = dummyLogger

	// goroutine recover catch stack
	s.recovery = func() {
		defer func() {
			if e := recover(); e != nil {
				s.logger.Printf("exec recover err: %v\n", e)
				s.logger.Printf("full stack: %s\n", string(debug.Stack()))
			}
		}()
	}

	// default interrupt signals to catch, you can use InterruptSignal option to append more
	s.interruptSignals = InterruptSignals

	// register interceptor
	s.streamInterceptors = make([]grpc.StreamServerInterceptor, 0, 20)
	s.unaryInterceptors = make([]grpc.UnaryServerInterceptor, 0, 20)

	// install panic handler which will turn panics into gRPC errors.
	s.streamInterceptors = append(s.streamInterceptors, gRecovery.StreamServerInterceptor())
	s.unaryInterceptors = append(s.unaryInterceptors, gRecovery.UnaryServerInterceptor())

	// install validator interceptor.
	s.streamInterceptors = append(s.streamInterceptors, gValidator.StreamServerInterceptor())
	s.unaryInterceptors = append(s.unaryInterceptors, gValidator.UnaryServerInterceptor())

	// apply default marshaler option for mux, can be replaced by using MuxOption
	s.muxOptions = append(s.muxOptions, defaultMuxOption)

	return &s
}

// NewService creates a new microservice
func NewService(opts ...Option) *Service {
	s := defaultService()

	// app option functions.
	s.apply(opts)

	// install request interceptor
	if s.enableRequestAccess {
		s.unaryInterceptors = append(s.unaryInterceptors, s.RequestInterceptor)
	}

	// default dial option is using insecure connection
	if len(s.gRPCDialOptions) == 0 {
		s.gRPCDialOptions = append(s.gRPCDialOptions, grpc.WithInsecure())
	}

	// install prometheus interceptor
	if s.enablePrometheus {
		s.streamInterceptors = append(s.streamInterceptors, gPrometheus.StreamServerInterceptor)
		s.unaryInterceptors = append(s.unaryInterceptors, gPrometheus.UnaryServerInterceptor)

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
	s.muxOptions = append(s.muxOptions, gRuntime.WithProtoErrorHandler(s.errorHandler))

	// init annotators
	for _, annotator := range s.annotators {
		s.muxOptions = append(s.muxOptions, gRuntime.WithMetadata(annotator))
	}

	s.mux = gRuntime.NewServeMux(s.muxOptions...)

	s.gRPCServerOptions = append(s.gRPCServerOptions,
		middleware.WithStreamServerChain(s.streamInterceptors...))
	s.gRPCServerOptions = append(s.gRPCServerOptions,
		middleware.WithUnaryServerChain(s.unaryInterceptors...))

	s.GRPCServer = grpc.NewServer(
		s.gRPCServerOptions...,
	)

	// default http server config
	// http server addr is specified in the startGRPCGateway method below
	if s.HTTPServer == nil {
		s.HTTPServer = &http.Server{
			ReadHeaderTimeout: 5 * time.Second,  // read header timeout
			ReadTimeout:       5 * time.Second,  // read request timeout
			WriteTimeout:      10 * time.Second, // write timeout
			IdleTimeout:       20 * time.Second, // tcp idle time
		}
	}

	return s
}

// RequestInterceptor request interceptor to record basic information of the request
func (s *Service) RequestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (res interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			// the error format defined by grpc must be used here to return code, desc
			err = status.Errorf(codes.Internal, "%s", "server inner error")

			s.logger.Printf("reply: %v\n", res)
			s.logger.Printf("exec panic: %v\n", r)
			s.logger.Printf("full stack: %s\n", string(debug.Stack()))
		}
	}()

	t := time.Now()
	clientIP, _ := GetGRPCClientIP(ctx)

	s.logger.Printf("exec begin\n")
	s.logger.Printf("client_ip: %s\n", clientIP)
	// s.logger.Printf("request: %v\n", req)

	// request ctx key
	if logID := ctx.Value(XRequestID); logID == nil {
		ctx = context.WithValue(ctx, XRequestID, RndUUID())
	}

	ctx = context.WithValue(ctx, GRPCClientIP, clientIP)
	ctx = context.WithValue(ctx, RequestMethod, info.FullMethod)
	ctx = context.WithValue(ctx, RequestURI, info.FullMethod)

	res, err = handler(ctx, req)
	ttd := time.Since(t).Milliseconds()
	if err != nil {
		s.logger.Printf("trace_error: %s\n", err.Error())
		s.logger.Printf("exec time: %v\n", ttd)
		s.logger.Printf("reply: %v\n", res)

		return nil, err
	}

	s.logger.Printf("exec end,cost time: %v ms\n", ttd)

	return res, err
}

// GetPid gets the process id of server
func (s *Service) GetPid() int {
	return os.Getpid()
}

// AddHandlerFromEndpoint add HandlerFromEndpoint.
func (s *Service) AddHandlerFromEndpoint(h ...HandlerFromEndpoint) {
	s.handlerFromEndpoints = append(s.handlerFromEndpoints, h...)
}

// AddRoute add some route to routes
func (s *Service) AddRoute(routes ...Route) {
	s.routes = append(s.routes, routes...)
}

// Start starts the microservice with listening on the ports
// start grpc gateway and http server on different port
func (s *Service) Start(httpPort int, grpcPort int) error {
	// http gw host and grpc host
	s.httpServerAddress = fmt.Sprintf("0.0.0.0:%d", httpPort)
	s.gRPCAddress = fmt.Sprintf("0.0.0.0:%d", grpcPort)

	// intercept interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, s.interruptSignals...)

	// channels to receive error
	errChan1 := make(chan error, 1)
	errChan2 := make(chan error, 1)

	// start gRPC server
	go func() {
		defer s.recovery()

		s.logger.Printf("Starting gPRC server listening on %d\n", grpcPort)
		errChan1 <- s.startGRPCServer()
	}()

	// start HTTP/1.0 gateway server
	go func() {
		defer s.recovery()

		s.logger.Printf("Starting http server listening on %d\n", httpPort)
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
		s.logger.Printf("Interrupt signal received: %v\n", sig)
		s.Stop()
		return nil
	}
}

// startGRPCServer start grpc server.
func (s *Service) startGRPCServer() error {
	// register reflection service on gRPC server.
	reflection.Register(s.GRPCServer)

	if s.gRPCNetwork == "" {
		s.gRPCNetwork = "tcp"
	}

	lis, err := net.Listen(s.gRPCNetwork, s.gRPCAddress)
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
	for _, h := range s.handlerFromEndpoints {
		err = h(ctx, s.mux, s.gRPCAddress, s.gRPCDialOptions)
		if err != nil {
			s.logger.Printf("register handler from endPoint error: %s\n", err.Error())
			return err
		}
	}

	// static file access
	if s.enableStaticAccess {
		// this is the fallback handler that will serve static files,
		// if file does not exist, then a 404 error will be returned.
		s.mux.Handle("GET", AllPattern(), func(w http.ResponseWriter, r *http.Request,
			pathParams map[string]string) {
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
	}

	// http server
	s.HTTPServer.Addr = s.httpServerAddress
	s.HTTPServer.Handler = s.httpHandler(s.mux)
	s.HTTPServer.RegisterOnShutdown(s.shutdownFunc)

	return s.HTTPServer.ListenAndServe()
}

// Stop stops the microservice gracefully.
func (s *Service) Stop() {
	// disable keep-alives on existing connections
	s.HTTPServer.SetKeepAlivesEnabled(false)

	// we wait for a duration of preShutdownDelay for running goroutines to finish their jobs
	if s.preShutdownDelay > 0 {
		s.logger.Printf("Waiting for %v before shutdown start\n", s.preShutdownDelay)
		time.Sleep(s.preShutdownDelay)
	}

	// gracefully stop gRPC server first
	s.GRPCServer.GracefulStop()

	// gracefully stop http server
	s.httpServerShutdown()
}

// httpServerShutdown http gateway server graceful shutdown.
func (s *Service) httpServerShutdown() {
	done := make(chan struct{}, 1)
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
	// gracefully stop http server
	go func() {
		defer s.recovery()
		defer close(done)

		if err := s.HTTPServer.Shutdown(ctx); err != nil {
			s.logger.Printf("Http server shutdown error: %v", err.Error())
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Printf("Server shutdown ctx cancel error: %v", ctx.Err())
	case <-done:
		s.logger.Printf("Server shutdown success")
	}
}

// ===The following method is mainly for grpc server and http gw server to start on one port==//
// referr: https://github.com/daheige/go-proj/blob/master/cmd/rpc/http/server.go#L123

// StartGRPCAndHTTPServer grpc server and grpc gateway port share a port
// error: rpc error: code = Unavailable desc = all SubConns are in TransientFailure,
// latest connection error: timed out waiting for server handshake
// please set this gRPC var.
// export GRPC_GO_REQUIRE_HANDSHAKE=off
func (s *Service) StartGRPCAndHTTPServer(port int) error {
	// http gw host and grpc host
	s.httpServerAddress = fmt.Sprintf("0.0.0.0:%d", port)
	s.gRPCAddress = s.httpServerAddress

	// intercept interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, s.interruptSignals...)

	// channels to receive error
	errChan := make(chan error, 1)

	// start HTTP/1.0 gateway server and  gRPC server.
	go func() {
		defer s.recovery()

		s.logger.Printf("Starting http server and grpc server listening on %d\n", port)
		errChan <- s.startGRPCAndHTTPServer()
	}()

	// wait for context cancellation or shutdown signal
	select {
	// if http server and gRPC server fail to start
	case err := <-errChan:
		return err
	// if we received an interrupt signal
	case sig := <-sigChan:
		s.logger.Printf("Interrupt signal received: %v\n", sig)
		s.stopGRPCAndHTTPServer()
		return nil
	}
}

func (s *Service) startGRPCAndHTTPServer() error {
	// apply routes
	for _, route := range s.routes {
		s.mux.Handle(route.Method, route.Pattern, route.Handler)
	}

	ctx := context.Background()
	var err error
	for _, h := range s.handlerFromEndpoints {
		err = h(ctx, s.mux, s.gRPCAddress, s.gRPCDialOptions)
		if err != nil {
			s.logger.Printf("register handler from endPoint error: %s\n", err.Error())
		}
	}

	// http server and h2c handler
	// create a http mux
	httpMux := http.NewServeMux()
	httpMux.Handle("/", s.mux)

	s.HTTPServer.Addr = s.httpServerAddress

	// gRPC server handler convert to http handler.
	s.HTTPServer.Handler = GRPCHandlerFunc(s.GRPCServer, httpMux)
	s.HTTPServer.RegisterOnShutdown(s.shutdownFunc)

	return s.HTTPServer.ListenAndServe()
}

func (s *Service) stopGRPCAndHTTPServer() {
	// disable keep-alives on existing connections
	s.HTTPServer.SetKeepAlivesEnabled(false)

	// we wait for a duration of preShutdownDelay for running goroutines to finish their jobs
	if s.preShutdownDelay > 0 {
		s.logger.Printf("Waiting for %v before shutdown start\n", s.preShutdownDelay)
		time.Sleep(s.preShutdownDelay)
	}

	// graceful server shutdown
	s.httpServerShutdown()
}

// The following method is only used to start the grpc server, but not start http gw.

// NewServiceWithoutGateway new a service without http gw.
func NewServiceWithoutGateway(opts ...Option) *Service {
	s := defaultService()

	// app option functions.
	s.apply(opts)

	// install request interceptor
	if s.enableRequestAccess {
		s.unaryInterceptors = append(s.unaryInterceptors, s.RequestInterceptor)
	}

	// default dial option is using insecure connection
	if len(s.gRPCDialOptions) == 0 {
		s.gRPCDialOptions = append(s.gRPCDialOptions, grpc.WithInsecure())
	}

	// install prometheus interceptor
	if s.enablePrometheus {
		s.streamInterceptors = append(s.streamInterceptors, gPrometheus.StreamServerInterceptor)
		s.unaryInterceptors = append(s.unaryInterceptors, gPrometheus.UnaryServerInterceptor)
	}

	s.muxOptions = nil

	s.gRPCServerOptions = append(s.gRPCServerOptions,
		middleware.WithStreamServerChain(s.streamInterceptors...))
	s.gRPCServerOptions = append(s.gRPCServerOptions,
		middleware.WithUnaryServerChain(s.unaryInterceptors...))

	s.GRPCServer = grpc.NewServer(
		s.gRPCServerOptions...,
	)

	return s
}

// StartGRPCWithoutGateway start gRPC without gw.
func (s *Service) StartGRPCWithoutGateway(grpcPort int) error {
	s.gRPCAddress = fmt.Sprintf("0.0.0.0:%d", grpcPort)

	// intercept interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, s.interruptSignals...)

	// channels to receive error
	errChan := make(chan error, 1)

	// start gRPC server
	go func() {
		defer s.recovery()

		s.logger.Printf("Starting gPRC server listening on %d\n", grpcPort)
		errChan <- s.startGRPCServer()
	}()

	// wait for context cancellation or shutdown signal
	select {
	// if gRPC server fail to start
	case err := <-errChan:
		return err
	// if we received an interrupt signal
	case sig := <-sigChan:
		s.logger.Printf("Interrupt signal received: %v\n", sig)
		s.StopGRPCWithoutGateway()
		return nil
	}
}

// StopGRPCWithoutGateway stop the gRPC server gracefully
func (s *Service) StopGRPCWithoutGateway() {
	// we wait for a duration of preShutdownDelay for running goroutines to finish their jobs
	if s.preShutdownDelay > 0 {
		s.logger.Printf("Waiting for %v before shutdown start\n", s.preShutdownDelay)
		time.Sleep(s.preShutdownDelay)
	}

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)

	defer cancel()

	// gracefully stop gRPC server
	go func() {
		defer s.recovery()
		defer close(done)

		s.GRPCServer.GracefulStop()
	}()

	select {
	case <-ctx.Done():
		s.logger.Printf("Grpc server shutdown ctx cancel error: %v", ctx.Err())
	case <-done:
		s.logger.Printf("Grpc server shutdown success")
	}
}

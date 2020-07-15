package gmicro

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/daheige/gmicro/example/pb"
)

var reverseProxyFunc HandlerFromEndpoint
var httpPort, grpcPort, sharePort int
var shutdownFunc func()

func initConf() {
	reverseProxyFunc = func(
		ctx context.Context,
		mux *runtime.ServeMux,
		grpcHostAndPort string,
		opts []grpc.DialOption,
	) error {
		return nil
	}

	httpPort = 8888
	grpcPort = 9999
	sharePort = 8081

	shutdownFunc = func() {
		fmt.Println("Server shutting down")
	}
}

/** TestNewService
% go test -v -test.run=TestNewService
=== RUN   TestNewService
2020/06/27 18:56:52 Starting gPRC server listening on 9999
2020/06/27 18:56:52 Starting http server listening on 8888
2020/06/27 18:56:53 req data:  name:"daheige"
2020/06/27 18:56:53 resp code:  200
--- PASS: TestNewService (6.02s)
PASS
ok  	github.com/daheige/gmicro	6.034s
*/
func TestNewService(t *testing.T) {
	initConf()
	var should = require.New(t)

	// add the /test endpoint
	route := Route{
		Method:  "GET",
		Pattern: PathPattern("test"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte("Hello!"))
		},
	}

	// test Option func
	s := NewService(
		WithRouteOpt(route),
		WithShutdownFunc(shutdownFunc),
		WithPreShutdownDelay(1*time.Second),
		WithHandlerFromEndpoint(pb.RegisterGreeterServiceHandlerFromEndpoint),
		// WithHandlerFromEndpoint(HandlerFromEndpoint(pb.RegisterGreeterServiceHandlerFromEndpoint)),
		WithLogger(LoggerFunc(log.Printf)),
		WithPrometheus(true),
	)

	// register grpc service
	pb.RegisterGreeterServiceServer(s.GRPCServer, &greeterService{})

	newRoute := Route{
		Method:  "GET",
		Pattern: PathPattern("health"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute)

	newRoute2 := Route{
		Method:  "GET",
		Pattern: PathPattern("info"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute2)

	s.AddHandlerFromEndpoint(reverseProxyFunc)

	go func() {
		err := s.Start(httpPort, grpcPort)
		should.NoError(err)
	}()

	// wait 1 second for the server start
	time.Sleep(1 * time.Second)

	// check if the http server is up
	httpHost := fmt.Sprintf(":%d", httpPort)
	_, err := net.Listen("tcp", httpHost)
	should.Error(err)

	// check if the grpc server is up
	grpcHost := fmt.Sprintf(":%d", grpcPort)
	_, err = net.Listen("tcp", grpcHost)
	should.Error(err)

	// check if the http endpoint works
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/", httpPort))
	should.NoError(err)
	should.Equal(http.StatusNotFound, resp.StatusCode)

	resp, err = client.Get(fmt.Sprintf("http://127.0.0.1:%d/v1/say/%s", httpPort, "daheige"))
	log.Println("resp code: ", resp.StatusCode)
	should.NoError(err)
	should.Equal(http.StatusOK, resp.StatusCode)

	resp, err = client.Get(fmt.Sprintf("http://127.0.0.1:%d/health", httpPort))
	should.NoError(err)
	should.Equal(http.StatusOK, resp.StatusCode)
	b, err := ioutil.ReadAll(resp.Body)
	should.NoError(err)
	should.Equal("OK", string(b))

	resp, err = client.Get(fmt.Sprintf("http://127.0.0.1:%d/info", httpPort))
	should.NoError(err)
	should.Equal(http.StatusOK, resp.StatusCode)

	resp, err = client.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics", httpPort))
	should.NoError(err)
	should.Equal(http.StatusOK, resp.StatusCode)

	// create service s2 to trigger errChan1
	s2 := NewService()

	s2.AddHandlerFromEndpoint(reverseProxyFunc)

	// grpc port 9999 alreday in use
	err = s2.Start(httpPort, grpcPort)
	should.Error(err)

	// create service s3 to trigger errChan2
	s3 := NewService()

	// http port 8888 already in use
	s.GRPCServer.Stop()
	s3.AddHandlerFromEndpoint(reverseProxyFunc)

	err = s3.Start(httpPort, grpcPort)
	should.Error(err)

	// wait 1 second for s3 gRPC server start
	time.Sleep(1 * time.Second)

	// close all previous services
	s.HTTPServer.Close()
	s3.GRPCServer.Stop()

	// run a new service, we use different ports to make sure ci not complain
	httpPort = 18888
	grpcPort = 19999
	s4 := NewService(
		WithShutdownTimeout(10 * time.Second),
	)
	s4.AddHandlerFromEndpoint(reverseProxyFunc)

	go func() {
		err := s4.Start(httpPort, grpcPort)
		should.NoError(err)
	}()

	// wait 1 second for the server start
	time.Sleep(1 * time.Second)

	// the redoc is not up for the second server
	resp, err = client.Get(fmt.Sprintf("http://127.0.0.1:%d/docs", httpPort))
	should.NoError(err)
	should.Equal(http.StatusNotFound, resp.StatusCode)

	// send an interrupt signal to stop s4
	syscall.Kill(s4.GetPid(), syscall.SIGINT)

	// wait 3 second for the server shutdown
	time.Sleep(3 * time.Second)
}

/** TestErrorReverseProxyFunc
% go test -v -test.run=TestErrorReverseProxyFunc
=== RUN   TestErrorReverseProxyFunc
--- PASS: TestErrorReverseProxyFunc (0.00s)
PASS
ok  	github.com/daheige/gmicro	0.012s
*/
func TestErrorReverseProxyFunc(t *testing.T) {
	initConf()

	var should = require.New(t)

	// mock error from reverseProxyFunc
	errText := "reverse proxy func error"
	reverseProxyFunc = func(
		ctx context.Context,
		mux *runtime.ServeMux,
		grpcHostAndPort string,
		opts []grpc.DialOption,
	) error {
		return errors.New(errText)
	}

	s := NewService(WithHandlerFromEndpoint(reverseProxyFunc))

	// http gw host and grpc host
	s.httpServerAddress = fmt.Sprintf("0.0.0.0:%d", httpPort)
	s.gRPCAddress = fmt.Sprintf("0.0.0.0:%d", grpcPort)

	err := s.startGRPCGateway()
	should.EqualError(err, errText)
}

// rpc service entry
type greeterService struct{}

func (s *greeterService) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloReply, error) {
	// panic(111)
	// The panic simulated here can be automatically captured in the request
	// interceptor to record the operation log
	log.Println("req data: ", in)
	return &pb.HelloReply{
		Name:    "hello," + in.Name,
		Message: "call ok",
	}, nil
}

/** TestGRPCAndHttpServer
% go test -v -test.run=TestGRPCAndHttpServer
=== RUN   TestGRPCAndHttpServer
2020/07/12 22:26:54 Starting http server and grpc server listening on 8081
2020/07/12 22:26:55 exec begin
2020/07/12 22:26:55 client_ip: 127.0.0.1
2020/07/12 22:26:55 req data:  name:"daheige"
2020/07/12 22:26:55 exec end,cost time: 0 ms
2020/07/12 22:26:55 resp code:  200
2020/07/12 22:27:08 exec begin
2020/07/12 22:27:08 client_ip: 127.0.0.1
2020/07/12 22:27:08 req data:  name:"daheige123456"
2020/07/12 22:27:08 exec end,cost time: 0 ms
2020/07/12 22:27:15 exec begin
2020/07/12 22:27:15 client_ip: 127.0.0.1
2020/07/12 22:27:15 req data:  name:"daheige123456"
2020/07/12 22:27:15 exec end,cost time: 0 ms
2020/07/12 22:27:20 exec begin
2020/07/12 22:27:20 client_ip: 127.0.0.1
2020/07/12 22:27:20 req data:  name:"daheige123"
2020/07/12 22:27:20 exec end,cost time: 0 ms
2020/07/12 22:27:23 exec begin
2020/07/12 22:27:23 client_ip: 127.0.0.1
2020/07/12 22:27:23 req data:  name:"daheige"
2020/07/12 22:27:23 exec end,cost time: 0 ms
*/

func TestGRPCAndHttpServer(t *testing.T) {
	initConf()

	var should = require.New(t)

	// add the /test endpoint
	route := Route{
		Method:  "GET",
		Pattern: PathPattern("test"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte("Hello!"))
		},
	}

	// test Option func
	s := NewService(
		WithRouteOpt(route),
		WithShutdownFunc(shutdownFunc),
		WithPreShutdownDelay(2*time.Second),
		WithHandlerFromEndpoint(pb.RegisterGreeterServiceHandlerFromEndpoint),
		WithLogger(LoggerFunc(log.Printf)),
		WithRequestAccess(true),
		WithPrometheus(true),
		WithGRPCServerOption(grpc.ConnectionTimeout(10*time.Second)),
	)

	// register grpc service
	pb.RegisterGreeterServiceServer(s.GRPCServer, &greeterService{})

	newRoute := Route{
		Method:  "GET",
		Pattern: PathPattern("health"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute)

	newRoute2 := Route{
		Method:  "GET",
		Pattern: PathPattern("info"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute2)

	go func() {
		err := s.StartGRPCAndHTTPServer(sharePort)
		should.NoError(err)
	}()

	// wait 1 second for the server start
	time.Sleep(1 * time.Second)

	// check if the http server is up
	httpHost := fmt.Sprintf(":%d", sharePort)
	_, err := net.Listen("tcp", httpHost)
	should.Error(err)

	// check if the grpc server is up
	grpcHost := fmt.Sprintf(":%d", sharePort)
	_, err = net.Listen("tcp", grpcHost)
	should.Error(err)

	// check if the http endpoint works
	// Visit this address in the browser
	// http://localhost:8081/v1/say/daheige
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/v1/say/%s", sharePort, "daheige"))
	log.Println("resp code: ", resp.StatusCode)
	should.NoError(err)
	should.Equal(http.StatusOK, resp.StatusCode)

	resp, err = client.Get(fmt.Sprintf("http://127.0.0.1:%d/health", sharePort))
	should.NoError(err)
	should.Equal(http.StatusOK, resp.StatusCode)
	b, err := ioutil.ReadAll(resp.Body)
	should.NoError(err)
	should.Equal("OK", string(b))

	resp, err = client.Get(fmt.Sprintf("http://127.0.0.1:%d/info", sharePort))
	should.NoError(err)
	should.Equal(http.StatusOK, resp.StatusCode)

	resp, err = client.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics", sharePort))
	should.NoError(err)
	should.Equal(http.StatusOK, resp.StatusCode)
	time.Sleep(100 * time.Second)
}

/**
=== RUN   TestGRPCServerWithoutGateway
2020/07/12 21:58:17 Starting gPRC server listening on 9999
2020/07/12 21:58:18 exec begin
2020/07/12 21:58:18 client_ip: 127.0.0.1
2020/07/12 21:58:18 req data:  name:"daheige"
2020/07/12 21:58:18 exec end,cost time: 0 ms
2020/07/12 21:58:18 name:hello,daheige,message:call ok
2020/07/12 21:58:19 Waiting for 2s before shutdown starts
2020/07/12 21:58:31 gRPC server shutdown success
--- PASS: TestGRPCServerWithoutGateway (14.02s)
PASS
*/
func TestGRPCServerWithoutGateway(t *testing.T) {
	initConf()

	var should = require.New(t)

	// test Option func
	s := NewServiceWithoutGateway(
		WithShutdownFunc(shutdownFunc),
		WithPreShutdownDelay(2*time.Second),
		WithLogger(LoggerFunc(log.Printf)),
		WithRequestAccess(true),
		WithPrometheus(true),
		WithGRPCServerOption(grpc.ConnectionTimeout(10*time.Second)),
	)

	// register grpc service
	pb.RegisterGreeterServiceServer(s.GRPCServer, &greeterService{})

	go func() {
		err := s.StartGRPCWithoutGateway(grpcPort)
		should.NoError(err)
	}()

	// wait 1 second for the server start
	time.Sleep(1 * time.Second)

	// check if the grpc server is up
	grpcHost := fmt.Sprintf(":%d", grpcPort)
	_, err := net.Listen("tcp", grpcHost)
	should.Error(err)

	// Set up a connection to the server.
	conn, err := grpc.Dial(s.gRPCAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewGreeterServiceClient(conn)

	// Contact the server and print out its response.
	res, err := c.SayHello(context.Background(), &pb.HelloReq{
		Name: "daheige",
	})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("name:%s,message:%s", res.Name, res.Message)

	time.Sleep(1 * time.Second)
	s.StopGRPCWithoutGateway()
}

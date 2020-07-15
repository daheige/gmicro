package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/daheige/gmicro"
	"github.com/daheige/gmicro/example/pb"
	"google.golang.org/grpc"
)

var sharePort int
var shutdownFunc func()

func init() {
	sharePort = 8081

	shutdownFunc = func() {
		fmt.Println("Server shutting down")
	}
}

// http://localhost:8081/v1/say/daheige123
/**
% go run server.go
2020/06/27 23:25:43 Starting http server and grpc server listening on 8081
2020/06/27 23:25:53 exec begin
2020/06/27 23:25:53 client_ip: 127.0.0.1
2020/06/27 23:25:53 req data:  name:"daheige123"
2020/06/27 23:25:53 exec end,cost time: 0 ms
*/

func main() {
	// add the /test endpoint
	route := gmicro.Route{
		Method:  "GET",
		Pattern: gmicro.PathPattern("test"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte("Hello!"))
		},
	}

	// test Option func
	s := gmicro.NewService(
		gmicro.WithRouteOpt(route),
		gmicro.WithShutdownFunc(shutdownFunc),
		gmicro.WithPreShutdownDelay(2*time.Second),
		gmicro.WithHandlerFromEndpoint(pb.RegisterGreeterServiceHandlerFromEndpoint),
		gmicro.WithLogger(gmicro.LoggerFunc(log.Printf)),
		gmicro.WithRequestAccess(true),
		gmicro.WithPrometheus(true),
		gmicro.WithGRPCServerOption(grpc.ConnectionTimeout(10*time.Second)),
	)

	// register grpc service
	pb.RegisterGreeterServiceServer(s.GRPCServer, &greeterService{})

	newRoute := gmicro.Route{
		Method:  "GET",
		Pattern: gmicro.PathPattern("health"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute)

	newRoute2 := gmicro.Route{
		Method:  "GET",
		Pattern: gmicro.PathPattern("info"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute2)

	log.Fatalln(s.StartGRPCAndHTTPServer(sharePort))
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

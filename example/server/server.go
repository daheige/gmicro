package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/daheige/gmicro/v2"
	"github.com/daheige/gmicro/v2/example/pb"
)

var (
	sharePort    int
	shutdownFunc func()
)

func init() {
	sharePort = 8081
	shutdownFunc = func() {
		log.Println("Server shutting down")
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
		Method: "GET",
		Path:   "/test",
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte("Hello!"))
		},
	}

	// test Option func
	s := gmicro.NewService(
		gmicro.WithRouteOpt(route),
		gmicro.WithShutdownFunc(shutdownFunc),
		gmicro.WithPreShutdownDelay(2*time.Second),
		gmicro.WithShutdownTimeout(6*time.Second),
		gmicro.WithHandlerFromEndpoint(pb.RegisterGreeterServiceHandlerFromEndpoint),
		gmicro.WithLogger(gmicro.LoggerFunc(log.Printf)),
		gmicro.WithRequestAccess(true),
		gmicro.WithPrometheus(true),
		gmicro.WithGRPCServerOption(grpc.ConnectionTimeout(10*time.Second)),
		gmicro.WithGRPCNetwork("tcp"), // grpc server start network
		gmicro.WithStaticAccess(true), // enable static file access,if use http gw
	)

	// register grpc service
	pb.RegisterGreeterServiceServer(s.GRPCServer, &greeterService{})

	newRoute := gmicro.Route{
		Method: "GET",
		Path:   "health",
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute)

	newRoute2 := gmicro.Route{
		Method: "GET",
		Path:   "/info",
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute2)

	// you can start grpc server and http gateway on one port
	log.Fatalln(s.StartGRPCAndHTTPServer(sharePort))

	// you can also specify ports for grpc and http gw separately
	// log.Fatalln(s.Start(sharePort, 50051))

	// you can start server without http gateway
	// log.Fatalln(s.StartGRPCWithoutGateway(50051))
}

// rpc service entry
type greeterService struct {
	// 这里必须包含这个解构体才可以，否则就是没有实现
	/*
		// All implementations must embed UnimplementedGreeterServiceServer
		// for forward compatibility
		type GreeterServiceServer interface {
			SayHello(context.Context, *HelloReq) (*HelloReply, error)
			mustEmbedUnimplementedGreeterServiceServer()
		}
	*/
	pb.UnimplementedGreeterServiceServer
}

func (s *greeterService) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloReply, error) {
	// panic(111)
	// The panic simulated here can be automatically captured in the request
	// interceptor to record the operation log
	log.Println("req data: ", in)
	time.Sleep(12 * time.Millisecond)

	md := gmicro.GetIncomingMD(ctx)
	log.Println("request md: ", md)
	return &pb.HelloReply{
		Name:    "hello," + in.Name,
		Message: "call ok",
	}, nil
}

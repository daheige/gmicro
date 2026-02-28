package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/daheige/gmicro/v2"
	"github.com/daheige/gmicro/v2/example/clients/go/pb"
)

var (
	address = "localhost:8081" // grpc server and http gateway on share port
	// address = "localhost:50051" // grpc server port without http gateway
	// address     = "localhost:50050" // nginx grpc_pass port
	defaultName = "golang grpc"
)

/**
% go run client.go daheige
2026/02/28 11:12:08 x-request-id:  b790f23fc93743f8aafe1e09c8335e27
2026/02/28 11:12:08 name:hello,daheige,message:call ok
*/

func main() {
	// Set up a connection to the server.
	// please note the following settings
	// Deprecated: use WithTransportCredentials and insecure.NewCredentials()
	// instead. Will be supported throughout 1.x.
	// conn, err := grpc.Dial(address, grpc.WithInsecure())
	// so use grpc.WithTransportCredentials(insecure.NewCredentials()) as default grpc.DialOption
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(30*time.Minute), // 连接生命周期
		grpc.WithMaxCallAttempts(3),          // 最大重试次数
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewGreeterServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	requestID := gmicro.Uuid()
	log.Println("x-request-id: ", requestID)
	md := metadata.New(map[string]string{
		"x-request-id": requestID,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	res, err := c.SayHello(ctx, &pb.HelloReq{
		Name: name,
	})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("name:%s,message:%s", res.Name, res.Message)
}

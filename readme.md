# gmicro
  
    Golang grpc micro library.
    Microservice prototype with gRPC + http +h2c+ gRPC gateway + logger + prometheus.
    Require Go version >= v1.13.
    Reference project：https://github.com/dakalab/micro
   
# supported features

    Golang grpc microservice scaffolding mainly encapsulates some components of grpc,
    log tracking, link tracking, traffic h2c conversion, service monitoring prometheus
    and other functions,as far as possible to maintain the kiss principle, 
    so that these components are pluggable of.

    The framework supports 2 different methods such as http api, grpc server pb. 
    The protocol called by the client can be http or grpc pb format. 
    At the same time,it supports the same port, while providing 
    http services (supporting http1.x protocol requests) and the processing 
    capabilities of grpc server.

# installation 
  
  go get -v github.com/daheige/gmicro/v2
  
# quick start
    
    please see example demo

# grpc-go
    
https://github.com/grpc/grpc-go

# grpc-gateway doc

https://grpc-ecosystem.github.io/grpc-gateway/

# gmicro v1.2.x action

https://github.com/daheige/goapp

https://github.com/daheige/gmicro-demo

Old projects can continue to use the v1 version, and new projects can start to use the v2 version.
Note that the gmicro address has changed, and the v2 version is github.com/daheige/gmicro/v2

# change log

| options  | desc | time |
| :-----       | :---- |:----|
| grpc gateway | upgrade grpc gateway to v2.4.0 | 2021-06-19 |
| grpc         | upgrade go grpc to v1.38.0 | 2021-06-19 |
| protobuf     | upgrade protobuf to v1.26.0 | 2021-06-19 |

# grpc tools

    please do the following:
    # go gRPC tools
    go get -v \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
    
    # go validator
    go get github.com/go-playground/validator/v10
    
    #This will place four binaries in your $GOBIN;
    #    protoc-gen-grpc-gateway
    #    protoc-gen-openapiv2
    #    protoc-gen-go
    #    protoc-gen-go-grpc
    
    # protoc inject tag
    # go get -u github.com/favadi/protoc-go-inject-tag
    
# Google APIs

    # googleapis link:https://github.com/googleapis/googleapis
    api/rpc link: https://github.com/googleapis/googleapis/tree/master/google

    ============
    
    Project: Google APIs
    URL: https://github.com/google/googleapis
    Revision: 3544ab16c3342d790b00764251e348705991ea4b
    License: Apache License 2.0
    
    
    Imported Files
    ---------------
    
    - google/api/annotations.proto
    - google/api/http.proto
    - google/api/httpbody.proto
    
    
    Generated Files
    ----------------
    
    They are generated from the .proto files by protoc-gen-go.
    - google/api/annotations.pb.go
    - google/api/http.pb.go

# License

  MIT

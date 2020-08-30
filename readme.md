# gmicro
  
    Golang grpc micro library.
    Microservice prototype with gRPC + http +h2c+ logger + prometheus.
    Require Go version >= v1.13.
    Reference project：https://github.com/dakalab/micro

# supported features

    Golang grpc microservice scaffolding mainly encapsulates some components of grpc,
    log tracking, link tracking, traffic h2c conversion, service monitoring prometheus and other functions, 
    as far as possible to maintain the kiss principle, so that these components are pluggable of. 
    The framework supports 2 different methods such as http api, grpc server pb. 
    The protocol called by the client can be http or grpc pb format. At the same time, 
    it supports the same port, while providing http services (supporting http1.x protocol requests) 
    and the processing capabilities of grpc server.

# gmicro action

    see example or see https://github.com/daheige/goapp
    
# upgrade go grpc tools
  
    please perform the following operations in order
    # go gRPC tools
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
    go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
    
    # go validator
    go get github.com/go-playground/validator/v10
    
    # after successful execution, 3 binary files will be generated under the $GOBIN directory.
    # protoc-gen-grpc-gateway
    # protoc-gen-grpc-swagger
    # protoc-gen-go
    # google api link: github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
    
    # protoc inject tag
    # go get -u github.com/favadi/protoc-go-inject-tag  
    
# Google APIs

    https://github.com/grpc-ecosystem/grpc-gateway/tree/master/third_party/googleapis

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

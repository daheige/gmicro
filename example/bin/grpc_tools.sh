#!/bin/bash

# go gRPC tools
go get -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -v github.com/golang/protobuf/proto

go get -v google.golang.org/protobuf/cmd/protoc-gen-go

# if you can't use it, please use go get below
# go get -v github.com/golang/protobuf/protoc-gen-go

go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc

# go validator
go get -v github.com/go-playground/validator/v10

# after successful execution, 3 binary files will be generated under the $GOBIN directory.
# protoc-gen-grpc-gateway
# protoc-gen-grpc-swagger
# protoc-gen-go
# google api link: github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

# protoc inject tag
go get -v github.com/favadi/protoc-go-inject-tag

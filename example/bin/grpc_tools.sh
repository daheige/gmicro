#!/bin/bash

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

#!/bin/bash

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

# google api link:https://github.com/googleapis/googleapis

# protoc inject tag
# go get -u github.com/favadi/protoc-go-inject-tag

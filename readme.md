# gmicro
  
    Golang grpc micro library.
    Microservice prototype with gRPC + http +h2c+ gRPC gateway + logger + prometheus.
    Require Go version >= 1.23.0
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

# go version
    if you use go version < 1.16,please use gmicro tag v1.3.3
    else use gmicro tag v2.3.0 higher version.

# installation
  
  go get -v github.com/daheige/gmicro/v2
  
# quick start
    
    please see example demo

# grpc-go
    
https://github.com/grpc/grpc-go

# grpc-gateway doc

https://grpc-ecosystem.github.io/grpc-gateway/

# gmicro action

- https://github.com/daheige/athena
- https://github.com/daheige/goapp
- https://github.com/daheige/gmicro-demo

Old projects can continue to use the v1 version, and new projects can start to use the v2 version.
Note that the gmicro address has changed, and the v2 version is github.com/daheige/gmicro/v2

# change log
| options           | desc                                | time       |
|:------------------|:------------------------------------|:-----------|
| go mod update     | update grpc version to v1.70.0      | 2025-03-05 |
| fix shutdown func | fix bug                             | 2025-03-05 |
| go version        | Require Go version >= v1.23.0       | 2025-03-05 |

| options           | desc                                     | time       |
|:------------------|:-----------------------------------------|:-----------|
| go mod update     | update grpc version to v1.60.1           | 2024-01-10 |
| grpc nodejs tools | update grpc nodejs client code generated | 2024-01-10 |
| grpc mac tools    | add mac grpc install shell               | 2024-01-10 |

| options          | desc                                  | time       |
|:-----------------|:--------------------------------------|:-----------|
| go mod update    | update grpc mod                       | 2023-06-04 |
| grpc metadata    | fix request md for RequestInterceptor | 2023-06-04 |
| add x-request-id | fix request x-request-id from md      | 2023-06-04 |

| options         | desc                 | time       |
|:----------------|:---------------------|:-----------|
| grpc tools      | update docs          | 2022-05-15 |
| grpc dockerfile | update go grpc tools | 2022-05-15 |

| options         | desc                                      | time       |
|:----------------|:------------------------------------------|:-----------|
| grpc gateway    | upgrade grpc gateway to v2.10.0           | 2022-05-11 |
| grpc            | upgrade go grpc to v1.46.0                | 2022-05-11 |
| protobuf        | upgrade protobuf to v1.28.0               | 2022-05-11 |
| grpc dockerfile | upgrade go to v1.16.15 and protoc v3.15.8 | 2022-05-11 |

| options  | desc                              | time       |
| :-----       |:----------------------------------|:-----------|
| grpc gateway | upgrade grpc gateway to v2.4.0 | 2021-06-19 |
| grpc         | upgrade go grpc to v1.38.0 | 2021-06-19 |
| protobuf     | upgrade protobuf to v1.26.0 | 2021-06-19 |

# grpc tools dockerfile
https://github.com/daheige/gmicro-grpc-tools
gmicro v2版本对应的grpc v1.60.1 待更新。。。

# install grpc tools
```shell
sh example/bin/grpc_tools.sh
```

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

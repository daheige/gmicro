# gmicro
  
  Golang grpc micro library.
  Microservice prototype with gRPC + http +h2c+ logger + prometheus.
  Require Go version >= v1.11.
  Reference project：https://github.com/dakalab/micro

# supported features

  Golang grpc microservice scaffolding mainly encapsulates some components of grpc,
  log tracking, link tracking, traffic h2c conversion, service monitoring prometheus and other functions, 
  as far as possible to maintain the kiss principle, so that these components are pluggable of. 
  The framework supports 2 different methods such as http api, grpc server pb. 
  The protocol called by the client can be http or grpc pb format. At the same time, 
  it supports the same port, while providing http services (supporting http1.x protocol requests) 
  and the processing capabilities of grpc server.

# License

  MIT

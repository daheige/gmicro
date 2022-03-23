FROM golang:1.16.3

# build golang environment
ENV PROTOC_VER 3.15.6

RUN apt-get update && apt-get install -y apt-utils zip unzip;

RUN wget -q -P /tmp/temp/ https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip && \
    cd /usr && unzip /tmp/temp/protoc-${PROTOC_VER}-linux-x86_64.zip;


RUN go env -w GO111MODULE=on; go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct;go env -w CGO_ENABLED=0 &&  \
    go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u google.golang.org/grpc  && \
    go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc  && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway  && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    go get -u github.com/go-playground/validator/v10 && \
    go get  -u  github.com/golang/mock/mockgen  && \
    go get -u github.com/favadi/protoc-go-inject-tag

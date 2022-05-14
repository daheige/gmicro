FROM golang:1.16.3

# build golang environment
ENV PROTOC_VER 3.15.6

RUN apt-get update && apt-get install -y apt-utils zip unzip;

RUN wget -q -P /tmp/temp/ https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip && \
    cd /usr && unzip /tmp/temp/protoc-${PROTOC_VER}-linux-x86_64.zip;


RUN go env -w GO111MODULE=on; go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct;go env -w CGO_ENABLED=0 &&  \
    go get -u github.com/golang/protobuf/proto && \
    go install github.com/golang/protobuf/protoc-gen-go@latest && \
    go install google.golang.org/grpc@latest  && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest  && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest && \
    go get -u github.com/go-playground/validator/v10 && \
    go install github.com/golang/mock/mockgen  && \
    go install github.com/favadi/protoc-go-inject-tag@latest

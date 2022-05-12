FROM golang:1.16.15

# build golang environment
ENV PROTOC_VER 3.15.8

# 解决docker时区问题和中文乱码问题
ENV TZ=Asia/Shanghai LANG="zh_CN.UTF-8"

# 安装protoc工具
RUN apt-get update && apt-get install -y apt-utils zip unzip && \
    wget -q -P /tmp/temp/ https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip && \
    cd /usr && unzip /tmp/temp/protoc-${PROTOC_VER}-linux-x86_64.zip && \
    && mkdir -p /go/logs && mkdir /go/go-grpc && \
    echo "export LC_ALL=$LANG" >> /etc/profile && \
    echo $TZ > /etc/timezone

# 安装go grpc需要的工具包
RUN go env -w GO111MODULE=on; go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct;go env -w CGO_ENABLED=0 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest  && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest && \
    go get -u github.com/go-playground/validator/v10 && \
    go install github.com/golang/mock/mockgen@latest  && \
    go install github.com/favadi/protoc-go-inject-tag@latest

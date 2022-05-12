FROM golang:1.16.12

RUN apt-get update && apt-get -y install unzip && apt-get clean

# install protobuf
ENV PROTOC_VER 3.15.8
ENV PROTOC_URL https://github.com/google/protobuf/releases/download/v${PB_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip

# 解决docker时区问题和中文乱码问题
ENV TZ=Asia/Shanghai LANG="zh_CN.UTF-8"

# nodejs version
ENV NODE_VER v12.16.2

# 安装protoc工具
RUN mkdir -p /tmp/protoc && \
    curl -L ${PROTOC_URL} > /tmp/protoc/protoc.zip && \
    cd /tmp/protoc && \
    unzip protoc.zip && \
    cp /tmp/protoc/bin/protoc /usr/local/bin && \
    cp -R /tmp/protoc/include/* /usr/local/include && \
    chmod go+rx /usr/local/bin/protoc && \
    cd /tmp && \
    rm -r /tmp/protoc && \
    mkdir -p /go/logs && mkdir /go/go-grpc && \
    echo "export LC_ALL=$LANG" >> /etc/profile && \
    echo $TZ > /etc/timezone

# 设置golang环境变量和安装grpc工具
RUN go env -w GO111MODULE=on; go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct;go env -w CGO_ENABLED=0 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest  && \
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest && \
    go get -u github.com/go-playground/validator/v10 && \
    go install github.com/golang/mock/mockgen@latest  && \
    go install github.com/favadi/protoc-go-inject-tag@latest

# install nodejs
RUN cd /usr/local/ && \
     wget https://nodejs.org/dist/$NODE_VER/node-$NODE_VER-linux-x64.tar.xz && \
     xz -d node-$NODE_VER-linux-x64.tar.xz && \
     tar xvf node-$NODE_VER-linux-x64.tar.xz && \
     mv node-$NODE_VER-linux-x64 nodejs && \
     chown -R $USER /usr/local/nodejs && \
     ln -s /usr/local/nodejs/bin/npm /usr/bin/npm && \
     chmod +x /usr/bin/npm && \
     echo "export NODEJS_HOME=/usr/local/nodejs" >> ~/.bashrc && \
     echo 'export PATH=$NODEJS_HOME/bin:$PATH' >> ~/.bashrc && \
     source ~/.bashrc && \
     npm install -g cnpm --registry=https://registry.npm.taobao.org

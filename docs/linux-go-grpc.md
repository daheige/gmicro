# go v1.16.15安装

    cd /usr/local
    sudo wget https://golang.google.cn/dl/go1.16.15.linux-amd64.tar.gz
    sudo tar zxvf go1.16.15.linux-amd64.tar.gz
    sudo chown -R $USER /usr/local/go
    
    设置好go env环境变量
    go env -w GO111MODULE=on
    go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
    go env -w CGO_ENABLED=0

# centos7 protoc工具安装

    1、下载https://github.com/protocolbuffers/protobuf/archive/v3.15.8.tar.gz
        cd /usr/local/src
        sudo wget https://github.com/protocolbuffers/protobuf/archive/v3.15.8.tar.gz
    
    2、开始安装
        sudo mv v3.15.8.tar.gz protobuf-3.15.8.tar.gz
        sudo tar zxvf protobuf-3.15.8.tar.gz
        cd protobuf-3.15.8
        sudo yum install gcc-c++ cmake libtool
        # 对于ubuntu系统 sudo apt install gcc cmake make libtool
        $ sudo mkdir /usr/local/protobuf

        需要编译, 在新版的 PB 源码中，是不包含 .configure 文件的，需要生成
        此时先执行 sudo ./autogen.sh 
        脚本说明如下:
        # Run this script to generate the configure script and other files that will
        # be included in the distribution. These files are not checked in because they
        # are automatically generated.

        此时生成了 .configure 文件，可以开始编译了
        sudo ./configure --prefix=/usr/local/protobuf
        sudo make && make install

        安装完成后,查看版本:
        $ cd /usr/local/protobuf/bin
        $ ./protoc --version
        libprotoc 3.15.8
        
        建立软链接
        $ sudo ln -s /usr/local/protobuf/bin/protoc /usr/bin/protoc
        $ sudo chmod +x /usr/bin/protoc

# go protoc工具安装
    
    执行如下命令:
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
    go install github.com/go-playground/validator/v10@latest
    go install github.com/golang/mock/mockgen@latest
    go install github.com/favadi/protoc-go-inject-tag@latest

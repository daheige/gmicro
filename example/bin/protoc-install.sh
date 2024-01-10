# mac 源码安装protoc，当然你可以使用 mac-install-grpc.sh 快速安装
protoExec=$(which "protoc")
if [ -z $protoExec ]; then
    mkdir ~/web/
    cd ~/web/
    git clone https://github.com/google/protobuf.git
    cd protobuf
    sh ./autogen.sh
    ./configure
    make
    sudo make install
fi

# centos7 安装参考地址 https://github.com/daheige/rs-rpc?tab=readme-ov-file#centos7-install-protoc

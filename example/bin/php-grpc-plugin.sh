# php grpc工具安装
# grpc_php_plugin check.
grpc_php_plugin=$(which "grpc_php_plugin")
if [ -z $grpc_php_plugin ]; then
    mkdir ~/web/
    cd ~/web
    git clone -b v1.60.0 https://github.com/grpc/grpc
    cd ~/web/grpc
    git checkout -b v1
    git submodule update --init --recursive
    mkdir -p cmake/build
    cd cmake/build
    cmake ../..
    make
    make grpc_php_plugin
    make protoc grpc_php_plugin
    cp grpc_php_plugin /usr/local/bin/
    chmod +x /usr/local/bin/grpc_php_plugin

    # 安装php grpc拓展（推荐pecl方式安装）
    # pecl install grpc
    # 安装完成后，修改php.ini文件，增加配置 extension = grpc.so
else
    echo "grpc_php_plugin exists!"
fi
exit 0

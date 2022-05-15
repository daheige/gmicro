#!/usr/bin/env bash

#编译golang可执行文件
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

#docker image版本
version=$1
if [ -z $version ];then
    version=v1
fi

grpc_tools="go-grpc-tools"
cnt=`docker image ls | grep "$grpc_tools:$version" | wc -l`
if [ $cnt -gt 0 ];then
    #删除之前的image
    docker rmi -f $grpc_tools:$version
fi

#重新生成镜像
cd $root_dir
docker build -t $grpc_tools:$version .

echo "构建grpc tools image $grpc_tools:$version 成功"

exit 0

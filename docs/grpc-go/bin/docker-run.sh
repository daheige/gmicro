#!/usr/bin/env bash

root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

#docker数据卷映射到主机的目录
workdir=$HOME

#docker image版本
version=$1
if [ -z $version ];then
    version=v1
fi

# 构建镜像
# sh $root_dir/bin/docker-build.sh $version
# echo "构建image完成"

#docker容器name名称
containerName=$2
if [ -z $containerName ];then
    containerName=go-grpc-dev
fi

#创建docker映射到当前主机上的目录
mkdir -p $workdir/www/go-grpc
mkdir -p $workdir/logs/go-grpc
chmod 755 $workdir/logs/go-grpc

#停止之前的容器
cnt=`docker container ls -a | grep $containerName | grep -v grep | wc -l`
if [ $cnt -gt 0 ];then
    docker stop $containerName
    docker rm $containerName
fi

grpc_tools="go-grpc-tools"
#运行容器,如果需要运行端口，请自行修改这一行就可以
docker run -it -d --name $containerName -v $workdir/www/go-grpc:/go/go-grpc -v $workdir/logs:/go/logs $grpc_tools:$version

echo "go-grpc docker运行成功!"
echo "开始你的go-grpc应用之旅吧!"
echo "docker exec -it $containerName /bin/bash"

exit 0

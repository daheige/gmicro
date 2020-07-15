#!/usr/bin/env bash
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

protoExec=$(which "protoc")
if [ -z $protoExec ]; then
    echo 'Please install protoc!'
    echo "Please look readme.md to install proto3"
    exit 0
fi

mkdir -p $root_dir/clients/nodejs/pb

cd $root_dir/protos

$protoExec --js_out=import_style=commonjs,binary:$root_dir/clients/nodejs/pb/ --plugin=protoc-gen-grpc=/usr/local/grpc/bins/opt/grpc_node_plugin --grpc_out=$root_dir/clients/nodejs/pb/ *.proto

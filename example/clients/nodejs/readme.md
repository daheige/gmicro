# grpc nodejs package
```json
{
    "dependencies": {
        "google-protobuf": "^3.21.2",
        "@grpc/grpc-js": "^1.9.11",
        "grpc-tools": "^1.12.4"
    }
}
```

# yarn install
```shell
npm config set registry https://registry.npmmirror.com/ 
sudo npm install -g yarn
# 设置国内镜像
yarn config set registry https://registry.npmmirror.com/
```

# run nodejs

    First start the server server.go
    % cd gmicro/example/server
    % go run server.go
    2022/05/13 23:09:27 Starting http server and grpc server listening on 8081

    Second, start the client written in Node.js to make a gRPC request
    % yarn install
    % node hello.js
    {
        wrappers_: null,
        messageId_: undefined,
        arrayIndexOffset_: -1,
        array: [ 'hello,heige', 'call ok' ],
        pivot_: 1.7976931348623157e+308,
        convertedPrimitiveFields_: {}
    }
    call ok
    hello,heige
    
    the server will output:
    2022/05/13 23:10:02 exec begin
    2022/05/13 23:10:02 client_ip: 127.0.0.1
    2022/05/13 23:10:02 req data:  name:"heige"
    2022/05/13 23:10:02 exec end,cost time: 0 ms

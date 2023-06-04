# grpc nodejs package
```json
{
    "dependencies": {
        "google-protobuf": "^3.20.1",
        "grpc": "^1.24.11",
        "grpc-tools": "^1.11.2"
    }
}
```

# run nodejs

    First start the server server.go
    % cd gmicro/example/server
    % go run server.go
    2022/05/13 23:09:27 Starting http server and grpc server listening on 8081

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

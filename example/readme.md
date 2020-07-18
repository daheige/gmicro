# Code generation

    % sh bin/go-generate.sh
    
    Generating codes...
    
    generating golang stubs...
    
    generating golang code success
    
    Generate codes successfully!
    
    
    % sh bin/php-generate.sh
    
    Generating codes...
    
    generating php stubs...
    generating php stubs from: /web/go/gmicro/example/protos/hello.proto
    	[DONE]
    
    Generate codes successfully!

    % sh bin/nodejs-generate.sh
    
    Generating codes...
    
    generating nodejs stubs...
    generating nodejs code success
    
    Generate codes successfully!
    
# service run

    % go run server/server.go
    
    % go run clients/go/client.go
    2020/07/15 22:56:51 name:hello,golang grpc,message:call ok
    
    % cd clients/php
    % composer install
    Loading composer repositories with package information
    Installing dependencies (including require-dev) from lock file
    Nothing to install or update
    Generating autoload files
    
    % cd ../../
    % php clients/hello_client.php daheige
    check App\Grpc\GPBMetadata\Hello\HelloReq exist
    bool(true)
    status code: 0
    name:hello,daheige
    call ok

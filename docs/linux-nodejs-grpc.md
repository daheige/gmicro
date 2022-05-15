# install nodejs

    cd /usr/local/
    sudo wget https://nodejs.org/dist/v12.16.2/node-v12.16.2-linux-x64.tar.xz
    # 这里我用的ubuntu20.04系统，如果使用的是centos系统，请自行先安装好xz
    sudo apt-get install xz-utils
    sudo xz -d node-v12.16.2-linux-x64.tar.xz
    sudo tar xvf node-v12.16.2-linux-x64.tar
    sudo mv node-v12.16.2-linux-x64 nodejs
    
    # For the convenience of the current user, set nodejs to belong to the current user,
    # which is not recommended in the production environment

    sudo chown -R $USER /usr/local/nodejs
    sudo ln -s /usr/local/nodejs/bin/npm /usr/bin/npm
    sudo chmod +x /usr/bin/npm
    
    # add env path
    vim ~/.bashrc 
    export NODEJS_HOME=/usr/local/nodejs
    export PATH=$NODEJS_HOME/bin:$PATH
    
    :wq
    source ~/.bashrc

# install cnpm

    npm install -g cnpm --registry=https://registry.npm.taobao.org

# nodejs gprc install

    cnpm install grpc-tools -g
    cnpm install google-protobuf -g
    cnpm install grpc -g

# nodejs generate code

    sh bin/nodejs-generate.sh

# nodejs client run

    cd gmicro/example/clients/nodejs
    cnpm install
    node hello.js
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

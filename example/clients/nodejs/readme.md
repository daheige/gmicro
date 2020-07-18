# install nodejs

    cd /usr/local/
    sudo wget https://nodejs.org/dist/v12.16.2/node-v12.16.2-linux-x64.tar.xz
    sudo xz -d node-v12.16.2-linux-x64.tar.xz
    sudo tar xvf node-v12.16.2-linux-x64.tar.xz
    sudo mv node-v12.16.2-linux-x64 nodejs
    
    For the convenience of the current user, set nodejs to belong to the current user,
    which is not recommended in the production environment
    
    sudo chown -R $USER /usr/local/nodejs
    sudo ln -s /usr/local/nodejs/bin/npm /usr/bin/npm
    sudo chmod +x /usr/bin/npm
    vim ~/.bashrc 
    export NODEJS_HOME=/usr/local/nodejs
    export PATH=$NODEJS_HOME/bin:$PATH

    source ~/.bashrc

# install cnpm

    npm install -g cnpm --registry=https://registry.npm.taobao.org

# nodejs gprc install

    cnpm install grpc-tools --save-dev
    cnpm install google-protobuf --save
    cnpm install grpc --save

    or 
    
    cnpm install
     
# nodejs client run

    % node hello.js
    (node:52641) [DEP0005] DeprecationWarning: Buffer() is deprecated due to security and usability issues. Please use the Buffer.alloc(), Buffer.allocUnsafe(), or Buffer.from() methods instead.
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

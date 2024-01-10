# install php grpc and protobuf
```shell
pecl install grpc
pecl install protobuf
```
查看php.ini文件在哪里，执行如下命令就可以获得路径
```shell
php -i | grep "Loaded Configuration File"
```
比如我的就在/usr/local/etc/php/8.2/php.ini中，执行完毕后，就将拓展加入到php.ini中即可
```ini
extension=protobuf.so
extension=grpc.so
```

# php grpc run
     关于composer安装参考：example/bin/mac-php-grpc.sh

     composer install
     php hello_client.php

# About whether to use protobuf.so

     For php7.0+, protoc3 can install php protobuf extension
     vim php.ini
     ; It is not necessary to install, it is generally recommended to use protobuf to expand it is better
     extension=protobuf.so
     extension=grpc.so

     At this time, you can remove the composer2.json
     "google/protobuf": "^3.8"
     mv composer2.json composer.json
     Then composer update

     For those who do not support php protobuf expansion, you can replace composer2.json with composer.json
    
# composer mirror settings

     Use composer config -g repo.packagist composer https://mirrors.aliyun.com/composer/

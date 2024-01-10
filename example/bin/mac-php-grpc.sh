# mac protobuf安装
```shell
brew install protobuf
```
或者编译安装：
```shell
git clone https://github.com/google/protobuf.git

cd protobuf

sh ./autogen.sh

./configure

make
sudo make install

```
# php grpc_php_plugin 插件安装
参考文档：https://www.jianshu.com/p/bb15ad7532be

```shell
make protoc grpc_php_plugin
```

# php protobuf拓展安装
先判断是否安装php grpc.so and protobuf.so
```shell
php -m | grep proto
php -m | grep grpc
```
如果没有任何输出就需要安装grpc和protobuf拓展
```shell
pecl install grpc
pecl install protobuf
```
安装grpc拓展成功提示（这里我是php8.1版本）：
```
Build process completed successfully
Installing '/usr/local/Cellar/php@8.1/8.1.25/pecl/20210902/grpc.so'
install ok: channel://pecl.php.net/grpc-1.59.1
Extension grpc enabled in php.ini
```

安装protobuf拓展成功提示：
```
Build process completed successfully
Installing '/usr/local/Cellar/php@8.1/8.1.25/pecl/20210902/protobuf.so'
install ok: channel://pecl.php.net/protobuf-3.25.1
Extension protobuf enabled in php.ini
```

查看php.ini路径目录
```shell
php --ini
```
将输出如下内容：
```
% php --ini
Configuration File (php.ini) Path: /usr/local/etc/php/8.1
Loaded Configuration File:         /usr/local/etc/php/8.1/php.ini
Scan for additional .ini files in: /usr/local/etc/php/8.1/conf.d
Additional .ini files parsed:      /usr/local/etc/php/8.1/conf.d/ext-opcache.ini
```

安装好后，在php.ini中添加对应的配置即可
(一般来说上面安装后，会自动添加依赖，如果没有添加，请手动添加即可)
```ini
extension=protobuf.so
extension=grpc.so
```
查看拓展是否安装成功：
```shell
php -m | grep grpc
php -m | grep protobuf
```

# php composer安装
```shell
cd ~
php -r "copy('https://install.phpcomposer.com/installer', 'composer-setup.php');"
php composer-setup.php
php -r "unlink('composer-setup.php');"
```
将安装好的composer.phar移动到/usr/local/bin/下面即可
```shell
sudo mv composer.phar /usr/local/bin/composer
```
设置镜像：
```shell
composer config -g repo.packagist composer https://packagist.org
```

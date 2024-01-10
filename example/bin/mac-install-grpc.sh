# mac protoc安装
protocExec=`which protoc`
if [ ! -z $phpExec ]; then
    echo "you has installed protoc"
    exit 0
fi

brew install autoconf
brew install automake
brew install libtool
brew install curl
brew install make
brew install g++
brew install unzip
brew install cmake
brew install protobuf

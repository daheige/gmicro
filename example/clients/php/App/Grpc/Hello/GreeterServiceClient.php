<?php
// GENERATED CODE -- DO NOT EDIT!

namespace App\Grpc\Hello;

/**
 * service 定义开放调用的服务
 */
class GreeterServiceClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \App\Grpc\Hello\HelloReq $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     * @return \Grpc\UnaryCall
     */
    public function SayHello(\App\Grpc\Hello\HelloReq $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/App.Grpc.Hello.GreeterService/SayHello',
        $argument,
        ['\App\Grpc\Hello\HelloReply', 'decode'],
        $metadata, $options);
    }

}

# fake_push 工具介绍

## 目的

* 封住grpc中push协议服务端。
* 如果需要可以通过在响应的回调函中添加响应的逻辑，做到快速开发 。


##使用介绍
fake_push 提供如下参数：

 * etcd string ：pushd向etcd服务注册,etcd的地址。默认为："http://127.0.0.1:2379"
    	
 * register string ：回调服务的注册地址，即使本程序。默认"localhost:20081"

 * service string : 注册分服务名称，应该和connd配置相同。默认为"live-broadcast-bifrost-pushd"

 * group string : 注册服务组，应该和connd配置相同。默认为"live-broadcast"

 * region string : 注册服务区域。默认为""
## 使用方式
 
```
./fake_push -etcd="http://192.168.22.11:4445" -service="live-broadcast-bifrost-pushd"
```


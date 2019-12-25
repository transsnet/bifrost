# fake_conn 工具介绍

## 目的

* 封住grpc中connd协议服务端。
* 如果需要可以通过在响应的回调函中添加响应的逻辑，做到快速开发 。


##使用介绍
conn 提供如下参数：

 * etcd string ：conn向etcd服务注册,etcd的地址。默认为："http://127.0.0.1:2379"
    	
 * register string ：回调服务的注册地址，即使本程序。默认"localhost:20081"

 * service string : 注册分服务名称，应该和pushd配置相同。默认为"conn-live"


## 使用方式
 
```
./fake_conn -etcd="http://192.168.22.11:4445" -register="192.168.22.12:5041"
```


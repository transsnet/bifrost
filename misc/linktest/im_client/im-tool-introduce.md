# im工具介绍

## 目的

* 作为打桩程序，保持N个用户在线。
* 订阅相同的topic，统计最后一次消息达到时间。丢失的消息是什么，到达率多少。


##使用介绍
im 提供如下参数：

 * clientid string ； 订阅消息的clentid，默认前缀为fperf-clientid。实际默认发送clientid为fperf-clientid-0，fperf-clientid-1, 依此递增。
    	
 * count int ：保持多少客户端在线。默认情况为1个。
 
 * idc int ： 客户端尾缀起始递增的起点，默认为0
    	
 * exit ：在所有客户端都收到消息之后，true时所有客户端自动断开，进程退出。默认为true
    	
 * interval duration ：打印信息时间间隔,默认 2s。
    	
 * password string ：MQTT客户端连接服务的密码。默认为test
  
 * qos int ：MQTT客户端连接服务器的qos级别。默认订阅为1。
  
 * server string ：MQTT客户端服务器的地址。默认为tcp://127.0.0.1:8000"
    	
 * topic string ：MQTT客户端订阅的topic。
   
 * username string ：MQTT客户端连接用户名字。默认为test。

 * appkey string : bifrost 服务中的验证密码 默认为空

 * pubs string : publish 接口服务地址 默认"127.0.0.1:5053"

 * payload string : 发送给client的信息，如果为空默认每次发送信息的时间



## 常见使用场景
 
* N个客户端在线，消息下发一次，统计客户端的接收情况。

```
./im -server="ssl://192.168.22.12:1884" -topic="fperf/topic" -exit=true -qos=1 -count=1 -idc=0
```

* N个客户端在线，消息从客户端id=100开始下发，下发时间为10s，客户端消息接收情况。

```
./im -server="ssl://192.168.22.12:1884" -topic="fperf/topic" -exit=false -qos=1 -count=1 -idc=100  -keep=10s
```

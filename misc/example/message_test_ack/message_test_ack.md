# message_test_a2a 工具介绍

## 目的

* 测试离线消息是否可以到达。


##使用介绍
message_test_ack 提供如下参数：

 * a string ：设置MQTT中服务的地址。默认为："tcp://127.0.0.1:1883"

 * r bool ：设置MQTT中接收客户端是否自动重连。默认为false
  
 * rp string" ：设置MQTT中接收客户端连接的密码。默认为:"recevier-password"

 * ru string :设置MQTT中接收客户端的用户名字。默认为"recevier-username"

 * rid string ：设置MQTT中接收客户端的clientid。默认为："mqtt-receiver"

 * sp string" ：设置MQTT中发送客户端连接的密码。默认为:"sender-password"

 * su string :设置MQTT中发送客户端的用户名字。默认为"sender-username"

 * sid string ：设置MQTT中发送客户端的clientid。默认为："mqtt-sender"
   
## 使用方式
 
```
./fake_message_test_ack -c="192.168.22.12:23381"  -r=true 
```


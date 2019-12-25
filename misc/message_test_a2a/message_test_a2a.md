# message_test_a2a 工具介绍

## 目的

* 保持用户在线完成一次连接和消息的回复。


##使用介绍
message_test_a2a 提供如下参数：

 * c string ：设置MQTT中的clientid。默认为："ssl-sample-sender"
    	
 * h string ：设置MQTT中服务的地址。默认为："tcp://127.0.0.1:1883"
  
 * p string" ：设置MQTT中连接的密码。默认为:"password"
    	
 * r bool ：设置MQTT中客户端是否自动重连。默认为false

 * s bool ：设置MQTT中客户端是否清理绘画。默认为true
   
 * u string :设置MQTT中客户端的用户名字。默认为"username"

## 使用方式
 
```
./fake_message_test_a2a -c="192.168.22.12:23381" -s=false -r=true 
```


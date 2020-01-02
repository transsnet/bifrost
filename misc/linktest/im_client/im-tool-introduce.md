# Introduction to IM tools

## Purpose

* As a piling procedure, keep N users online which subscribe the same topic
* Print time the last message reached. the missing messages and the arrival rate


## Used to introduce
Provide the following parameters: 

 * clientid string : 订阅消息的clentid，默认前缀为fperf-clientid。实际默认发送clientid为fperf-clientid-0，fperf-clientid-1, 依此递增。
    	
 * count int : How many clients are kept online. The default is 1
 
 * idc int : The starting point at which the client end starts incrementing , The default is 0
    	
 * exit : When true, all clients automatically disconnect and the process exits after all clients receive the message. The default is true
    	
 * interval duration : Print message interval, the default is 2s
    	
 * password string : The password for the MQTT client connection service. The default is the test
  
 * qos int : The qos level of the MQTT client connection server. The default is 1.
  
 * server string : The address of the MQTT client server. The default is "tcp://127.0.0.1:8000"
    	
 * topic string : Topic to which MQTT clients subscribe.
   
 * username string : The MQTT client connects to the username. The default for the test.

 * appkey string : The authentication password in the service ,Th defaults is null

 * pubs string : Publish interface service address , the default is "127.0.0.1:5053".

 * payload string : The message sent to the client, if empty, defaults to the time each message is sent



## Common usage scenarios
 
* N clients are online, the message is sent once, and the reception of the client is counted

```
./im -server="ssl://192.168.22.12:1884" -topic="fperf/topic" -exit=true -qos=1 -count=1 -idc=0
```

* N clients are online. The message is sent from client id 100, and the time of sending is 10s. The client message is received.

```
./im -server="ssl://192.168.22.12:1884" -topic="fperf/topic" -exit=false -qos=1 -count=1 -idc=100  -keep=10s
```

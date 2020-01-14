# fake_conn tool is introduced

## Purpose

* Provide the connd protocol server in GRPC.

## Directions for use

conn parameters that：

 * etcd string : Conn registers with etcd service. default:"http://127.0.0.1:2379"
    	
 * register string : The registered address of the callback service ,default:"http://localhost:20081"

 * service string : Register the service name, which should be the same as the pushd configuration. default:"conn-live"


## Examplet
 
```
./fake_conn -etcd="http://192.168.22.11:4445" -register="192.168.22.12:5041"
```


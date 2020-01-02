## The MQTT protocol supports test automation

#### IM

    * c1.c2 login. subscribe to a topic. Send a message
    * C1 is offline. messages remain delivered as normal
    * message sent normally
    * c2 offline, message sent. c1,c2 did not receive the message

#### Push

    * c1.c2 login and subscribe no topic. Messages are issued normally

#### retain msg

    // Two case
    * c1.c2 is not logged in. send message retain empty message. c1.c2 is logged in and received.
    * c1.c2 not logged in. send message retain empty message. c1.c2login not received

#### cleansession

    * c1 log in to receive messages. Log in to receive messages, log out, log out, log in to receive messages

#### will msg 

    *  //TODO

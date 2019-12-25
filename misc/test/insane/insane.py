import imqtt
import argparse
import traceback

def connectN(host, port, n = 2):
    for i in range(n):
        c = imqtt.TCPClient(host, port)
        ack = c.Send(imqtt.ConnectPacket()).Recv()
        print(ack)

def disconn_before_conn(host, port):
    c = imqtt.TCPClient(host, port)
    c.Send(imqtt.DisconnectPacket())
    c.Send(imqtt.ConnectPacket()).Recv()

def sub_before_conn(host, port):
    c = imqtt.TCPClient(host, port)
    c.Send(imqtt.SubscribePacket()).Recv()
    
def pub_before_conn(host, port):
    c = imqtt.TCPClient(host, port)
    c.Send(imqtt.PublishPacket(Payload='hello', QoS=1)).Recv()

def sub_and_pub(host, port):
    # connect to server
    c1 = imqtt.TCPClient(host, port)
    c1.Send(imqtt.ConnectPacket(ClientID='imqtt1')).Recv()
    c2 = imqtt.TCPClient(host, port)
    c2.Send(imqtt.ConnectPacket(ClientID='imqtt2')).Recv()
    # sub and pub
    c1.Send(imqtt.SubscribePacket()).Recv()
    c2.Send(imqtt.PublishPacket(Payload='hello', QoS=1)).Recv()

    # recv the message
    print(c1.Recv())

def sub_after_disconn(host, port):
    c = imqtt.TCPClient(host, port)
    c.Send(imqtt.ConnectPacket()).Recv()
    c.Send(imqtt.DisconnectPacket())

    c.Send(imqtt.SubscribePacket()).Recv()

def pub_after_disconn(host, port):
    c = imqtt.TCPClient(host, port)
    c.Send(imqtt.ConnectPacket()).Recv()
    c.Send(imqtt.DisconnectPacket())

    c.Send(imqtt.PublishPacket(Payload='hello', QoS=1)).Recv()

def duplicate_clientid(host, port):
    c1 = imqtt.TCPClient(host, port)
    c1.Send(imqtt.ConnectPacket(ClientID='insane-client')).Recv()
    c2 = imqtt.TCPClient(host, port)
    c2.Send(imqtt.ConnectPacket(ClientID='insane-client')).Recv()

commands = {"connectN": connectN, 
        "disconn_before_conn": disconn_before_conn,
        "sub_before_conn": sub_before_conn,
        "pub_before_conn": pub_before_conn,
        "sub_and_pub": sub_and_pub,
        "sub_after_disconn": sub_after_disconn,
        "pub_after_disconn": pub_after_disconn,
        "duplicate_clientid": duplicate_clientid}

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='insane test client')
    parser.add_argument('--host') 
    parser.add_argument('--port', type=int) 
    parser.add_argument('testcase', nargs='+') 
    args = parser.parse_args()
    for tc in args.testcase:
        cmd = commands[tc]
        try:
            cmd(args.host, args.port)
        except Exception as e:
            traceback.print_exc()

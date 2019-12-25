package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/meitu/bifrost/misc/autotest/conndtest/callbacktest"
	"github.com/meitu/bifrost/misc/autotest/conndtest/mqtttest"
)

var (
	client  *mqtttest.Client
	callcli *callbacktest.Client
	//callback config
	service string
	group   string
	region  string
	//MQTT appkey
	username string
)

func main() {
	var (
		debug      bool
		pubsAddr   string
		etcdAddr   string
		mqttsAddr  string
		password   string
		testcaseid int
	)
	flag.IntVar(&testcaseid, "testcase", 0, "you need test caseid")
	flag.BoolVar(&debug, "debug", false, "display specific information")
	flag.StringVar(&pubsAddr, "pubs", "127.0.0.1:5053", "publish of address")
	flag.StringVar(&mqttsAddr, "mqtts", "tcp://127.0.0.1:1883", "address of mqtt broker")
	flag.StringVar(&etcdAddr, "etcd", "http://127.0.0.1:2379", "address of etcd")
	flag.StringVar(&username, "username", "test", "username")
	flag.StringVar(&password, "password", "test", "password")

	flag.StringVar(&service, "service", "service", "service callbck")
	flag.Parse()

	mqtttest.SetPassword(password)
	mqtttest.SetUsername(username)

	client = mqtttest.NewClient(pubsAddr, mqttsAddr)
	callcli = callbacktest.NewClient(pubsAddr, mqttsAddr, []string{etcdAddr})

	if debug {
		fmt.Printf("%#v\n", client)
	}
	run(testcaseid)
}

func run(testcaseid int) {
	switch testcaseid {
	case 1:
		if err := client.IMLiveCase(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("IM Live case sucess")
	case 2:
		if err := client.PushCase(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Push case sucess")
	case 3:
		if err := client.Retain(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Retain case sucess")
	case 4:
		if err := client.WillMsg(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Will case sucess")
	case 5:
		if err := client.CleanSession(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("cleansession case sucess")
	case 6:
		// if err := callcli.Service(service, group, region); err != nil {
		// log.Fatal(err)
		// }
		// fmt.Println("callback service sucess")
	case 7:
		// if err := callcli.CookieEuqal(service, group, region); err != nil {
		// log.Fatal(err)
		// }
		// fmt.Println("callback cookie sucess")
	default:
		if err := client.IMLiveCase(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("IM Live case sucess")
		time.Sleep(time.Millisecond * 100)
		if err := client.PushCase(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Push case sucess")
		time.Sleep(time.Millisecond * 100)
		if err := client.Retain(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Retain case sucess")
		time.Sleep(time.Millisecond * 100)
		if err := client.WillMsg(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Will case sucess")
		time.Sleep(time.Millisecond * 100)
		if err := client.CleanSession(); err != nil {
			log.Fatal(err)
		}
	}
}

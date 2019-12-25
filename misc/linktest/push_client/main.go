package main

import (
	"flag"
	"net/url"
	"strings"
	"time"

	"github.com/meitu/bifrost/misc/linktest/scene"
)

func main() {
	t := scene.NewThord()
	var appkey string
	var addrs string
	flag.StringVar(&addrs, "server", "tcp://127.0.0.1:1883", "address of mqtt broker")
	flag.StringVar(&appkey, "appkey", "", "appkey of bifrost server")
	flag.StringVar(&t.PubAddr, "pubs", "127.0.0.1:5053", "address of publish server")
	flag.StringVar(&t.Payload, "payload", "", "publish meessage")
	flag.DurationVar(&t.Interval, "interval", 2*time.Second, "the sending msg of interval time")
	flag.DurationVar(&t.KeepAlive, "keep", 0, "service life")

	flag.IntVar(&t.Cli.QoS, "qos", 1, "publish of qos which is range of 0-2")
	// flag.StringVar(&t.username, "username", "test", "username")
	flag.StringVar(&t.Cli.Password, "password", "test", "password")
	flag.StringVar(&t.Cli.Clientid, "clientid", "fperf-clientid", "prefix of clientid")
	flag.StringVar(&t.Cli.Topic, "topic", "/fperf/topic", "prefix of topic")
	flag.IntVar(&t.Cli.IDPoint, "id", 0, "the start of client point")

	flag.BoolVar(&t.Exit, "exit", false, "the server is exit after client receive all megssage")
	flag.IntVar(&t.Stat.Count, "count", 1, "the num of client")

	flag.IntVar(&t.Clis.DropCount, "drop", 0, "the num of client is closed in once")
	flag.Parse()

	t.Pubc = scene.NewPubcli(t.PubAddr, appkey)
	t.Clis.SetHandler(t.Cli.MqttConnect)
	t.Cli.Address = strings.Split(addrs, ",")
	t.Cli.SubSame = false
	app := make(url.Values)
	app.Set("bifrost-appkey", appkey)
	t.Cli.Username = app.Encode()
	t.Start()
}

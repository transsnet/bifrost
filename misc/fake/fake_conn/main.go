package main

import (
	"flag"

	"github.com/meitu/bifrost/misc/fake/fake_conn/fakeConnd"
)

type Settings struct {
	service string
	address string
	etcd    string
}

var settings Settings

func main() {
	flag.StringVar(&settings.address, "register", "localhost:20081", "register address")
	flag.StringVar(&settings.etcd, "etcd", "http://127.0.0.1:2379", "etcd address")
	flag.StringVar(&settings.service, "service", "connd-live", "service name of callback")
	flag.Parse()
	fakeConnd.DefaultConndServer([]string{settings.etcd}, settings.service, settings.address)
}

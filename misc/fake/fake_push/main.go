package main

import (
	"flag"
)

type Settings struct {
	service string
	group   string
	region  string
	address string
	etcd    string
}

var settings Settings

func main() {
	flag.StringVar(&settings.service, "service", "live-broadcast-bifrost-pushd", "service name")
	flag.StringVar(&settings.address, "register", "localhost:20081", "register address")
	flag.StringVar(&settings.etcd, "etcd", "http://127.0.0.1:2379", "etcd address")
	flag.StringVar(&settings.group, "group", "live-broadcast", "group")
	flag.StringVar(&settings.region, "region", "", "region")
	flag.Parse()
	// fake.DefaultPushdServer(settings.service, settings.address, settings.etcd, settings.group, settings.region)
}

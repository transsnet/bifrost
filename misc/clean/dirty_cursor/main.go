package main

import (
	"flag"
	"log"
	"time"

	"github.com/distributedio/configo"
)

type Conf struct {
	Redis     RedisConf
	Nap       time.Duration `cfg: "nap; 10ms;; sleep time between operating"`
	Cursor    string        `cfg: "cursor; 0;; start cursor(default) to scan"`
	Count     int           `cfg: "count; 50;; a hint for number(default) of keys to scan"`
	DelFields []string      `cfg:"fields; required; ; delete fields in key"`
}

type RedisConf struct {
	Servers []string `cfg:"servers; required; url; master cluster"`
	Auth    string   `cfg:"auth;;; password to connect to redis"`
	Timeout int64    `cfg:"timeout; 5000; >10; timeout of redis request in milliseconds"`
}

func main() {
	var path string
	var verbose bool
	flag.StringVar(&path, "c", "thor-gc.conf", "path of configration file")
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.Parse()

	conf := &Conf{}
	if err := configo.Load(path, conf); err != nil {
		log.Fatal(err)
	}

	if verbose {
		log.Println(conf)
	}

	gc := Gc{conf: conf}

	// var wg sync.WaitGroup
	for _, urlstring := range conf.Redis.Servers {
		// wg.Add(1)
		// go func(urlstring string) {
		gc.Purge(urlstring, conf.Redis.Auth)
		// wg.Done()
		// }(urlstring)
	}
	// wg.Wait()
}

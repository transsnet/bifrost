package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/meitu/bifrost/commons/auth"
)

func main() {
	var namespace string
	var key string

	flag.StringVar(&namespace, "n", "", "namespace")
	flag.StringVar(&key, "k", "", "server key")
	flag.Parse()

	token, err := auth.Token([]byte(key), []byte(namespace), time.Now().Unix())
	fmt.Println(string(token), err)
}

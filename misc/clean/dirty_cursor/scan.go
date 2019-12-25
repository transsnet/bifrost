package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Gc struct {
	conf *Conf
}

// Purge detele keys that has been expired
func (gc *Gc) Purge(urlstring, auth string) {
	u, err := url.Parse(urlstring)
	if err != nil {
		log.Fatal(err)
	}
	dialUrl := fmt.Sprintf("%s://%s:%s", u.Scheme, u.Hostname(), u.Port())

	cursor := gc.conf.Cursor
	count := gc.conf.Count

	conn, err := redis.DialURL(dialUrl, redis.DialPassword(auth))
	if err != nil {
		log.Println(err)
	}

	for {
		log.Printf("%s\n", dialUrl)

		cursor, err = gc.PurgeScan(conn, cursor, count)
		if err != nil {
			// try next redis
			log.Println(err)
			break
		}

		if cursor == "0" {
			log.Println(dialUrl, "scan done")
			time.Sleep(gc.conf.Nap)
			break
		}
	}
}

// PurgeOnce delete keys that has been expired by one scan
func (gc *Gc) PurgeScan(c redis.Conn, cursor string, count int) (string, error) {
	reply, err := redis.Values(c.Do("SCAN", cursor, "COUNT", count, "MATCH", "bifrost:route:callback:*"))
	if err != nil {
		return "0", err
	}

	// Get cursor from reply
	reply, err = redis.Scan(reply, &cursor)
	if err != nil {
		return "0", err
	}

	// Get keys from reply
	keys := make([]string, 0)
	reply, err = redis.Scan(reply, &keys)
	if err != nil {
		return cursor, err
	}

	cc := 0
	for _, key := range keys {
		// ignore bifrost keys
		for _, field := range gc.conf.DelFields {
			num, err := redis.Int(c.Do("HDEL", key, field))
			if err != nil {
				log.Println(err)
			}
			cc += num
		}
	}
	if cc != 0 {
		log.Printf("key count %d, clear count %d\n", len(keys), cc)
	}
	return cursor, nil
}

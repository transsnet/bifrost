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
	stat *statistics
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

	for _, key := range keys {
		fields, err := redis.StringMap(c.Do("HGETALL", key))
		if err != nil {
			log.Println(err)
		}
		if len(fields) > 1 {
			for k, _ := range fields {
				gc.stat.Add(k)
				// fmt.Println(k, v)
			}
			/*
				_, err := c.Do("del", key)
				if err != nil {
					fmt.Println(err)
				}
			*/
		}
	}
	return cursor, nil
}

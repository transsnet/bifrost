package main

import (
	"fmt"
	"sync"
)

type statistics struct {
	slock sync.Mutex
	smap  map[string]int
}

func NewStatistics() *statistics {
	return &statistics{
		smap: make(map[string]int),
	}
}

func (s *statistics) Add(ip string) {
	s.slock.Lock()
	s.smap[ip]++
	s.slock.Unlock()
}

func (s *statistics) Print() {
	maxk := ""
	maxv := 0
	sum := 0
	for k, v := range s.smap {
		if maxv < v {
			maxk = k
			maxv = v
		}
		sum += v
		fmt.Printf("key : %s ,fileds : %d\n", k, v)
	}
	fmt.Printf("the max number of node : %s ,count: %d\n", maxk, maxv)
	fmt.Printf("the sum of fileds %d\n", sum)
}

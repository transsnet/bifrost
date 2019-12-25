package scene

import (
	"log"
	"sync"
	"time"
)

type Statistics struct {
	Lock          sync.Mutex
	StartSendMsg  time.Time
	LatestRecvMsg time.Duration
	Count         int
	Total         int
}

func NewStat() *Statistics {
	return &Statistics{}
}

func StatDefault() *Statistics {
	stat := &Statistics{}
	return stat
}

func (stat *Statistics) Calculate() bool {
	stat.Lock.Lock()
	defer stat.Lock.Unlock()
	stat.Total++
	latest := time.Now().Sub(stat.StartSendMsg)
	if latest > stat.LatestRecvMsg {
		stat.LatestRecvMsg = latest
	}
	if stat.Count > stat.Total {
		return false
	}
	return true
}

func (stat *Statistics) Print(clean bool) {
	stat.Lock.Lock()
	if stat.Count != stat.Total {
		log.Printf("the expect count of recving messages %d\n", stat.Count)
		log.Printf("the real count of recving messages %d\n", stat.Total)
	}
	log.Printf("the latest recevied message time is %v\n", stat.LatestRecvMsg)
	log.Printf("the arrival rate is %v\n", 100*stat.Total/stat.Count)
	if clean {
		stat.Total = 0
		stat.LatestRecvMsg = 0
	}
	stat.Lock.Unlock()
}

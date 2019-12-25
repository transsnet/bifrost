package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/coreos/etcd/clientv3"
)

type GoddessOptions struct {
	EtcdCluster []string
	ServiceName string
	Nap         time.Duration
}

type Metric struct {
	Connection int64
	Publish    int64
}
type Metrics []*Metric

func (metrics Metrics) Sum() (int64, int64) {
	c := int64(0)
	p := int64(0)
	for i := range metrics {
		c += metrics[i].Connection
		p += metrics[i].Publish
	}
	return p, c
}

func main() {
	var opt GoddessOptions
	var cluster string
	flag.StringVar(&cluster, "etcd", "etcd01.bjbx.m.com:2379", "etcd cluster,seperate by comma")
	flag.StringVar(&opt.ServiceName, "service", "product-goddess-qa-connd", "service name of bifrost")
	flag.DurationVar(&opt.Nap, "nap", time.Second, "time duration between fetch")
	flag.Parse()

	opt.EtcdCluster = strings.Split(cluster, ",")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   opt.EtcdCluster,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	prefix := "/radar/service/" + opt.ServiceName + "/"
	ctx := context.Background()
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	var status = struct {
		Status string
	}{}

	connds := make([]string, 0, resp.Count)
	for i := 0; i < int(resp.Count); i++ {
		kv := resp.Kvs[i]
		node := path.Dir(string(kv.Key[len(prefix):]))
		key := path.Base(string(kv.Key))
		if key != "status" {
			continue
		}
		if err := json.Unmarshal(kv.Value, &status); err != nil {
			log.Fatalln(err)
		}
		if status.Status != "up" {
			continue
		}
		host, _, err := net.SplitHostPort(node)
		if err != nil {
			log.Println(err)
			continue
		}
		connds = append(connds, host)
	}

	begin := &Metric{}
	last := &Metric{}
	var bar *pb.ProgressBar
	for {
		time.Sleep(opt.Nap)

		pub, conn := fetchMetrics(connds).Sum()
		if begin.Publish == 0 {
			begin.Publish = pub
			begin.Connection = conn
			last.Connection = conn
			last.Publish = pub
			bar = pb.StartNew(int(conn))
			continue
		}
		if pub-begin.Publish > 0 {
			// bar.Set64(pub - begin.Publish)
		}
		if pub-last.Publish == 0 {
			begin.Publish = pub
			begin.Connection = conn
			bar.Finish()
			bar = pb.StartNew(int(conn))
		}
		last.Connection = conn
		last.Publish = pub
	}

}

func parse(line string) (key, val string) {
	key_begin := 0
	key_end := 0
	val_begin := 0
	val_end := 0

	for i := range line {
		c := line[i]
		if (c == ' ' || c == '{') && key_end == 0 {
			key_end = i
			key = line[key_begin:key_end]
			continue
		}
		if c == ' ' {
			val_begin = i + 1
			continue
		}

		if c == '\n' || c == '\r' {
			val_end = i
			val = line[val_begin:val_end]
		}
	}
	return
}

func fetchMetric(url string, out chan *Metric) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("fetch failed:", err)
		return
	}
	defer resp.Body.Close()

	metric := &Metric{}
	reader := bufio.NewReader(resp.Body)
	for {
		raw, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("read body failed:", err)
		}
		line := string(raw)
		if strings.Index(line, "bifrost_connd_down_packet_publish_count") >= 0 && !strings.HasPrefix(line, "#") {
			_, val := parse(line)
			ival, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				log.Println("parse int faile:", err)
				continue
			}
			metric.Publish += ival
		}
		if strings.Index(line, "bifrost_connd_current_connect_count") >= 0 && !strings.HasPrefix(line, "#") {
			if strings.Index(line, "ConnectFailure") >= 0 {
				continue
			}
			_, val := parse(line)
			ival, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				log.Println("parse int faile:", err)
				continue
			}
			metric.Connection += ival
		}
	}
	select {
	case out <- metric:
	default:
		log.Println("out chan blocked")
	}
}

func fetchMetrics(connds []string) Metrics {
	var wg sync.WaitGroup
	var metrics Metrics
	out := make(chan *Metric, 1024)
	for i := range connds {
		url := fmt.Sprintf("http://%s:2355/metrics", connds[i])
		wg.Add(1)
		go func() {
			fetchMetric(url, out)
			wg.Done()
		}()
	}
	wg.Wait()
	close(out)
loop:
	for {
		select {
		case m, ok := <-out:
			if !ok || m == nil {
				break loop
			}
			metrics = append(metrics, m)
		}
	}
	return metrics
}

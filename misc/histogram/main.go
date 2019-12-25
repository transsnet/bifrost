package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/benchmark/stats"
)

type HistogramOptions struct {
	NumBuckets     int
	GrowthFactor   float64
	BaseBucketSize float64
	MinValue       int64
}

func main() {
	opt := &HistogramOptions{}

	duration := false
	unit := "1"

	fs := flag.NewFlagSet("histogram", flag.ExitOnError)
	fs.IntVar(&opt.NumBuckets, "buckets", 16, "the number of buckets.")
	fs.Float64Var(&opt.GrowthFactor, "growth", 0.8, "growth factor of the buckets, A value of 0.1 indicates that bucket N+1 will be 10% larger than bucket N")

	fs.Float64Var(&opt.BaseBucketSize, "base", 100, "the size of the first bucket")
	fs.Int64Var(&opt.MinValue, "min", 10, "the lower bound of the first bucket")
	fs.BoolVar(&duration, "d", false, "parse value use time.Duration(golang)")
	fs.StringVar(&unit, "u", "1", "unit to dispaly histogram")
	fs.Parse(os.Args[1:])

	var u int64 = 1
	var d time.Duration
	var err error
	// try to parse unit as integer
	if c := unit[len(unit)-1]; c >= '0' && c <= '9' {
		u, err = strconv.ParseInt(unit, 10, 64)
	} else if duration {
		// parse as time.Duration
		d, err = time.ParseDuration("1" + unit)
		u = int64(d)
	} else {
		log.Fatalln("unexpected unit")
	}
	if err != nil {
		log.Fatalln(err)
	}

	histogram := stats.NewHistogram(stats.HistogramOptions{
		NumBuckets:     opt.NumBuckets,
		GrowthFactor:   opt.GrowthFactor,
		BaseBucketSize: opt.BaseBucketSize * float64(u),
		MinValue:       opt.MinValue * u,
	})

	var line string
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		var cost int64
		if duration {
			costDuration, err := time.ParseDuration(strings.TrimSpace(line))
			if err != nil {
				log.Println("parse failed", err)
				continue
			}
			cost = int64(costDuration)
		} else {
			cost, err = strconv.ParseInt(strings.TrimSpace(line), 10, 64)
			if err != nil {
				log.Println("parse failed", err)
				continue
			}
		}
		histogram.Add(cost)
	}
	histogram.PrintWithUnit(os.Stdout, float64(u))
}

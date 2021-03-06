// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package stats implements a generic benchmarking tool.
package stats

import (
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
	"time"
)

// Bar is true.
var Bar = true

// SetBar sets bar.
func SetBar(bar bool) {
	Bar = bar
}

// StartPrint prints the stats result.
func StartPrint(numParallels int, totalCalls int, clients []Client) {
	result := Start(numParallels, totalCalls, clients)
	fmt.Println(result.Format())
}

// Start returns the stats result.
func Start(numParallels int, totalCalls int, clients []Client) *Result {
	bodyChan := make(chan *body, totalCalls)
	startTime := time.Now()
	wg := &sync.WaitGroup{}
	numClients := len(clients)
	s := newStats(bodyChan, numClients, numParallels, totalCalls)
	cnt := &count{v: 0}
	for i := 0; i < numClients; i++ {
		go startClient(bodyChan, wg, numParallels, cnt, totalCalls, clients[i])
		wg.Add(1)
	}
	var stopLog = false
	var mut sync.Mutex
	if Bar {
		go func() {
			for {
				mut.Lock()
				if len(bodyChan) >= totalCalls || stopLog {
					mut.Unlock()
					break
				}
				mut.Unlock()
				i := int(cnt.load()) * 1e2 / totalCalls
				fmt.Fprintf(os.Stdout, "%d%% [%s]\r", i, getBar(i))
				time.Sleep(time.Millisecond * 1e2)
			}
		}()
	}
	wg.Wait()
	s.setTime(time.Now().Sub(startTime).Nanoseconds() / 1e3)
	mut.Lock()
	stopLog = true
	mut.Unlock()
	if Bar {
		fmt.Fprintf(os.Stdout, "%s\r", getStr(106, " "))
	}
	<-s.finish
	return s.result()
}

type stats struct {
	totalCalls        int
	finish            chan bool
	bodyChan          chan *body
	Clients           int
	Parallels         int
	Time              float64
	TotalTime         float64
	Times             []int
	TotalRequestSize  int64
	TotalResponseSize int64
	ResponseOk        int64
	Errors            int64
}

// Result is a stats result.
type Result struct {
	Clients                    int
	Parallels                  int
	TotalCalls                 int64
	TotalTime                  float64
	RequestsPerSecond          float64
	AverageTimePerRequest      float64
	FastestTimeForRequest      float64
	SlowestTimeForRequest      float64
	N001thThousandthTime       float64
	N010thThousandthTime       float64
	N050thThousandthTime       float64
	N100thThousandthTime       float64
	N250thThousandthTime       float64
	N500thThousandthTime       float64
	N750thThousandthTime       float64
	N900thThousandthTime       float64
	N950thThousandthTime       float64
	N990thThousandthTime       float64
	N999thThousandthTime       float64
	TotalRequestBodySizes      int64
	AverageBodySizePerRequest  float64
	RequestRateBytePerSecond   float64
	RequestRateMBytePerSecond  float64
	TotalResponseBodySizes     int64
	AverageBodySizePerResponse float64
	ResponseRateBytePerSecond  float64
	ResponseRateMBytePerSecond float64
	ResponseOk                 int64
	ResponseOkPercentile       float64
	Errors                     int64
	ErrorsPercentile           float64
}

func newStats(bodyChan chan *body, numClients int, numParallels int, totalCalls int) *stats {
	s := &stats{
		finish:     make(chan bool, 1),
		totalCalls: totalCalls,
		bodyChan:   bodyChan,
		Clients:    numClients,
		Parallels:  numParallels,
		Times:      make([]int, totalCalls),
	}
	go s.run()
	return s
}

func (s *stats) setTime(time int64) {
	s.Time = float64(time)
}

func (s *stats) run() {
	i := 0
	for body := range s.bodyChan {
		s.Times[i] = int(body.Time)
		i++
		s.TotalTime += float64(body.Time)
		s.TotalRequestSize += body.RequestSize
		s.TotalResponseSize += body.ResponseSize
		if body.Error {
			s.Errors++
		} else {
			s.ResponseOk++
		}
		bodyPool.Put(body)
		if i == s.totalCalls {
			break
		}
	}
	s.finish <- true
}

func (s *stats) result() *Result {
	sort.Ints(s.Times)
	total := float64(len(s.Times))
	totalInt := int64(total)
	var statsResult = &Result{}
	statsResult.Clients = s.Clients
	statsResult.Parallels = s.Parallels
	statsResult.TotalCalls = totalInt
	statsResult.TotalTime = s.Time / 1e6
	statsResult.RequestsPerSecond = total / (s.Time / 1e6)
	statsResult.AverageTimePerRequest = s.TotalTime / total / 1e3
	statsResult.FastestTimeForRequest = float64(s.Times[0]) / 1e3
	statsResult.SlowestTimeForRequest = float64(s.Times[totalInt-1]) / 1e3
	statsResult.N001thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*1))-1]) / 1e3
	statsResult.N010thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*10))-1]) / 1e3
	statsResult.N050thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*50))-1]) / 1e3
	statsResult.N100thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*100))-1]) / 1e3
	statsResult.N250thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*250))-1]) / 1e3
	statsResult.N500thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*500))-1]) / 1e3
	statsResult.N750thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*750))-1]) / 1e3
	statsResult.N900thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*900))-1]) / 1e3
	statsResult.N950thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*950))-1]) / 1e3
	statsResult.N990thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*990))-1]) / 1e3
	statsResult.N999thThousandthTime = float64(s.Times[int(math.Ceil(total/1e3*999))-1]) / 1e3
	statsResult.ResponseOk = s.ResponseOk
	statsResult.ResponseOkPercentile = float64(s.ResponseOk) / total * 1e2
	statsResult.Errors = s.Errors
	statsResult.ErrorsPercentile = float64(s.Errors) / total * 1e2
	if s.TotalRequestSize > 0 {
		statsResult.TotalRequestBodySizes = s.TotalRequestSize
		statsResult.AverageBodySizePerRequest = float64(s.TotalRequestSize) / total
		tr := float64(s.TotalRequestSize) / (s.Time / 1e6)
		statsResult.RequestRateBytePerSecond = tr
		statsResult.RequestRateMBytePerSecond = tr / 1e6
	}
	if s.TotalResponseSize > 0 {
		statsResult.TotalResponseBodySizes = s.TotalResponseSize
		statsResult.AverageBodySizePerResponse = float64(s.TotalResponseSize) / float64(s.ResponseOk)
		tr := float64(s.TotalResponseSize) / (s.Time / 1e6)
		statsResult.ResponseRateBytePerSecond = tr
		statsResult.ResponseRateMBytePerSecond = tr / 1e6
	}
	return statsResult
}

// Format returns the formatted string.
func (statsResult *Result) Format() string {
	format := ""
	format += fmt.Sprintln("Summary:")
	format += fmt.Sprintf("\tClients:\t%d\n", statsResult.Clients)
	format += fmt.Sprintf("\tParallel calls per client:\t%d\n", statsResult.Parallels)
	format += fmt.Sprintf("\tTotal calls:\t%d\n", statsResult.TotalCalls)
	format += fmt.Sprintf("\tTotal time:\t%.2fs\n", statsResult.TotalTime)
	format += fmt.Sprintf("\tRequests per second:\t%.2f\n", statsResult.RequestsPerSecond)
	format += fmt.Sprintf("\tFastest time for request:\t%.2fms\n", statsResult.FastestTimeForRequest)
	format += fmt.Sprintf("\tAverage time per request:\t%.2fms\n", statsResult.AverageTimePerRequest)
	format += fmt.Sprintf("\tSlowest time for request:\t%.2fms\n\n", statsResult.SlowestTimeForRequest)
	format += fmt.Sprintln("Time:")
	format += fmt.Sprintf("\t0.1%%\ttime for request:\t%.2fms\n", statsResult.N001thThousandthTime)
	format += fmt.Sprintf("\t1%%\ttime for request:\t%.2fms\n", statsResult.N010thThousandthTime)
	format += fmt.Sprintf("\t5%%\ttime for request:\t%.2fms\n", statsResult.N050thThousandthTime)
	format += fmt.Sprintf("\t10%%\ttime for request:\t%.2fms\n", statsResult.N100thThousandthTime)
	format += fmt.Sprintf("\t25%%\ttime for request:\t%.2fms\n", statsResult.N250thThousandthTime)
	format += fmt.Sprintf("\t50%%\ttime for request:\t%.2fms\n", statsResult.N500thThousandthTime)
	format += fmt.Sprintf("\t75%%\ttime for request:\t%.2fms\n", statsResult.N750thThousandthTime)
	format += fmt.Sprintf("\t90%%\ttime for request:\t%.2fms\n", statsResult.N900thThousandthTime)
	format += fmt.Sprintf("\t95%%\ttime for request:\t%.2fms\n", statsResult.N950thThousandthTime)
	format += fmt.Sprintf("\t99%%\ttime for request:\t%.2fms\n", statsResult.N990thThousandthTime)
	format += fmt.Sprintf("\t99.9%%\ttime for request:\t%.2fms\n\n", statsResult.N999thThousandthTime)
	if statsResult.TotalRequestBodySizes > 0 {
		format += fmt.Sprintln("Request:")
		format += fmt.Sprintf("\tTotal request body sizes:\t%d\n", statsResult.TotalRequestBodySizes)
		format += fmt.Sprintf("\tAverage body size per request:\t%.2f Byte\n", statsResult.AverageBodySizePerRequest)
		format += fmt.Sprintf("\tRequest rate per second:\t%.2f Byte/s (%.2f MByte/s)\n\n", statsResult.RequestRateBytePerSecond, statsResult.RequestRateMBytePerSecond)
	}
	if statsResult.TotalResponseBodySizes > 0 {
		format += fmt.Sprintln("Response:")
		format += fmt.Sprintf("\tTotal response body sizes:\t%d\n", statsResult.TotalResponseBodySizes)
		format += fmt.Sprintf("\tAverage body size per response:\t%.2f Byte\n", statsResult.AverageBodySizePerResponse)
		format += fmt.Sprintf("\tResponse rate per second:\t%.2f Byte/s (%.2f MByte/s)\n\n", statsResult.ResponseRateBytePerSecond, statsResult.ResponseRateMBytePerSecond)
	}
	format += fmt.Sprintln("Result:")
	format += fmt.Sprintf("\tResponse ok:\t%d (%.2f%%)\n", statsResult.ResponseOk, statsResult.ResponseOkPercentile)
	format += fmt.Sprintf("\tErrors:\t%d (%.2f%%)\n", statsResult.Errors, statsResult.ErrorsPercentile)
	return format
}

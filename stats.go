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
	"sync/atomic"
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
	numClients := len(clients)
	bodyChan := make(chan []*body, int(float64(numParallels*numClients)*1.25))
	bodyChans := make([]chan *body, numParallels*numClients)
	for i := 0; i < numParallels*numClients; i++ {
		ch := make(chan *body, int(float64(totalCalls)*1.5)/(numParallels*numClients))
		bodyChans[i] = ch
	}
	cwg := &sync.WaitGroup{}
	s := newStats(bodyChan, numClients, numParallels, totalCalls)
	cnt := &count{v: 0}
	start := make(chan struct{})
	for i := 0; i < numClients; i++ {
		go startClient(bodyChans[i*numParallels:(i+1)*numParallels], start, cwg, numParallels, cnt, totalCalls, clients[i])
		cwg.Add(1)
	}
	wg := &sync.WaitGroup{}
	var stopLog int32
	if Bar {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for {
				if atomic.LoadInt32(&stopLog) > 0 {
					break
				}
				i := int(cnt.load()) * 1e2 / totalCalls
				fmt.Fprintf(os.Stdout, "%d%% [%s]\r", i, getBar(i))
				time.Sleep(time.Millisecond * 1e2)
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		for {
			for i := 0; i < numParallels*numClients; i++ {
				ch := bodyChans[i]
				select {
				case b, ok := <-ch:
					if !ok {
						return
					}
					length := len(ch)
					bodies := make([]*body, 0, length+1)
					bodies = append(bodies, b)
					for j := 0; j < length; j++ {
						b := <-ch
						bodies = append(bodies, b)
					}
					bodyChan <- bodies
				default:
				}
			}
			sum := 0
			for i := 0; i < numParallels*numClients; i++ {
				sum += len(bodyChans[i])
			}
			if sum == 0 {
				time.Sleep(time.Millisecond)
			}
		}
	}()
	startTime := time.Now()
	close(start)
	cwg.Wait()
	s.setTime(time.Now().Sub(startTime).Nanoseconds())
	atomic.StoreInt32(&stopLog, 1)
	if Bar {
		fmt.Fprintf(os.Stdout, "%s\r", getStr(106, " "))
	}
	<-s.finish
	for i := 0; i < numParallels*numClients; i++ {
		close(bodyChans[i])
	}
	wg.Wait()
	return s.result()
}

type stats struct {
	totalCalls        int
	finish            chan bool
	bodyChan          chan []*body
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
	N000001thMillionthTime     float64
	N000010thMillionthTime     float64
	N000100thMillionthTime     float64
	N001000thMillionthTime     float64
	N010000thMillionthTime     float64
	N050000thMillionthTime     float64
	N100000thMillionthTime     float64
	N250000thMillionthTime     float64
	N500000thMillionthTime     float64
	N750000thMillionthTime     float64
	N900000thMillionthTime     float64
	N950000thMillionthTime     float64
	N990000thMillionthTime     float64
	N999000thMillionthTime     float64
	N999900thMillionthTime     float64
	N999990thMillionthTime     float64
	N999999thMillionthTime     float64
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

func newStats(bodyChan chan []*body, numClients int, numParallels int, totalCalls int) *stats {
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
	for bodies := range s.bodyChan {
		for idx := range bodies {
			body := bodies[idx]
			s.Times[i] = int(body.Time)
			i++
			s.TotalTime += float64(body.Time)
			s.TotalRequestSize += body.RequestSize
			if body.Error {
				s.Errors++
			} else {
				s.TotalResponseSize += body.ResponseSize
				s.ResponseOk++
			}
			bodyPool.Put(body)
		}
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
	statsResult.TotalTime = s.Time / 1e9
	statsResult.RequestsPerSecond = total / (s.Time / 1e9)
	statsResult.AverageTimePerRequest = s.TotalTime / total / 1e6
	statsResult.FastestTimeForRequest = float64(s.Times[0]) / 1e6
	statsResult.SlowestTimeForRequest = float64(s.Times[totalInt-1]) / 1e6
	statsResult.N000001thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*1))-1]) / 1e6
	statsResult.N000010thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*10))-1]) / 1e6
	statsResult.N000100thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*100))-1]) / 1e6
	statsResult.N001000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*1000))-1]) / 1e6
	statsResult.N010000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*10000))-1]) / 1e6
	statsResult.N050000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*50000))-1]) / 1e6
	statsResult.N100000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*100000))-1]) / 1e6
	statsResult.N250000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*250000))-1]) / 1e6
	statsResult.N500000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*500000))-1]) / 1e6
	statsResult.N750000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*750000))-1]) / 1e6
	statsResult.N900000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*900000))-1]) / 1e6
	statsResult.N950000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*950000))-1]) / 1e6
	statsResult.N990000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*990000))-1]) / 1e6
	statsResult.N999000thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*999000))-1]) / 1e6
	statsResult.N999900thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*999900))-1]) / 1e6
	statsResult.N999990thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*999990))-1]) / 1e6
	statsResult.N999999thMillionthTime = float64(s.Times[int(math.Ceil(total/1e6*999999))-1]) / 1e6
	statsResult.ResponseOk = s.ResponseOk
	statsResult.ResponseOkPercentile = float64(s.ResponseOk) / total * 1e2
	statsResult.Errors = s.Errors
	statsResult.ErrorsPercentile = float64(s.Errors) / total * 1e2
	if s.TotalRequestSize > 0 {
		statsResult.TotalRequestBodySizes = s.TotalRequestSize
		statsResult.AverageBodySizePerRequest = float64(s.TotalRequestSize) / total
		tr := float64(s.TotalRequestSize) / (s.Time / 1e9)
		statsResult.RequestRateBytePerSecond = tr
		statsResult.RequestRateMBytePerSecond = tr / 1e6
	}
	if s.TotalResponseSize > 0 {
		statsResult.TotalResponseBodySizes = s.TotalResponseSize
		statsResult.AverageBodySizePerResponse = float64(s.TotalResponseSize) / float64(s.ResponseOk)
		tr := float64(s.TotalResponseSize) / (s.Time / 1e9)
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
	format += fmt.Sprintf("\tTotal time:\t%.3fs\n", statsResult.TotalTime)
	format += fmt.Sprintf("\tRequests per second:\t%.3f\n", statsResult.RequestsPerSecond)
	format += fmt.Sprintf("\tFastest time for request:\t%.6fms\n", statsResult.FastestTimeForRequest)
	format += fmt.Sprintf("\tAverage time per request:\t%.6fms\n", statsResult.AverageTimePerRequest)
	format += fmt.Sprintf("\tSlowest time for request:\t%.6fms\n\n", statsResult.SlowestTimeForRequest)
	format += fmt.Sprintln("Time:")
	format += fmt.Sprintf("\t00.0001%%\ttime for request:\t%.6fms\n", statsResult.N000001thMillionthTime)
	format += fmt.Sprintf("\t00.0010%%\ttime for request:\t%.6fms\n", statsResult.N000010thMillionthTime)
	format += fmt.Sprintf("\t00.0100%%\ttime for request:\t%.6fms\n", statsResult.N000100thMillionthTime)
	format += fmt.Sprintf("\t00.1000%%\ttime for request:\t%.6fms\n", statsResult.N001000thMillionthTime)
	format += fmt.Sprintf("\t01.0000%%\ttime for request:\t%.6fms\n", statsResult.N010000thMillionthTime)
	format += fmt.Sprintf("\t05.0000%%\ttime for request:\t%.6fms\n", statsResult.N050000thMillionthTime)
	format += fmt.Sprintf("\t10.0000%%\ttime for request:\t%.6fms\n", statsResult.N100000thMillionthTime)
	format += fmt.Sprintf("\t25.0000%%\ttime for request:\t%.6fms\n", statsResult.N250000thMillionthTime)
	format += fmt.Sprintf("\t50.0000%%\ttime for request:\t%.6fms\n", statsResult.N500000thMillionthTime)
	format += fmt.Sprintf("\t75.0000%%\ttime for request:\t%.6fms\n", statsResult.N750000thMillionthTime)
	format += fmt.Sprintf("\t90.0000%%\ttime for request:\t%.6fms\n", statsResult.N900000thMillionthTime)
	format += fmt.Sprintf("\t95.0000%%\ttime for request:\t%.6fms\n", statsResult.N950000thMillionthTime)
	format += fmt.Sprintf("\t99.0000%%\ttime for request:\t%.6fms\n", statsResult.N990000thMillionthTime)
	format += fmt.Sprintf("\t99.9000%%\ttime for request:\t%.6fms\n", statsResult.N999000thMillionthTime)
	format += fmt.Sprintf("\t99.9900%%\ttime for request:\t%.6fms\n", statsResult.N999900thMillionthTime)
	format += fmt.Sprintf("\t99.9990%%\ttime for request:\t%.6fms\n", statsResult.N999990thMillionthTime)
	format += fmt.Sprintf("\t99.9999%%\ttime for request:\t%.6fms\n\n", statsResult.N999999thMillionthTime)

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
	format += fmt.Sprintf("\tResponse ok:\t%d (%.3f%%)\n", statsResult.ResponseOk, statsResult.ResponseOkPercentile)
	format += fmt.Sprintf("\tErrors:\t%d (%.3f%%)\n", statsResult.Errors, statsResult.ErrorsPercentile)
	return format
}

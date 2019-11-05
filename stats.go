package stats

import (
	"fmt"
	"sort"
	"math"
)
var Log =true

func SetLog(log bool)  {
	Log=log
}

type Stats struct {
	totalCalls 			int
	finish 				chan bool
	bodyChan 			chan *Body
	Conns				int
	Parallels			int
	Time				float64
	TotalTime			float64
	Times				[]int
	TotalRequestSize	int64
	TotalResponseSize	int64
	ResponseOk			int64
	Errors				int64
}
func newStats(bodyChan chan *Body,conns int,parallels int,totalCalls int)*Stats{
	stats := &Stats{
		finish:			make(chan bool,1),
		totalCalls:		totalCalls,
		bodyChan:		bodyChan,
		Conns:			conns,
		Parallels:		parallels,
		Times:			make([]int, totalCalls),
	}
	go stats.run()
	return stats
}
func (stats *Stats)SetTime(time int64){
	stats.Time=float64(time)
}
func (stats *Stats)run(){
	i := 0
	for res := range stats.bodyChan {
		stats.Times[i] = int(res.Time)
		i++
		stats.TotalTime += float64(res.Time)
		stats.TotalRequestSize += res.RequestSize
		stats.TotalResponseSize += res.ResponseSize
		if res.Error {
			stats.Errors++
		}else {
			stats.ResponseOk++
		}
		if i==stats.totalCalls{
			break
		}
	}
	stats.finish<-true
}
func (stats *Stats)Result()*StatsResult{
	sort.Ints(stats.Times)
	total := float64(len(stats.Times))
	totalInt := int64(total)
	var statsResult =&StatsResult{}
	statsResult.Conns=stats.Conns
	statsResult.Parallels=stats.Parallels
	statsResult.TotalCalls=totalInt
	statsResult.TotalTime=stats.Time/1E6
	statsResult.RequestsPerSecond=total/(stats.Time/1E6)
	statsResult.AverageTimePerRequest=stats.TotalTime/total/1000
	statsResult.FastestTimeForRequest=float64(stats.Times[0])/1000
	statsResult.SlowestTimeForRequest=float64(stats.Times[totalInt-1])/1000
	statsResult.N001thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*1))-1])/1000
	statsResult.N010thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*10))-1])/1000
	statsResult.N050thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*50))-1])/1000
	statsResult.N100thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*100))-1])/1000
	statsResult.N250thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*250))-1])/1000
	statsResult.N500thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*500))-1])/1000
	statsResult.N750thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*750))-1])/1000
	statsResult.N900thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*900))-1])/1000
	statsResult.N950thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*950))-1])/1000
	statsResult.N990thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*990))-1])/1000
	statsResult.N999thThousandthTime=float64(stats.Times[int(math.Ceil(total/1000*999))-1])/1000
	statsResult.ResponseOk=stats.ResponseOk
	statsResult.ResponseOkPercentile=float64(stats.ResponseOk)/total*1e2
	statsResult.Errors=stats.Errors
	statsResult.ErrorsPercentile=float64(stats.Errors)/total*1e2
	if stats.TotalRequestSize>0{
		statsResult.TotalRequestBodySizes=stats.TotalRequestSize
		statsResult.AverageBodySizePerRequest=float64(stats.TotalRequestSize)/total
		tr := float64(stats.TotalRequestSize) / (stats.Time / 1E6)
		statsResult.RequestRateBytePerSecond=tr
		statsResult.RequestRateMBytePerSecond=tr/1E6
	}
	if stats.TotalResponseSize>0{
		statsResult.TotalResponseBodySizes=stats.TotalResponseSize
		statsResult.AverageBodySizePerResponse=float64(stats.TotalResponseSize)/total
		tr := float64(stats.TotalResponseSize) / (stats.Time / 1E6)
		statsResult.ResponseRateBytePerSecond=tr
		statsResult.ResponseRateMBytePerSecond=tr/1E6
	}
	return statsResult
}
type StatsResult struct {
	Conns 					int
	Parallels     			int
	TotalCalls				int64
	TotalTime				float64
	RequestsPerSecond 		float64
	AverageTimePerRequest		float64
	FastestTimeForRequest 	float64
	N001thThousandthTime		float64
	N010thThousandthTime		float64
	N050thThousandthTime		float64
	N100thThousandthTime		float64
	N250thThousandthTime		float64
	N500thThousandthTime		float64
	N750thThousandthTime		float64
	N900thThousandthTime		float64
	N950thThousandthTime		float64
	N990thThousandthTime		float64
	N999thThousandthTime	float64
	SlowestTimeForRequest 	float64
	TotalRequestBodySizes  int64
	AverageBodySizePerRequest  float64
	RequestRateBytePerSecond  float64
	RequestRateMBytePerSecond  float64
	TotalResponseBodySizes  int64
	AverageBodySizePerResponse  float64
	ResponseRateBytePerSecond  float64
	ResponseRateMBytePerSecond  float64
	ResponseOk				int64
	ResponseOkPercentile	float64
	Errors      			int64
	ErrorsPercentile		float64
}
func (statsResult *StatsResult)Format() {
	fmt.Println("Summary:")
	fmt.Printf("\tConns:%d\n", statsResult.Conns)
	fmt.Printf("\tParallels:%d\n", statsResult.Parallels)
	fmt.Printf("\tTotal Calls:%d\n", statsResult.TotalCalls)
	fmt.Printf("\tTotal time:%.2fs\n", statsResult.TotalTime)
	fmt.Printf("\tRequests per second:%.2f\n", statsResult.RequestsPerSecond)
	fmt.Printf("\tFastest time for request:%.2fms\n", statsResult.FastestTimeForRequest)
	fmt.Printf("\tAverage time per request:%.2fms\n", statsResult.AverageTimePerRequest)
	fmt.Printf("\tSlowest time for request:%.2fms\n\n", statsResult.SlowestTimeForRequest)
	fmt.Println("Time:")
	fmt.Printf("\t0.1%%\ttime for request:%.2fms\n", statsResult.N001thThousandthTime)
	fmt.Printf("\t1%%\t\ttime for request:%.2fms\n", statsResult.N010thThousandthTime)
	fmt.Printf("\t5%%\t\ttime for request:%.2fms\n", statsResult.N050thThousandthTime)
	fmt.Printf("\t10%%\t\ttime for request:%.2fms\n", statsResult.N100thThousandthTime)
	fmt.Printf("\t25%%\t\ttime for request:%.2fms\n", statsResult.N250thThousandthTime)
	fmt.Printf("\t50%%\t\ttime for request:%.2fms\n", statsResult.N500thThousandthTime)
	fmt.Printf("\t75%%\t\ttime for request:%.2fms\n", statsResult.N750thThousandthTime)
	fmt.Printf("\t90%%\t\ttime for request:%.2fms\n", statsResult.N900thThousandthTime)
	fmt.Printf("\t95%%\t\ttime for request:%.2fms\n", statsResult.N950thThousandthTime)
	fmt.Printf("\t99%%\t\ttime for request:%.2fms\n", statsResult.N990thThousandthTime)
	fmt.Printf("\t99.9%%\ttime for request:%.2fms\n\n", statsResult.N999thThousandthTime)
	if statsResult.TotalRequestBodySizes>0{
		fmt.Println("Request:")
		fmt.Printf("\tTotal request body sizes:%d\n", statsResult.TotalRequestBodySizes)
		fmt.Printf("\tAverage body size per request:%.2f Byte\n", statsResult.AverageBodySizePerRequest)
		fmt.Printf("\tRequest rate per second:%.2f Byte/s (%.2f MByte/s)\n\n", statsResult.RequestRateBytePerSecond,statsResult.RequestRateMBytePerSecond)
	}

	if statsResult.TotalResponseBodySizes>0{
		fmt.Println("Response:")
		fmt.Printf("\tTotal response body sizes:%d\n", statsResult.TotalResponseBodySizes)
		fmt.Printf("\tAverage body size per response:%.2f Byte\n", statsResult.AverageBodySizePerResponse)
		fmt.Printf("\tResponse rate per second:%.2f Byte/s (%.2f MByte/s)\n\n", statsResult.ResponseRateBytePerSecond,statsResult.ResponseRateMBytePerSecond)
	}
	fmt.Println("Result:")
	fmt.Printf("\tResponseOk:%d (%.2f%%)\n", statsResult.ResponseOk, statsResult.ResponseOkPercentile)
	fmt.Printf("\tErrors:%d (%.2f%%)\n", statsResult.Errors, statsResult.ErrorsPercentile)
}
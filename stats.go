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

func (statsResult *StatsResult)Format()string{
	format:=""
	format+=fmt.Sprintln("Summary:")
	format+=fmt.Sprintf("\tConns:\t%d\n", statsResult.Conns)
	format+=fmt.Sprintf("\tParallels:\t%d\n", statsResult.Parallels)
	format+=fmt.Sprintf("\tTotal Calls:\t%d\n", statsResult.TotalCalls)
	format+=fmt.Sprintf("\tTotal time:\t%.2fs\n", statsResult.TotalTime)
	format+=fmt.Sprintf("\tRequests per second:\t%.2f\n", statsResult.RequestsPerSecond)
	format+=fmt.Sprintf("\tFastest time for request:\t%.2fms\n", statsResult.FastestTimeForRequest)
	format+=fmt.Sprintf("\tAverage time per request:\t%.2fms\n", statsResult.AverageTimePerRequest)
	format+=fmt.Sprintf("\tSlowest time for request:\t%.2fms\n\n", statsResult.SlowestTimeForRequest)
	format+=fmt.Sprintln("Time:")
	format+=fmt.Sprintf("\t0.1%%\ttime for request:\t%.2fms\n", statsResult.N001thThousandthTime)
	format+=fmt.Sprintf("\t1%%\ttime for request:\t%.2fms\n", statsResult.N010thThousandthTime)
	format+=fmt.Sprintf("\t5%%\ttime for request:\t%.2fms\n", statsResult.N050thThousandthTime)
	format+=fmt.Sprintf("\t10%%\ttime for request:\t%.2fms\n", statsResult.N100thThousandthTime)
	format+=fmt.Sprintf("\t25%%\ttime for request:\t%.2fms\n", statsResult.N250thThousandthTime)
	format+=fmt.Sprintf("\t50%%\ttime for request:\t%.2fms\n", statsResult.N500thThousandthTime)
	format+=fmt.Sprintf("\t75%%\ttime for request:\t%.2fms\n", statsResult.N750thThousandthTime)
	format+=fmt.Sprintf("\t90%%\ttime for request:\t%.2fms\n", statsResult.N900thThousandthTime)
	format+=fmt.Sprintf("\t95%%\ttime for request:\t%.2fms\n", statsResult.N950thThousandthTime)
	format+=fmt.Sprintf("\t99%%\ttime for request:\t%.2fms\n", statsResult.N990thThousandthTime)
	format+=fmt.Sprintf("\t99.9%%\ttime for request:\t%.2fms\n\n", statsResult.N999thThousandthTime)
	if statsResult.TotalRequestBodySizes>0{
		format+=fmt.Sprintln("Request:")
		format+=fmt.Sprintf("\tTotal request body sizes:\t%d\n", statsResult.TotalRequestBodySizes)
		format+=fmt.Sprintf("\tAverage body size per request:\t%.2f Byte\n", statsResult.AverageBodySizePerRequest)
		format+=fmt.Sprintf("\tRequest rate per second:\t%.2f Byte/s (%.2f MByte/s)\n\n", statsResult.RequestRateBytePerSecond,statsResult.RequestRateMBytePerSecond)
	}

	if statsResult.TotalResponseBodySizes>0{
		format+=fmt.Sprintln("Response:")
		format+=fmt.Sprintf("\tTotal response body sizes:\t%d\n", statsResult.TotalResponseBodySizes)
		format+=fmt.Sprintf("\tAverage body size per response:\t%.2f Byte\n", statsResult.AverageBodySizePerResponse)
		format+=fmt.Sprintf("\tResponse rate per second:\t%.2f Byte/s (%.2f MByte/s)\n\n", statsResult.ResponseRateBytePerSecond,statsResult.ResponseRateMBytePerSecond)
	}
	format+=fmt.Sprintln("Result:")
	format+=fmt.Sprintf("\tResponseOk:\t%d (%.2f%%)\n", statsResult.ResponseOk, statsResult.ResponseOkPercentile)
	format+=fmt.Sprintf("\tErrors\t:%d (%.2f%%)\n", statsResult.Errors, statsResult.ErrorsPercentile)
	return format
}
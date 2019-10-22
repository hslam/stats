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
	Connections 		int
	Parallels     		int
	AvgDuration 		float64
	Duration    		float64
	Sum         		float64
	Times       		[]int
	TransferedRequest	int64
	TransferedResponse	int64
	ResponseOk			int64
	Errors      		int64
}

func getS(n int,char string) (s string) {
	if n<1{
		return
	}
	for i:=1;i<=n;i++{
		s+=char
	}
	return
}

func CalcStats(bodyChan chan *Body, duration int64,numConnections int,numParallels int) *Stats {
	stats := &Stats{
		Connections: numConnections,
		Parallels:     numParallels,
		Times:       make([]int, len(bodyChan)),
		Duration:    float64(duration),
		AvgDuration: float64(duration),
	}

	i := 0
	for res := range bodyChan {
		stats.Sum += float64(res.Duration)
		stats.Times[i] = int(res.Duration)
		i++
		stats.TransferedRequest += res.RequestSize
		stats.TransferedResponse += res.ResponseSize
		if res.Error {
			stats.Errors++
		}else {
			stats.ResponseOk++
		}

		if len(bodyChan) == 0 {
			break
		}
	}
	return stats
}

func CalcStatsResult(allStats *Stats) *StatsResult{
	sort.Ints(allStats.Times)
	total := float64(len(allStats.Times))
	totalInt := int64(total)
	var statsResult =&StatsResult{}
	statsResult.Connections=allStats.Connections
	statsResult.Parallels=allStats.Parallels
	statsResult.TotalCalls=totalInt
	statsResult.TotalTimePassed=allStats.AvgDuration/1E6
	statsResult.RequestsPerSecond=total/(allStats.AvgDuration/1E6)
	statsResult.FastestTimeForRequest=float64(allStats.Times[0])/1000
	statsResult.AvgTimePerRequest=allStats.Sum/total/1000
	statsResult.MedianTimePerRequest=float64(allStats.Times[totalInt/2-1])/1000
	statsResult.N10thPercentileTime=float64(allStats.Times[int(math.Ceil(total/10*1))-1])/1000
	statsResult.N25thPercentileTime=float64(allStats.Times[int(math.Ceil(total/100*25))-1])/1000
	statsResult.N75thPercentileTime=float64(allStats.Times[int(math.Ceil(total/100*75))-1])/1000
	statsResult.N90thPercentileTime=float64(allStats.Times[int(math.Ceil(total/10*9))-1])/1000
	statsResult.N99thPercentileTime=float64(allStats.Times[int(math.Ceil(total/100*99))-1])/1000
	statsResult.N999thThousandthTime=float64(allStats.Times[int(math.Ceil(total/1000*999))-1])/1000
	statsResult.SlowestTimeForRequest=float64(allStats.Times[totalInt-1])/1000
	statsResult.ResponseOk=allStats.ResponseOk
	statsResult.ResponseOkPercentile=float64(allStats.ResponseOk)/total*1e2
	statsResult.Errors=allStats.Errors
	statsResult.ErrorsPercentile=float64(allStats.Errors)/total*1e2
	if allStats.TransferedRequest>0{
		statsResult.TotalRequestBodySizes=allStats.TransferedRequest
		statsResult.AvgBodySizePerRequest=float64(allStats.TransferedRequest)/total
		tr := float64(allStats.TransferedRequest) / (allStats.AvgDuration / 1E6)
		statsResult.TransferRequestRateBytePerSecond=tr
		statsResult.TransferRequestRateMBytePerSecond=tr/1E6
	}
	if allStats.TransferedResponse>0{
		statsResult.TotalResponseBodySizes=allStats.TransferedResponse
		statsResult.AvgBodySizePerResponse=float64(allStats.TransferedResponse)/total
		tr := float64(allStats.TransferedResponse) / (allStats.AvgDuration / 1E6)
		statsResult.TransferResponseRateBytePerSecond=tr
		statsResult.TransferResponseRateMBytePerSecond=tr/1E6
	}
	return statsResult
}
type StatsResult struct {
	Connections 			int
	Parallels     			int
	TotalCalls				int64
	TotalTimePassed			float64
	RequestsPerSecond 		float64
	FastestTimeForRequest 	float64
	AvgTimePerRequest		float64
	N10thPercentileTime		float64
	N25thPercentileTime		float64
	MedianTimePerRequest	float64
	N75thPercentileTime		float64
	N90thPercentileTime		float64
	N99thPercentileTime		float64
	N999thThousandthTime	float64
	SlowestTimeForRequest 	float64
	TotalRequestBodySizes  int64
	AvgBodySizePerRequest  float64
	TransferRequestRateBytePerSecond  float64
	TransferRequestRateMBytePerSecond  float64
	TotalResponseBodySizes  int64
	AvgBodySizePerResponse  float64
	TransferResponseRateBytePerSecond  float64
	TransferResponseRateMBytePerSecond  float64
	ResponseOk				int64
	ResponseOkPercentile	float64
	Errors      			int64
	ErrorsPercentile		float64
}
func PrintStatsResult(statsResult *StatsResult) {
	fmt.Println("==========================BENCHMARK==========================")
	fmt.Printf("Used Connections:\t\t\t%d\n", statsResult.Connections)
	fmt.Printf("Concurrent Calls Per Connection:\t%d\n", statsResult.Parallels)
	fmt.Printf("Total Number Of Calls:\t\t\t%d\n\n", statsResult.TotalCalls)
	fmt.Println("===========================TIMINGS===========================")
	fmt.Printf("Total time passed:\t\t\t%.2fs\n", statsResult.TotalTimePassed)
	fmt.Printf("Requests per second:\t\t\t%.2f\n", statsResult.RequestsPerSecond)
	fmt.Printf("Avg time per request:\t\t\t%.2fms\n", statsResult.AvgTimePerRequest)
	fmt.Printf("Fastest time for request:\t\t%.2fms\n", statsResult.FastestTimeForRequest)
	fmt.Printf("10%% time:\t\t\t\t\t\t%.2fms\n", statsResult.N10thPercentileTime)
	fmt.Printf("25%% time:\t\t\t\t\t\t%.2fms\n", statsResult.N25thPercentileTime)
	fmt.Printf("Median time per request:\t\t%.2fms\n", statsResult.MedianTimePerRequest)
	fmt.Printf("75%% time:\t\t\t\t\t\t%.2fms\n", statsResult.N75thPercentileTime)
	fmt.Printf("90%% time:\t\t\t\t\t\t%.2fms\n", statsResult.N90thPercentileTime)
	fmt.Printf("99%% time:\t\t\t\t\t\t%.2fms\n", statsResult.N99thPercentileTime)
	fmt.Printf("99.9%% time:\t\t\t\t\t\t%.2fms\n", statsResult.N999thThousandthTime)
	fmt.Printf("Slowest time for request:\t\t%.2fms\n\n", statsResult.SlowestTimeForRequest)
	if statsResult.TotalRequestBodySizes>0{
		fmt.Println("=========================REQUESTDATA=========================")
		fmt.Printf("Total request body sizes:\t\t%d\n", statsResult.TotalRequestBodySizes)
		fmt.Printf("Avg body size per request:\t\t%.2f Byte\n", statsResult.AvgBodySizePerRequest)
		fmt.Printf("Transfer request rate per second:\t%.2f Byte/s (%.2f MByte/s)\n\n", statsResult.TransferRequestRateBytePerSecond,statsResult.TransferRequestRateMBytePerSecond)
	}
	if statsResult.TotalResponseBodySizes>0{
		fmt.Println("=========================RESPONSEDATA========================")
		fmt.Printf("Total response body sizes:\t\t%d\n", statsResult.TotalResponseBodySizes)
		fmt.Printf("Avg body size per response:\t\t%.2f Byte\n", statsResult.AvgBodySizePerResponse)
		fmt.Printf("Transfer response rate per second:\t%.2f Byte/s (%.2f MByte/s)\n\n", statsResult.TransferResponseRateBytePerSecond,statsResult.TransferResponseRateMBytePerSecond)
	}
	fmt.Println("==========================RESPONSES==========================")
	fmt.Printf("ResponseOk:\t\t\t\t%d (%.2f%%)\n", statsResult.ResponseOk, statsResult.ResponseOkPercentile)
	fmt.Printf("Errors:\t\t\t\t\t%d (%.2f%%)\n", statsResult.Errors, statsResult.ErrorsPercentile)
}
package stats

import (
	"sync"
	"os"
	"time"
	"fmt"
	"encoding/json"
)
type Client interface {
	Call()(RequestSize int64,ResponseSize int64,ok bool)
}
func StartClientStats(numParallels int,totalCalls int, clients []Client) []byte {
	bodyChan := make(chan *Body, totalCalls*2)
	benchTime := NewTimer()
	benchTime.Reset()
	wg := &sync.WaitGroup{}
	numConnections:=len(clients)
	for i := 0; i < numConnections; i++ {
		go startClient(
			bodyChan,
			wg,
			numParallels,
			totalCalls,
			clients[i],
		)
		wg.Add(1)
	}
	var stopLog=false
	if Log{
		go func() {
			for {
				if len(bodyChan) >= totalCalls||stopLog{
					break
				}
				i:=len(bodyChan)*100/totalCalls
				fmt.Fprintf(os.Stdout, "%d%% [%s]\r",i,getS(i,"#") + getS(100-i," "))
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}
	wg.Wait()
	stopLog=true
	if Log{
		fmt.Fprintf(os.Stdout, "%s\r",getS(106," "))
	}
	stats := CalcStats(
		bodyChan,
		benchTime.Duration(),
		numConnections,
		numParallels,
	)
	statsResult:=CalcStatsResult(stats)
	PrintStatsResult(statsResult)
	b, err := json.Marshal(&stats)
	if err != nil {
		fmt.Println(err)
	}
	return b
}
func startClient(bodyChan chan *Body, waitGroup *sync.WaitGroup, numParallels int,totalCalls int,c Client) {
	defer waitGroup.Done()
	wg := &sync.WaitGroup{}
	for i:=0;i<numParallels;i++{
		go run(bodyChan,wg,totalCalls,c)
		wg.Add(1)
	}
	wg.Wait()
}

func run (bodyChan chan *Body,waitGroup *sync.WaitGroup,totalCalls int, c Client){
	defer waitGroup.Done()
	timer := NewTimer()
	for {
		timer.Reset()
		body := &Body{}
		RequestSize,ResponseSize,ok:=c.Call()
		body.Error=!ok
		body.RequestSize=RequestSize
		body.ResponseSize=ResponseSize
		body.Duration = timer.Duration()
		if len(bodyChan) >= totalCalls {
			break
		}
		bodyChan <- body
	}
}
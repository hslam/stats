package stats

import (
	"sync"
	"os"
	"time"
	"fmt"
)
type Client interface {
	Call()(RequestSize int64,ResponseSize int64,ok bool)
}

func getStr(n int,char string) (s string) {
	if n<1{
		return
	}
	for i:=1;i<=n;i++{
		s+=char
	}
	return
}

func Start(parallels int,totalCalls int, clients []Client) {
	bodyChan := make(chan *Body, totalCalls)
	startTime := time.Now()
	wg := &sync.WaitGroup{}
	conns:=len(clients)
	stats := newStats(bodyChan, conns, parallels,totalCalls)
	count:=&Count{v:0}
	for i := 0; i < conns; i++ {
		go startClient(bodyChan, wg, parallels, count, totalCalls, clients[i])
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
				fmt.Fprintf(os.Stdout, "%d%% [%s]\r",i,getStr(i,"#") + getStr(100-i," "))
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}
	wg.Wait()
	stopLog=true
	if Log{
		fmt.Fprintf(os.Stdout, "%s\r",getStr(106," "))
	}
	stats.SetTime(time.Now().Sub(startTime).Nanoseconds()/1000)
	<-stats.finish
	statsResult:=stats.Result()
	statsResult.Format()
}
func startClient(bodyChan chan *Body, waitGroup *sync.WaitGroup, numParallels int,count *Count,totalCalls int,c Client) {
	defer waitGroup.Done()
	wg := &sync.WaitGroup{}
	for i:=0;i<numParallels;i++{
		go run(bodyChan,wg,count,totalCalls,c)
		wg.Add(1)
	}
	wg.Wait()
}

func run (bodyChan chan *Body,waitGroup *sync.WaitGroup,count *Count, totalCalls int, c Client){
	defer waitGroup.Done()
	for {
		if count.add()>int64(totalCalls){
			break
		}
		startTime := time.Now()
		body := &Body{}
		RequestSize,ResponseSize,ok:=c.Call()
		body.Error=!ok
		body.RequestSize=RequestSize
		body.ResponseSize=ResponseSize
		body.Time = time.Now().Sub(startTime).Nanoseconds()/1000
		bodyChan <- body
	}
}
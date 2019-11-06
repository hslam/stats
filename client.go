package stats

import (
	"sync"
	"time"
)

type Client interface {
	Call()(RequestSize int64,ResponseSize int64,ok bool)
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
	startTime := time.Now()
	for {
		if count.add()>int64(totalCalls){
			break
		}
		startTime = time.Now()
		body := bodyPool.Get().(*Body)
		RequestSize,ResponseSize,ok:=c.Call()
		body.Error=!ok
		body.RequestSize=RequestSize
		body.ResponseSize=ResponseSize
		body.Time = time.Now().Sub(startTime).Nanoseconds()/1E3
		bodyChan <- body
	}
}

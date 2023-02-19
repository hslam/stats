// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package stats

import (
	"sync"
	"time"
)

// Client represents a client.
type Client interface {
	// Call returns a request size, a response size and its ok status.
	Call() (RequestSize int64, ResponseSize int64, Ok bool)
}

func startClient(bodyChans []chan *body, start <-chan struct{}, waitGroup *sync.WaitGroup, numParallels int, cnt *count, totalCalls int, c Client) {
	defer waitGroup.Done()
	wg := &sync.WaitGroup{}
	for i := 0; i < numParallels; i++ {
		go run(bodyChans[i], start, wg, cnt, totalCalls, c)
		wg.Add(1)
	}
	wg.Wait()
}

func run(bodyChan chan *body, start <-chan struct{}, waitGroup *sync.WaitGroup, cnt *count, totalCalls int, c Client) {
	defer waitGroup.Done()
	var startTime time.Time
	<-start
	for {
		if cnt.add(1) > int64(totalCalls) {
			break
		}
		startTime = time.Now()
		body := bodyPool.Get().(*body)
		RequestSize, ResponseSize, ok := c.Call()
		body.Error = !ok
		body.RequestSize = RequestSize
		body.ResponseSize = ResponseSize
		body.Time = time.Now().Sub(startTime).Nanoseconds()
		bodyChan <- body
	}
}

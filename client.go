// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package stats

import (
	"sync"
	"time"
)

//Client is the interface of client.
type Client interface {
	Call() (RequestSize int64, ResponseSize int64, Ok bool)
}

func startClient(bodyChan chan *Body, waitGroup *sync.WaitGroup, numParallels int, count *Count, totalCalls int, c Client) {
	defer waitGroup.Done()
	wg := &sync.WaitGroup{}
	for i := 0; i < numParallels; i++ {
		go run(bodyChan, wg, count, totalCalls, c)
		wg.Add(1)
	}
	wg.Wait()
}

func run(bodyChan chan *Body, waitGroup *sync.WaitGroup, count *Count, totalCalls int, c Client) {
	defer waitGroup.Done()
	var startTime time.Time
	for {
		if count.add(1) > int64(totalCalls) {
			break
		}
		startTime = time.Now()
		body := bodyPool.Get().(*Body)
		RequestSize, ResponseSize, ok := c.Call()
		body.Error = !ok
		body.RequestSize = RequestSize
		body.ResponseSize = ResponseSize
		body.Time = time.Now().Sub(startTime).Nanoseconds() / 1e3
		bodyChan <- body
	}
}

package main

import (
	"github.com/hslam/stats"
	"math/rand"
	"time"
)

func main() {
	var Clients []stats.Client
	for i := 0; i < 1e2; i++ {
		Clients = append(Clients, &Client{})
	}
	parallel := 32
	totalCalls := 1000000
	stats.StartPrint(parallel, totalCalls, Clients)
}

//Client implements the stats.Client interface.
type Client struct {
}

//Call returns RequestSize, ResponseSize, Ok.
func (c *Client) Call() (int64, int64, bool) {
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000))) //mock task time
	if rand.Intn(1000) == 1 {
		return 1e3, 0, false //mock error
	}
	return 1e3, 1e3, true //mock success
}

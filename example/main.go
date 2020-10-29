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

type Client struct {
}

func (c *Client) Call() (int64, int64, bool) {
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000)))
	if rand.Intn(1000) == 1 {
		return 1e3, 0, false
	}
	return 1e3, 1e3, true
}

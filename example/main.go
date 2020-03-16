package main

import (
	"github.com/hslam/stats"
	"math/rand"
	"time"
)

func main() {
	var Clients []stats.Client
	for i := 0; i < 1E2; i++ {
		Clients = append(Clients, &Client{})
	}
	parallel := 32
	total_calls := 1000000
	stats.StartPrint(parallel, total_calls, Clients)
}

//Client implements interface of client.
type Client struct {
}

//Call returns RequestSize, ResponseSize, Ok.
func (c *Client) Call() (int64, int64, bool) {
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000))) //to do time
	if rand.Intn(1000) == 1 {
		return 1E3, 0, false //error
	}
	return 1E3, 1E3, true //success
}

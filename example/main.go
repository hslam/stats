package main
import (
	"hslam.com/git/x/stats"
	"time"
	"math/rand"
)
func main()  {
	var Clients []stats.Client
	for i:=0;i<1E2 ;i++  {
		Clients=append(Clients, &Client{})
	}
	parallel:=32
	total_calls:=1000000
	stats.StartPrint(parallel,total_calls,Clients)
}
type Client struct {
}
func (c *Client)Call()(int64,int64,bool){
	time.Sleep(time.Microsecond*time.Duration(rand.Intn(1000)))
	return 1E3,1E3,true
}

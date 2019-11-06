# stats
Stats is written in Golang using for benchmark

## Get started

### Install
```
go get hslam.com/git/x/stats
```
### Import
```
import "hslam.com/git/x/stats"
```

### example
```
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
```

### Output
```
Summary:
	Clients:	100
	Parallels:	32
	Total Calls:	1000000
	Total time:	0.82s
	Requests per second:	1213991.49
	Fastest time for request:	0.00ms
	Average time per request:	2.54ms
	Slowest time for request:	29.81ms

Time:
	0.1%	time for request:	0.03ms
	1%	time for request:	0.40ms
	5%	time for request:	0.67ms
	10%	time for request:	0.86ms
	25%	time for request:	1.31ms
	50%	time for request:	2.04ms
	75%	time for request:	3.02ms
	90%	time for request:	4.78ms
	95%	time for request:	6.62ms
	99%	time for request:	10.18ms
	99.9%	time for request:	16.88ms

Request:
	Total request body sizes:	1000000000
	Average body size per request:	1000.00 Byte
	Request rate per second:	1213991494.78 Byte/s (1213.99 MByte/s)

Response:
	Total response body sizes:	1000000000
	Average body size per response:	1000.00 Byte
	Response rate per second:	1213991494.78 Byte/s (1213.99 MByte/s)

Result:
	ResponseOk:	1000000 (100.00%)
	Errors:	0 (0.00%)
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)


### Authors
stats was written by Mort Huang.

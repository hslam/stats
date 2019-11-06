# stats
Stats is written in golang using for benchmark

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
	time.Sleep(time.Microsecond*time.Duration(rand.Intn(1000)))	//to do time
	if rand.Intn(1000)==1{
		return 1E3,0,false		//error
	}
	return 1E3,1E3,true		//success
}
```

### Output
```
Summary:
	Clients:	100
	Parallels:	32
	Total Calls:	1000000
	Total time:	0.78s
	Requests per second:	1275476.04
	Fastest time for request:	0.00ms
	Average time per request:	2.41ms
	Slowest time for request:	30.80ms

Time:
	0.1%	time for request:	0.07ms
	1%	time for request:	0.51ms
	5%	time for request:	0.74ms
	10%	time for request:	0.92ms
	25%	time for request:	1.36ms
	50%	time for request:	2.07ms
	75%	time for request:	2.93ms
	90%	time for request:	4.02ms
	95%	time for request:	5.30ms
	99%	time for request:	9.30ms
	99.9%	time for request:	15.57ms

Request:
	Total request body sizes:	1000000000
	Average body size per request:	1000.00 Byte
	Request rate per second:	1275476039.54 Byte/s (1275.48 MByte/s)

Response:
	Total response body sizes:	998972000
	Average body size per response:	998.97 Byte
	Response rate per second:	1274164850.18 Byte/s (1274.16 MByte/s)

Result:
	ResponseOk:	998972 (99.90%)
	Errors:	1028 (0.10%)
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)


### Authors
stats was written by Mort Huang.

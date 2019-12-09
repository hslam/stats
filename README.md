# stats
Stats is written in golang using for benchmark

## Get started

### Install
```
go get github.com/hslam/stats
```
### Import
```
import "github.com/hslam/stats"
```

### example
```
package main
import (
	"github.com/hslam/stats"
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
	Parallel calls per client:	32
	Total calls:	1000000
	Total time:	0.70s
	Requests per second:	1427276.68
	Fastest time for request:	0.00ms
	Average time per request:	2.19ms
	Slowest time for request:	16.67ms

Time:
	0.1%	time for request:	0.00ms
	1%	time for request:	0.46ms
	5%	time for request:	0.71ms
	10%	time for request:	0.89ms
	25%	time for request:	1.32ms
	50%	time for request:	2.00ms
	75%	time for request:	2.78ms
	90%	time for request:	3.64ms
	95%	time for request:	4.33ms
	99%	time for request:	6.43ms
	99.9%	time for request:	9.79ms

Request:
	Total request body sizes:	1000000000
	Average body size per request:	1000.00 Byte
	Request rate per second:	1427276684.72 Byte/s (1427.28 MByte/s)

Response:
	Total response body sizes:	998982000
	Average body size per response:	1000.00 Byte
	Response rate per second:	1425823717.06 Byte/s (1425.82 MByte/s)

Result:
	Response ok:	998982 (99.90%)
	Errors:	1018 (0.10%)
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Meng Huang)


### Authors
stats was written by Meng Huang.

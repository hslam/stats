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
	Parallel calls per client:	32
	Total calls:	1000000
	Total time:	0.70s
	Requests per second:	1426356.50
	Fastest time for request:	0.00ms
	Average time per request:	2.18ms
	Slowest time for request:	17.98ms

Time:
	0.1%	time for request:	0.06ms
	1%	time for request:	0.40ms
	5%	time for request:	0.64ms
	10%	time for request:	0.80ms
	25%	time for request:	1.20ms
	50%	time for request:	1.80ms
	75%	time for request:	2.56ms
	90%	time for request:	3.62ms
	95%	time for request:	4.94ms
	99%	time for request:	10.39ms
	99.9%	time for request:	15.48ms

Request:
	Total request body sizes:	1000000000
	Average body size per request:	1000.00 Byte
	Request rate per second:	1426356500.69 Byte/s (1426.36 MByte/s)

Response:
	Total response body sizes:	998969000
	Average body size per response:	998.97 Byte
	Response rate per second:	1424885927.14 Byte/s (1424.89 MByte/s)

Result:
	Response ok:	998969 (99.90%)
	Errors:	1031 (0.10%)
```

### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)


### Authors
stats was written by Mort Huang.

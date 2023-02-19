# stats
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/stats)](https://pkg.go.dev/github.com/hslam/stats)
[![Build Status](https://github.com/hslam/stats/workflows/build/badge.svg)](https://github.com/hslam/stats/actions)
[![codecov](https://codecov.io/gh/hslam/stats/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/stats)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/stats?v=7e100)](https://goreportcard.com/report/github.com/hslam/stats)
[![GitHub release](https://img.shields.io/github/release/hslam/stats.svg)](https://github.com/hslam/stats/releases/latest)
[![LICENSE](https://img.shields.io/github/license/hslam/stats.svg?style=flat-square)](https://github.com/hslam/stats/blob/master/LICENSE)

Package stats implements a generic benchmarking tool.

## Get started

### Install
```
go get github.com/hslam/stats
```
### Import
```
import "github.com/hslam/stats"
```

### Example
```go
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
	parallel := 8
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
```

### Output
```
Summary:                                                                                                  
	Clients:	100
	Parallel calls per client:	8
	Total calls:	1000000
	Total time:	0.656s
	Requests per second:	1524557.683
	Fastest time for request:	0.000135ms
	Average time per request:	0.517865ms
	Slowest time for request:	3.276049ms

Time:
	00.0001%	time for request:	0.000135ms
	00.0010%	time for request:	0.000161ms
	00.0100%	time for request:	0.000265ms
	00.1000%	time for request:	0.002682ms
	01.0000%	time for request:	0.016050ms
	05.0000%	time for request:	0.061525ms
	10.0000%	time for request:	0.113671ms
	25.0000%	time for request:	0.265406ms
	50.0000%	time for request:	0.517284ms
	75.0000%	time for request:	0.768143ms
	90.0000%	time for request:	0.918319ms
	95.0000%	time for request:	0.968208ms
	99.0000%	time for request:	1.017011ms
	99.9000%	time for request:	1.381343ms
	99.9900%	time for request:	1.939697ms
	99.9990%	time for request:	2.706382ms
	99.9999%	time for request:	3.034106ms

Request:
	Total request body sizes:	1000000000
	Average body size per request:	1000.00 Byte
	Request rate per second:	1524557682.63 Byte/s (1524.56 MByte/s)

Response:
	Total response body sizes:	998969000
	Average body size per response:	1000.00 Byte
	Response rate per second:	1522985863.66 Byte/s (1522.99 MByte/s)

Result:
	Response ok:	998969 (99.897%)
	Errors:	1031 (0.103%)
```

### License
This package is licensed under a MIT license (Copyright (c) 2019 Meng Huang)


### Author
stats was written by Meng Huang.

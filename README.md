# stats
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/stats)](https://pkg.go.dev/github.com/hslam/stats)
[![Build Status](https://travis-ci.org/hslam/stats.svg?branch=master)](https://travis-ci.org/hslam/stats)
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

### License
This package is licensed under a MIT license (Copyright (c) 2019 Meng Huang)


### Author
stats was written by Meng Huang.

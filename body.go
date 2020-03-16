// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package stats

import (
	"sync"
)

var (
	bodyPool *sync.Pool
)

func init() {
	bodyPool = &sync.Pool{
		New: func() interface{} {
			return &Body{}
		},
	}
}

//Body defines the struct of response body.
type Body struct {
	RequestSize  int64
	ResponseSize int64
	Time         int64
	Error        bool
}

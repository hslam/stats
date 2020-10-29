// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package stats

import (
	"sync"
)

var (
	bodyPool = sync.Pool{
		New: func() interface{} {
			return &body{}
		},
	}
)

type body struct {
	RequestSize  int64
	ResponseSize int64
	Time         int64
	Error        bool
}

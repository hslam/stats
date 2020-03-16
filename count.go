// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package stats

import (
	"sync/atomic"
)

//Count is the count struct.
type Count struct {
	v int64
}

func (c *Count) add(delta int64) int64 {
	return atomic.AddInt64(&c.v, delta)
}

func (c *Count) load() int64 {
	return atomic.LoadInt64(&c.v)
}

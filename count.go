// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package stats

import (
	"sync/atomic"
)

type count struct {
	v int64
}

func (c *count) add(delta int64) int64 {
	return atomic.AddInt64(&c.v, delta)
}

func (c *count) load() int64 {
	return atomic.LoadInt64(&c.v)
}

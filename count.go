package stats

import (
	"sync/atomic"
)

type Count struct {
	v int64
}

func (c *Count)add()int64{
	return atomic.AddInt64(&c.v, 1);
}

func (c *Count)load()int64{
	return atomic.LoadInt64(&c.v);
}
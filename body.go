package stats

import (
	"sync"
)

var (
	bodyPool			*sync.Pool
)

func init() {
	bodyPool= &sync.Pool{
		New: func() interface{} {
			return &Body{}
		},
	}
}

type Body struct {
	RequestSize		int64
	ResponseSize	int64
	Time			int64
	Error			bool
}


// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package stats

import (
	"math/rand"
	"testing"
	"time"
)

func TestStats(t *testing.T) {
	var Clients []Client
	for i := 0; i < 1E2; i++ {
		Clients = append(Clients, &client{})
	}
	parallel := 32
	totalCalls := 1000000
	StartPrint(parallel, totalCalls, Clients)
}

type client struct {
}

func (c *client) Call() (int64, int64, bool) {
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000))) //to do time
	if rand.Intn(1000) == 1 {
		return 1E3, 0, false //error
	}
	return 1E3, 1E3, true //success
}

func TestBar(t *testing.T) {
	SetBar(false)
	getBar(0)
}

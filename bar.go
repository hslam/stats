// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package stats

func getBar(i int) (s string) {
	return getStr(i, "#") + getStr(1e2-i, " ")
}

func getStr(n int, char string) (s string) {
	if n < 1 {
		return
	}
	for i := 1; i <= n; i++ {
		s += char
	}
	return
}

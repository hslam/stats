package stats

func getBar(i int) (s string) {
	return getStr(i,"#") + getStr(1E2-i," ")
}

func getStr(n int,char string) (s string) {
	if n<1{
		return
	}
	for i:=1;i<=n;i++{
		s+=char
	}
	return
}

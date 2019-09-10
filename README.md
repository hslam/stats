```
type WrkClient struct {
}

func (c *WrkClient)Call()(int64,int64,bool){
    //To Do
    //return 1024,1024,true
	return 1024,0,false
}
```

```
	var wrkClients []stats.Client
	parallel:=1
	total_calls:=1000000
	wrkClients[0]= &WrkClient{}
	stats.StartClientStats(parallel,total_calls,wrkClients)
```


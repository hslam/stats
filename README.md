# stats
## benchmark test


## Get started

### Install
```
go get hslam.com/git/x/stats
```
### Import
```
import "hslam.com/git/x/stats"
```

### example
```
type Client struct {
}

func (c *Client)Call()(int64,int64,bool){
    //To Do
	//return 1024,0,false
    return 1024,1024,true
}

func example(){
    var Clients []stats.Client
	parallel:=1
	total_calls:=1000000
	Clients[0]= &Client{}
	stats.StartPrint(parallel,total_calls,Clients)
}
```


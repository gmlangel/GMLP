package main;

import(
	"fmt"
	"./proxy"
	"time"
	"./luboService"
	//"strings"
)

func main(){
	runloopChan := make(chan int);
	fmt.Println("I'm Socket's main, time:",uint32(time.Now().Unix())," version = ","1.0.1");

	//开启录播服务
	socketPro := proxy.NewLoBoSocket();
	socketPro.GInit();

	//开启定时删除过去数据的任务
	go luboService.ExecDeadLineTask();
	
	runloopChan  <- 1;//启动runloop
}
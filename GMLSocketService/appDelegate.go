package main;

import(
	"./proxy"
	//"./luboService"
)

func main(){
	runloopChan := make(chan int);

	// //开启录播服务
	// socketPro := proxy.NewLoBoSocket();
	// socketPro.GInit();

	// //开启定时删除过去数据的任务
	// go luboService.ExecDeadLineTask();

	//开启RemoteDecition服务
	rdSocketPro := proxy.NewRemoteDecitionSocket();
	rdSocketPro.GInit();
	runloopChan  <- 1;//启动runloop
}
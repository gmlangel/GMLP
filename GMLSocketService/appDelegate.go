package main;

import(
	"fmt"
	"./proxy"
	"time"
)

func main(){
	runloopChan := make(chan int);
	fmt.Println("I'm Socket's main, time:",uint32(time.Now().Unix()));
	//测试用
	socketPro := proxy.NewLoBoSocket();
	socketPro.GInit();
	// //测试关闭重启
	// time.Sleep(time.Second * time.Duration(70));
	// socketPro.DeInit();

	// time.Sleep(time.Second * time.Duration(3));
	// socketPro.GInit();
	
	runloopChan  <- 1;//启动runloop
}
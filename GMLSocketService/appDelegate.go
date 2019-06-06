package main;

import(
	"fmt"
	"./proxy"
	"time"
	//"strings"
)

func main(){
	runloopChan := make(chan int);
	fmt.Println("I'm Socket's main, time:",uint32(time.Now().Unix())," version = ",1.0);

	//开启录播服务
	socketPro := proxy.NewLoBoSocket();
	socketPro.GInit();
	
	runloopChan  <- 1;//启动runloop
}
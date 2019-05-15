package main;

import(
	"fmt"
	"./proxy"
)

func main(){
	runloopChan := make(chan int);
	fmt.Println("I'm Socket's main");
	//测试用
	socketPro := proxy.NewLoBoSocket();
	socketPro.Start();

	//测试关闭
	// time.Sleep(time.Second * 30);
	// socketPro.Stop();
	//socketPro.DeInit();
	runloopChan  <- 1;//启动runloop
}
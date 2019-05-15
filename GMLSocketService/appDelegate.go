package main;

import(
	"fmt"
	"./proxy"
)

func main(){
	fmt.Println("I'm Socket's main");
	//测试用
	socketPro := proxy.NewLoBoSocket();
	socketPro.Start();
	//socketPro.DeInit();
}
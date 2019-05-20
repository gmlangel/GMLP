package main;

import(
	"fmt"
	"./proxy"
	"time"
	//"strings"
)

func main(){
	runloopChan := make(chan int);
	fmt.Println("I'm Socket's main, time:",uint32(time.Now().Unix()));

	// str := "郭明龙<gmle>中风abc你是";
	// tbyte:= []byte{131,173,230 ,152, 142, 233, 190, 153, 60, 103, 109, 108, 101, 62, 228, 184, 173, 233, 163, 142, 97, 98, 99, 228, 189, 160, 230, 152, 175};
	// fmt.Println([]byte(str));
	
	// str2 := string(tbyte)
	// fmt.Println(str2);

	// j := strings.Index(str2,"<gmle>");
	// str3 := str2[0:j];
	// fmt.Println(str3);

	// delbyte:= []byte(str3);
	// j = len(delbyte);

	// fmt.Println(string(tbyte[j:]));

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
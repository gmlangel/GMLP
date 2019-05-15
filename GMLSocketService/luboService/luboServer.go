package luboService;

import(
	"fmt"
	model "../models"
	"log"
	"net"
	"os"
)

type LuboServer struct{
	isStarted bool;
	sockListener *net.Listener;/*socket监听器*/
}

/**
初始化录播服务器
*/
func (sev *LuboServer)Init(conf *model.LoBoServerConfig){
	if conf.Host != "" && conf.Port != 0 && conf.ServerType != ""{
		hostPath := fmt.Sprintf("%s:%v",conf.Host,conf.Port);
		tempListen,err := net.Listen(conf.ServerType,hostPath);
		if err != nil{
			log.Println(fmt.Fprintf(os.Stderr,err.Error()));
			return;
		}
		if tempListen == nil{
			return;
		}
		sev.sockListener = &tempListen;//获取socket服务的引用
		fmt.Println("录播服务初始化成功");
	}
}

/**
释放录播服务器
*/
func (sev *LuboServer)DeInit(){

}

/**
开始录播服务
*/
func (sev *LuboServer)OpenServer(){
	if sev.isStarted == true{
		return;
	}
	fmt.Println("服务器启动成功");
	sev.isStarted = true;
	go sev._openServer();
}

func (sev *LuboServer)_openServer(){
	for sev.isStarted == true{
		//接受客户端发来的建立socket的请求
		newClient,err := (*sev.sockListener).Accept();
		if err != nil{
			log.Println(fmt.Fprintf(os.Stderr,err.Error()));
			continue;
		}
		//打印客户端信息
		log.Println(newClient.RemoteAddr().String(), " tcp connect success");
	}
}

/**
停止录播服务
*/
func (sev *LuboServer)CloseServer(){
	if sev.isStarted == false{
		return;
	}
	fmt.Println("服务器停止");
	sev.isStarted = false;
}
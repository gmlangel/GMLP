package proxy;

import(
	server "../remoteDecitionService"
	model "../models"
	"fmt"
	"time"
)

/**
录播服务业务委托
*/
type RemoteDecitionSocketProxy struct{
	serv *server.RemoteDecitionServer;
	isInited bool
}

func NewRemoteDecitionSocket()(* RemoteDecitionSocketProxy){
	ins := &RemoteDecitionSocketProxy{};
	ins.isInited = false;
	return ins;
}

/**
启动服务
*/
func (pro *RemoteDecitionSocketProxy)GInit(){
	fmt.Println("I'm RemoteDecitionSocket, time:",uint32(time.Now().Unix())," version = ","1.0.0");
	if pro.isInited == true{
		return;
	}
	pro.isInited = true;
	conf := &model.LoBoServerConfig{Host:"0.0.0.0",Port:63333,ServerType:"tcp"};
	pro.serv = &server.RemoteDecitionServer{};//创建录播服务
	pro.serv.Init(conf);
	pro.serv.OpenServer();
}



/**
释放操作
*/
func (pro *RemoteDecitionSocketProxy)DeInit(){
	if pro.isInited == false{
		return;
	}
	pro.isInited = false;
	if pro.serv != nil{
		pro.serv.DeInit();
		pro.serv = nil;
	}
}
package proxy;

import(
	server "../luboService"
	model "../models"
)

/**
录播服务业务委托
*/
type LoboSocketProxy struct{
	serv *server.LuboServer;
	isInited bool
}

func NewLoBoSocket()(* LoboSocketProxy){
	ins := &LoboSocketProxy{};
	ins.isInited = false;
	return ins;
}

/**
启动服务
*/
func (pro *LoboSocketProxy)GInit(){
	if pro.isInited == true{
		return;
	}
	pro.isInited = true;
	conf := &model.LoBoServerConfig{Host:"0.0.0.0",Port:59999,ServerType:"tcp"};
	pro.serv = &server.LuboServer{};//创建录播服务
	pro.serv.Init(conf);
	pro.serv.OpenServer();
}



/**
释放操作
*/
func (pro *LoboSocketProxy)DeInit(){
	if pro.isInited == false{
		return;
	}
	pro.isInited = false;
	if pro.serv != nil{
		pro.serv.DeInit();
		pro.serv = nil;
	}
}
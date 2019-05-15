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

}

func NewLoBoSocket()(* LoboSocketProxy){
	ins := &LoboSocketProxy{};
	ins.serv = &server.LuboServer{};//创建录播服务

	conf := &model.LoBoServerConfig{Host:"0.0.0.0",Port:59999,ServerType:"tcp"};
	ins.serv.Init(conf);
	
	return ins;
}

/**
开始录播服务
*/
func (pro *LoboSocketProxy)Start(){
	if pro.serv != nil{
		pro.serv.OpenServer();
	}
}

/**
停止录播服务
*/
func (pro *LoboSocketProxy)Stop(){
	if pro.serv != nil{
		pro.serv.CloseServer();
	}
}


/**
释放操作
*/
func (pro *LoboSocketProxy)DeInit(){
	pro.Stop();
	if pro.serv != nil{
		pro.serv.DeInit();
		pro.serv = nil;
	}
}
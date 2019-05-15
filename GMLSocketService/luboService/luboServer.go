package luboService;

import(
	"fmt"
	model "../models"
	"log"
	"net"
	"os"
)

type LuboServerError struct{
	msg string
}

func (err *LuboServerError)Error()string{
	return err.msg;
}

type LuboServer struct{
	isStarted bool;
	sockListener net.Listener;/*socket监听器*/
	connectIdOffset int64;/*客户端连接的ID ,用于生成 '客户端连接ID池'*/
	connectIdPool []int64;/*客户端连接ID池*/
	unOwnedConnect map[int64]*LuBoClientConnection;/*无主连接字典.用于记录未进入教室的用户的socket链接 {sid:socket}*/
	ownedConnect map[int64]*LuBoClientConnection;/*有主连接字典{sid:socket}*/
	ownedConnectUIDMap map[int64]*LuBoClientConnection;/*有主连接字典{uid:socket}*/
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
		sev.sockListener = tempListen;//获取socket服务的引用
		sev.unOwnedConnect = map[int64]*LuBoClientConnection{};
		sev.ownedConnect = map[int64]*LuBoClientConnection{};
		sev.ownedConnectUIDMap = map[int64]*LuBoClientConnection{};

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
		newClient,err := sev.sockListener.Accept();
		if err != nil{
			log.Println(fmt.Fprintf(os.Stderr,err.Error()));
			continue;
		}
		//生成新的socket的id
		sid,err:= sev.createConnectId();
		if err != nil{
			log.Println("生成socketID出错，当前socket被丢弃.错误原因",err.Error());
			continue;
		}
		//生成socket的管理器
		luboclient := NewLuBoClientConn(sid,newClient);
		luboclient.OnTimeout = func(cli * LuBoClientConnection){
			sev.destroySocket(cli);//释放socket
		}
		luboclient.OnError = func(cli * LuBoClientConnection){
			sev.destroySocket(cli);//释放socket
		}
		//塞入无主socket记录集
		sev.unOwnedConnect[sid] = luboclient;
		//打印客户端信息
		log.Println("new Client join,Address:",newClient.RemoteAddr().String(), " discription:",luboclient);
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
	sev.sockListener.Close();
}

/*生成socketID*/
func (sev *LuboServer)createConnectId()(int64,error){
	if len(sev.connectIdPool) == 0{
		if sev.connectIdOffset < 0xfffffffffffffe - 10000{
			for i := 1;i<=10000;i++{
				sev.connectIdPool = append(sev.connectIdPool,sev.connectIdOffset + int64(i));
            }
            sev.connectIdOffset += 10000;
		}else{
			return 0,&LuboServerError{msg:"无法继续生成connectId,因为ID超出最大限制"};
		}
	}
	j := len(sev.connectIdPool);
	currentArr := sev.connectIdPool[:j-1];
	popArr := sev.connectIdPool[j-1:];
	result := popArr[0];
	sev.connectIdPool = currentArr;
	return result,nil;
}

/*释放一个client Socket*/
func (sev *LuboServer)destroySocket(cli * LuBoClientConnection){
	sid := cli.SID;
	uid := cli.UID;
	//rid := cli.RID;
	if sev.unOwnedConnect[sid] != nil{
		sev.unOwnedConnect[sid] = nil;
	}

	if sev.ownedConnect[sid] != nil{
		sev.ownedConnect[sid] = nil;
	}

	if sev.ownedConnectUIDMap[uid] != nil{
		sev.ownedConnectUIDMap[uid] = nil;
	}
	//返回socketID 到id池,以便之后的链接使用
	sev.connectIdPool = append(sev.connectIdPool,sid);
	cli.DestroySocket("你被删除了");
	
}
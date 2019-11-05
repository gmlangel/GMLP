package remoteDecitionService;

import(
	"fmt"
	model "../models"
	"log"
	"net"
	"errors"
	"os"
)

/**全局记录集合声明及操作函数*/
var(
	DeadLineInterval int64 = 24 * 3600;//数据的过期时间，默认为缓存1天
	unOwnedConnectChan chan int = make(chan int,1);
	unOwnedConnect map[int64]*RemoteDecitionConnection = map[int64]*RemoteDecitionConnection{};/*无主连接字典.用于记录未进入教室的用户的socket链接 {sid:socket}*/
	
	ownedConnectChan chan int = make(chan int,1);
	ownedConnect map[int64]*RemoteDecitionConnection = map[int64]*RemoteDecitionConnection{};/*有主连接字典{sid:socket}*/
	
	ownedConnectUIDMapChan chan int = make(chan int,1);
	ownedConnectUIDMap map[int64]*RemoteDecitionConnection = map[int64]*RemoteDecitionConnection{};/*有主连接字典{uid:socket}*/


	connectIdPool = []int64{};/*客户端连接ID池*/
	destroyChan = make(chan int,1);/*client socket释放操作时的互斥锁*/
)

//向unOwnedConnect 添加或者设置值
func UnOwnedConnect_SetValue(key int64,value *RemoteDecitionConnection){
	unOwnedConnectChan <- 1;
	unOwnedConnect[key] = value;
	<-unOwnedConnectChan;
}
/*获取unOwnedConnect中的指定key对应的值*/
func UnOwnedConnect_GetValue(key int64)*RemoteDecitionConnection{
	var result *RemoteDecitionConnection = nil;
	unOwnedConnectChan <- 1;
	result = unOwnedConnect[key];
	<-unOwnedConnectChan;
	return result;
}
/*清空unOwnedConnect*/
func UnOwnedConnect_Clear(){
	unOwnedConnectChan <- 1;
	unOwnedConnect = map[int64]*RemoteDecitionConnection{};
	<-unOwnedConnectChan;
}


/*向ownedConnect 添加或者设置值*/
func OwnedConnect_SetValue(key int64,value *RemoteDecitionConnection){
	ownedConnectChan <- 1;
	ownedConnect[key] = value;
	<-ownedConnectChan;
}
/*获取OwnedConnect中的指定key对应的值*/
func OwnedConnect_GetValue(key int64)*RemoteDecitionConnection{
	var result *RemoteDecitionConnection = nil;
	ownedConnectChan <- 1;
	result = ownedConnect[key];
	<-ownedConnectChan;
	return result;
}
/*清空OwnedConnect*/
func OwnedConnect_Clear(){
	ownedConnectChan <- 1;
	ownedConnect = map[int64]*RemoteDecitionConnection{};
	<-ownedConnectChan;
}


/*向ownedConnect 添加或者设置值*/
func OwnedConnectUIDMap_SetValue(key int64,value *RemoteDecitionConnection){
	ownedConnectUIDMapChan <- 1;
	ownedConnectUIDMap[key] = value;
	<-ownedConnectUIDMapChan;
}
/*获取ownedConnectUIDMap中的指定key对应的值*/
func OwnedConnectUIDMap_GetValue(key int64)*RemoteDecitionConnection{
	var result *RemoteDecitionConnection = nil;
	ownedConnectUIDMapChan <- 1;
	result = ownedConnectUIDMap[key];
	<-ownedConnectUIDMapChan;
	return result;
}
/*清空ownedConnectUIDMap*/
func OwnedConnectUIDMap_Clear(){
	ownedConnectUIDMapChan <- 1;
	ownedConnectUIDMap = map[int64]*RemoteDecitionConnection{};
	<-ownedConnectUIDMapChan;
}

type RemoteDecitionServer struct{
	isStarted bool;
	sockListener net.Listener;/*socket监听器*/
	connectIdOffset int64;/*客户端连接的ID ,用于生成 '客户端连接ID池'*/
}

/**
初始化录播服务器
*/
func (sev *RemoteDecitionServer)Init(conf *model.LoBoServerConfig){
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

		fmt.Println("RemoteDicition服务初始化成功");
	}
}

/**
释放录播服务器
*/
func (sev *RemoteDecitionServer)DeInit(){
	//关闭服务监听
	sev.CloseServer();
}

/**
开始录播服务
*/
func (sev *RemoteDecitionServer)OpenServer(){
	if sev.isStarted == true{
		return;
	}
	fmt.Println("服务器启动成功");
	sev.isStarted = true;
	go sev._openServer();
}

func (sev *RemoteDecitionServer)_openServer(){
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
		remoteDecitionclient := NewRemoteDecitionConn(sid,newClient);
		remoteDecitionclient.OnTimeout = func(cli * RemoteDecitionConnection){
			DestroySocket(cli,nil);//释放socket
		}
		remoteDecitionclient.OnError = func(cli * RemoteDecitionConnection){
			DestroySocket(cli,nil);//释放socket
		}
		//塞入无主socket记录集
		UnOwnedConnect_SetValue(sid,remoteDecitionclient);
		//打印客户端信息
		fmt.Println("new Client join,Address:",newClient.RemoteAddr().String(), " discription:",remoteDecitionclient);
	}
}


/*释放一个client Socket*/
func DestroySocket(cli * RemoteDecitionConnection,completeFunc func()){
	destroyChan <- 1
	sid := cli.SID;
	uid := cli.UID;
	//rid := cli.RID;
	UnOwnedConnect_SetValue(sid,nil);
	OwnedConnect_SetValue(sid,nil);
	OwnedConnectUIDMap_SetValue(uid,nil);
	//返回socketID 到id池,以便之后的链接使用
	connectIdPool = append(connectIdPool,sid);
	res := CreateProtocal(model.S_NOTIFY_C_OFFLINE).(*model.OfflineNotify_s2c);
	res.Code = 259;
	res.Reason = "您已经被踢";
	cli.OnSocketCloseComplete = completeFunc;//设置 cli完成关闭后的处理函数
	cli.DestroySocket(res);
	//fmt.Println(fmt.Sprintf("sev.unOwnedConnect = %v, sev.ownedConnect = %v, sev.ownedConnectUIDMap = %v",unOwnedConnect,ownedConnect,ownedConnectUIDMap))
	<- destroyChan
}

/**
停止录播服务
*/
func (sev *RemoteDecitionServer)CloseServer(){
	if sev.isStarted == false{
		return;
	}
	fmt.Println("服务器停止");
	sev.isStarted = false;
	sev.sockListener.Close();
	sev.sockListener = nil;

	//停止所有socket
	for key := range unOwnedConnect{
		if sock := unOwnedConnect[key];sock != nil{
			DestroySocket(sock,nil);
		}
	}

	for key := range ownedConnect{
		if sock := ownedConnect[key];sock != nil{
			DestroySocket(sock,nil);
		}
	}

	for key := range ownedConnectUIDMap{
		if sock := ownedConnectUIDMap[key];sock != nil{
			DestroySocket(sock,nil);
		}
	}
	//释放数组和集合
	UnOwnedConnect_Clear();
	OwnedConnect_Clear();
	OwnedConnectUIDMap_Clear();
}

/*生成socketID*/
func (sev *RemoteDecitionServer)createConnectId()(int64,error){
	if len(connectIdPool) == 0{
		if sev.connectIdOffset < 0xfffffffffffffe - 10000{
			for i := 1;i<=10000;i++{
				connectIdPool = append(connectIdPool,sev.connectIdOffset + int64(i));
            }
            sev.connectIdOffset += 10000;
		}else{
			return 0,errors.New("无法继续生成connectId,因为ID超出最大限制");
		}
	}
	j := len(connectIdPool);
	currentArr := connectIdPool[:j-1];
	popArr := connectIdPool[j-1:];
	result := popArr[0];
	connectIdPool = currentArr;
	return result,nil;
}

/*有一个新用户登录*/
func NewUserClientlogin(sid int64,uid int64,client *RemoteDecitionConnection){
	UnOwnedConnect_SetValue(sid,nil);
	OwnedConnect_SetValue(sid,client);
	OwnedConnectUIDMap_SetValue(uid,client);
}
package luboService;

import(
	"fmt"
	model "../models"
	"log"
	"net"
	"os"
)
/**全局记录集合声明及操作函数*/
var(
	unOwnedConnectChan chan int = make(chan int,1);
	unOwnedConnect map[int64]*LuBoClientConnection = map[int64]*LuBoClientConnection{};/*无主连接字典.用于记录未进入教室的用户的socket链接 {sid:socket}*/
	
	ownedConnectChan chan int = make(chan int,1);
	ownedConnect map[int64]*LuBoClientConnection = map[int64]*LuBoClientConnection{};/*有主连接字典{sid:socket}*/
	
	ownedConnectUIDMapChan chan int = make(chan int,1);
	ownedConnectUIDMap map[int64]*LuBoClientConnection = map[int64]*LuBoClientConnection{};/*有主连接字典{uid:socket}*/

	roomInfoMap map[int64]*model.RoomInfo = map[int64]*model.RoomInfo{};
	roomInfoMapChan = make(chan int,1);

	lessonResultMap map[string][]model.Answer_c2s = map[string][]model.Answer_c2s{};//课程报告数据集
	lessonResultMapChan = make(chan int,1);

	teachScriptMap map[int64]map[string]interface{} = map[int64]map[string]interface{}{};//教学脚本数据集
	teachScriptMapChan = make(chan int,1);
)

/*根据SID和UID获取 socket连接*/
func GetSocketByUIDAndSID(sid int64,uid int64)*LuBoClientConnection{
	var sock *LuBoClientConnection = UnOwnedConnect_GetValue(sid);
	if sock == nil{
		sock = OwnedConnect_GetValue(sid);
		if sock == nil{
			sock = OwnedConnectUIDMap_GetValue(uid);
		}
	}
	return sock;
}

/*有一个新用户进入了教室，需要同步各sock缓存集合的记录*/
func NewUserClientJoinRoom(sid int64,uid int64,client *LuBoClientConnection){
	UnOwnedConnect_SetValue(sid,nil);
	OwnedConnect_SetValue(sid,client);
	OwnedConnectUIDMap_SetValue(uid,client);
}


/*向teachScriptMap 添加或者设置值*/
func TeachScriptMap_SetValue(key int64,value map[string]interface{}){
	teachScriptMapChan <- 1;
	teachScriptMap[key] = value;
	<-teachScriptMapChan;
}
/*获取teachScriptMap中的指定key对应的值*/
func TeachScriptMap_GetValue(key int64)map[string]interface{}{
	var result map[string]interface{} = nil;
	teachScriptMapChan <- 1;
	result =teachScriptMap[key];
	<-teachScriptMapChan;
	return result;
}
/*清空teachScriptMap*/
func TeachScriptMap_Clear(){
	teachScriptMapChan <- 1;
	teachScriptMap = map[int64]map[string]interface{}{};
	<-teachScriptMapChan;
}

/*向lessonResultMap 添加或者设置值*/
func LessonResultMap_SetValue(key string,value []model.Answer_c2s){
	lessonResultMapChan <- 1;
	lessonResultMap[key] = value;
	<-lessonResultMapChan;
}
/*获取lessonResultMap中的指定key对应的值*/
func LessonResultMap_GetValue(key string)[]model.Answer_c2s{
	var result []model.Answer_c2s = nil;
	lessonResultMapChan <- 1;
	result = lessonResultMap[key];
	<-lessonResultMapChan;
	return result;
}
/*清空lessonResultMap*/
func LessonResultMap_Clear(){
	lessonResultMapChan <- 1;
	lessonResultMap = map[string][]model.Answer_c2s{};
	<-lessonResultMapChan;
}

/*向roomInfoMap 添加或者设置值*/
func RoomInfoMap_SetValue(key int64,value *model.RoomInfo){
	roomInfoMapChan <- 1;
	roomInfoMap[key] = value;
	<-roomInfoMapChan;
}
/*获取roomInfoMap中的指定key对应的值*/
func RoomInfoMap_GetValue(key int64)*model.RoomInfo{
	var result *model.RoomInfo = nil;
	roomInfoMapChan <- 1;
	result =roomInfoMap[key];
	<-roomInfoMapChan;
	return result;
}
/*清空roomInfoMap*/
func RoomInfoMap_Clear(){
	roomInfoMapChan <- 1;
	roomInfoMap = map[int64]*model.RoomInfo{};
	<-roomInfoMapChan;
}



/*向unOwnedConnect 添加或者设置值*/
func UnOwnedConnect_SetValue(key int64,value *LuBoClientConnection){
	unOwnedConnectChan <- 1;
	unOwnedConnect[key] = value;
	<-unOwnedConnectChan;
}
/*获取unOwnedConnect中的指定key对应的值*/
func UnOwnedConnect_GetValue(key int64)*LuBoClientConnection{
	var result *LuBoClientConnection = nil;
	unOwnedConnectChan <- 1;
	result = unOwnedConnect[key];
	<-unOwnedConnectChan;
	return result;
}
/*清空unOwnedConnect*/
func UnOwnedConnect_Clear(){
	unOwnedConnectChan <- 1;
	unOwnedConnect = map[int64]*LuBoClientConnection{};
	<-unOwnedConnectChan;
}


/*向ownedConnect 添加或者设置值*/
func OwnedConnect_SetValue(key int64,value *LuBoClientConnection){
	ownedConnectChan <- 1;
	ownedConnect[key] = value;
	<-ownedConnectChan;
}
/*获取OwnedConnect中的指定key对应的值*/
func OwnedConnect_GetValue(key int64)*LuBoClientConnection{
	var result *LuBoClientConnection = nil;
	ownedConnectChan <- 1;
	result = ownedConnect[key];
	<-ownedConnectChan;
	return result;
}
/*清空OwnedConnect*/
func OwnedConnect_Clear(){
	ownedConnectChan <- 1;
	ownedConnect = map[int64]*LuBoClientConnection{};
	<-ownedConnectChan;
}


/*向ownedConnect 添加或者设置值*/
func OwnedConnectUIDMap_SetValue(key int64,value *LuBoClientConnection){
	ownedConnectUIDMapChan <- 1;
	ownedConnectUIDMap[key] = value;
	<-ownedConnectUIDMapChan;
}
/*获取ownedConnectUIDMap中的指定key对应的值*/
func OwnedConnectUIDMap_GetValue(key int64)*LuBoClientConnection{
	var result *LuBoClientConnection = nil;
	ownedConnectUIDMapChan <- 1;
	result = ownedConnectUIDMap[key];
	<-ownedConnectUIDMapChan;
	return result;
}
/*清空ownedConnectUIDMap*/
func OwnedConnectUIDMap_Clear(){
	ownedConnectUIDMapChan <- 1;
	ownedConnectUIDMap = map[int64]*LuBoClientConnection{};
	<-ownedConnectUIDMapChan;
}

/**服务错误声明*/
type LuboServerError struct{
	msg string
}

func (err *LuboServerError)Error()string{
	return err.msg;
}


/**服务业务处理类*/
type LuboServer struct{
	isStarted bool;
	sockListener net.Listener;/*socket监听器*/
	connectIdOffset int64;/*客户端连接的ID ,用于生成 '客户端连接ID池'*/
	connectIdPool []int64;/*客户端连接ID池*/
	
	destroyChan chan int;/*client socket释放操作时的互斥锁*/
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
		
		sev.destroyChan = make(chan int,1);//初始化 client socket释放操作时的互斥锁

		fmt.Println("录播服务初始化成功");
	}
}

/**
释放录播服务器
*/
func (sev *LuboServer)DeInit(){
	//关闭服务监听
	sev.CloseServer();
	//释放互斥锁
	close(sev.destroyChan);
	sev.destroyChan = nil;
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
		UnOwnedConnect_SetValue(sid,luboclient);
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
	sev.sockListener = nil;

	//停止所有socket
	for key := range unOwnedConnect{
		if sock := unOwnedConnect[key];sock != nil{
			sev.destroySocket(sock);
		}
	}

	for key := range ownedConnect{
		if sock := ownedConnect[key];sock != nil{
			sev.destroySocket(sock);
		}
	}

	for key := range ownedConnectUIDMap{
		if sock := ownedConnectUIDMap[key];sock != nil{
			sev.destroySocket(sock);
		}
	}
	//释放数组和集合
	UnOwnedConnect_Clear();
	OwnedConnect_Clear();
	OwnedConnectUIDMap_Clear();
	RoomInfoMap_Clear();
	LessonResultMap_Clear();
	TeachScriptMap_Clear();
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
	sev.destroyChan <- 1
	sid := cli.SID;
	uid := cli.UID;
	//rid := cli.RID;
	UnOwnedConnect_SetValue(sid,nil);
	OwnedConnect_SetValue(sid,nil);
	OwnedConnectUIDMap_SetValue(uid,nil);
	//返回socketID 到id池,以便之后的链接使用
	sev.connectIdPool = append(sev.connectIdPool,sid);
	res := CreateProtocal(model.S_NOTIFY_C_OFFLINE).(*model.OfflineNotify_s2c);
	res.Code = 259;
	res.Reason = "您已经被踢";
	cli.DestroySocket(res);
	fmt.Println(fmt.Sprintf("sev.unOwnedConnect = %v, sev.ownedConnect = %v, sev.ownedConnectUIDMap = %v",unOwnedConnect,ownedConnect,ownedConnectUIDMap))
	<- sev.destroyChan
}
package remoteDecitionService;

import(
	"net"
	"time"
	"log"
	model "../models"
	"encoding/json"
	"fmt"
	"strings"
)

const(
	 socketTimeoutSecond time.Duration= time.Second * 60;
	 pkgHead = "<gmlb>";//包头
	 pkgFoot = "<gmle>";//包尾
)

/**
用于管理单个客户端的Socket连接
*/
type RemoteDecitionConnection struct{
	isConnected bool/*是否处于连接状态*/
	runLoopExecChan chan int;/*脚本循环处理锁*/
	dataSyncChan chan int;/*课中信息读取互斥锁*/
	writeChan chan int;/*socket写入操作，锁*/
	waitSendMSGBuffer []interface{};/*即将发送给客户端的消息的队列*/
	readChan chan int;/*数据读取锁*/
	readBuffer []byte;/*数据读取队列*/
	SID int64 /*socket id*/
	UID int64 /*用户ID默认为-1*/
	Sock net.Conn;/*真正的socket操作对象*/
	TimeoutSecond time.Duration;/*超时时长*/
	OnTimeout func (*RemoteDecitionConnection);/*当当前socket超时时触发*/
	OnError func (*RemoteDecitionConnection);/*当当前socket 发生除了Timeout错误以外，如EOF时触发*/
	OnSocketCloseComplete func();//当当前socket连接被close完毕后触发的处理函数。
	CurrentStrategyConfigPath string;
	CurrentConditionConfigPath string;
}


/**
创建新的socket管理工具
*/
func NewRemoteDecitionConn(sid int64,conn net.Conn)(*RemoteDecitionConnection){
	client := &RemoteDecitionConnection{SID:sid,UID:-1,Sock:conn,TimeoutSecond:socketTimeoutSecond,isConnected:true};
	client.writeChan = make(chan int,1);//初始化 socket写入操作锁
	client.readChan = make(chan int,1);
	client.runLoopExecChan = make(chan int,1);
	client.dataSyncChan = make(chan int,1);
	client.writeChan <- 1;
	client.readChan <- 1;
	client.runLoopExecChan <- 1;
	client.dataSyncChan <-1;
	client.UID = -1;
	go client.runLoopRead();//开socket read队列
	//封装并回执给客户端“客户端接入”信令
	pro := CreateProtocal(model.S_RES_C_HY);
	client.Write(pro);
	return client;
}

/**
关闭并释放socket,支持在关闭socket之前，依然向原有socket中写入一条消息
*/
func (lbc *RemoteDecitionConnection)DestroySocket(arg interface{}){
	lbc.UID = -1;
	if lbc.isConnected == false{
		return;
	}
	lbc.isConnected = false;
	lbc.OnTimeout = nil;
	lbc.OnError = nil;
	go lbc.writeLastMsgAndCloseSock(arg);//送最后一条消息后，关闭socket
	
}

/*发送最后一条数据给客户端，之后关闭socket*/
func (lbc *RemoteDecitionConnection)writeLastMsgAndCloseSock(arg interface{}){
	sock := lbc.Sock
	if sock == nil{
		return;
	}
	_,isOk := <- lbc.writeChan;
	if false == isOk{
		return;
	}
	//执行写入
	sock.SetDeadline(time.Now().Add(socketTimeoutSecond));//延长超时时间
	data,err := json.Marshal(arg);
	tj := len(data);
	if err != nil || tj <= 2{
		log.Println("sock:",lbc.SID,"数据转换出错:",err.Error())
		lbc.writeChan <- 1;
		return;
	}
	data = append([]byte(pkgHead),data...);
	data = append(data,[]byte(pkgFoot)...);
	fmt.Println("sock:",lbc.SID," 准备写入socket的数据:",string(data));
	n,err := sock.Write(data);
	_ = n;
	if err != nil{
		lbc.writeChan <- 1;
		lbc.closeSocket();
		return;
	}
	lbc.writeChan <- 1;
	lbc.closeSocket();
}

/*
发送数据给客户端
*/
func (lbc *RemoteDecitionConnection)Write(arg interface{}){
	_,isOk := <- lbc.writeChan;
	if false == isOk{
		return;
	}
	lbc.waitSendMSGBuffer = append(lbc.waitSendMSGBuffer,arg);//向 “等待发送的消息队列中”添加一条消息
	lbc.writeChan <- 1;
	go lbc.writeToSocket();//调用消息发送函数
}

/*检查是否需要关闭客户端socket连接*/
func (lbc *RemoteDecitionConnection)closeSocket(){
	lbc.SID = -1;
	lbc.Sock.Close();
	lbc.Sock = nil;
	<- lbc.writeChan
	close(lbc.writeChan)//关闭channel
	<- lbc.readChan
	close(lbc.readChan)
	<- lbc.runLoopExecChan
	close(lbc.runLoopExecChan);
	<- lbc.dataSyncChan
	close(lbc.dataSyncChan);
	if lbc.OnSocketCloseComplete != nil{
		lbc.OnSocketCloseComplete();
		lbc.OnSocketCloseComplete = nil;
	}
}

func (lbc *RemoteDecitionConnection)writeToSocket(){
	sock := lbc.Sock
	if sock == nil{
		return;
	}
	_,isOk := <- lbc.writeChan;
	if false == isOk{
		return;
	}
	j := len(lbc.waitSendMSGBuffer);
	if j <= 0{
		lbc.writeChan <- 1;
		return;
	}
	//取一条被写数据
	arg := lbc.waitSendMSGBuffer[0];
	lbc.waitSendMSGBuffer = lbc.waitSendMSGBuffer[1:j];
	//执行写入
	sock.SetDeadline(time.Now().Add(socketTimeoutSecond));//延长超时时间
	data,err := json.Marshal(arg);
	tj := len(data);
	if err != nil || tj <= 2{
		log.Println("sock:",lbc.SID,"数据转换出错:",err.Error())
		lbc.writeChan <- 1
		return;
	}
	data = append([]byte(pkgHead),data...);
	data = append(data,[]byte(pkgFoot)...);
	fmt.Println("sock:",lbc.SID," 准备写入socket的数据:",string(data));
	n,err := sock.Write(data);
	_ = n;
	if err != nil{
		lbc.writeChan <- 1;
		if operr,ok :=err.(*net.OpError);ok == true{
			if operr.Timeout() == true && lbc.OnTimeout != nil{
				lbc.OnTimeout(lbc);//socket超时， 通知管理器，进行处理
			}else if lbc.OnError != nil{
				lbc.OnError(lbc);//socket错误时， 通知管理器，进行处理
			}
		}else if lbc.OnError != nil{
			lbc.OnError(lbc);//socket错误时， 通知管理器，进行处理
		}
		return;
	}
	lbc.writeChan <- 1;
}

/**
从socket中循环读取数据
*/
func (lbc *RemoteDecitionConnection)runLoopRead(){
	buffer := make([]byte,1024);//每次读取1024个字节的数据
	sock := lbc.Sock
	if nil == sock{
		return;
	}
	sock.SetDeadline(time.Now().Add(socketTimeoutSecond));//延长超时时间
	for lbc.isConnected{
		n,err := sock.Read(buffer);//调用这行代码后，当前read 携程会默认阻塞，直到读取到数据或者timeout或者err之后才会执行下面的代码
		if err != nil{
			log.Println("socket:",lbc.SID,",数据读取错误",sock.RemoteAddr().String(), " connection error: ", err.Error());
			
			if operr,ok :=err.(*net.OpError);ok == true{
				if operr.Timeout() == true && lbc.OnTimeout != nil{
					lbc.OnTimeout(lbc);//socket超时， 通知管理器，进行处理
				}else if lbc.OnError != nil{
					lbc.OnError(lbc);//socket错误时， 通知管理器，进行处理
				}
			}else if lbc.OnError != nil{
				lbc.OnError(lbc);//socket错误时， 通知管理器，进行处理
			}
			break;//跳出消息循环
		}
		data := buffer[:n];//最终读取出的数据
		if len(data) > 0{
			//延长心跳超时时间
			sock.SetDeadline(time.Now().Add(socketTimeoutSecond));
			//如果有数据则异步进行解析，这样可以利用sock.read时，处理数据包。而不会因为read后同步执行checkPage或者execPackage时间过长，无形中推迟下一次read的时间点
			_,isOk := <- lbc.readChan;
			if false == isOk{
				break;
			}
			lbc.readBuffer = append(lbc.readBuffer,data...);
			lbc.readChan <- 1
			go lbc.checkPackage();
		}
	}
}

/*在readbuffer中检索有用的数据包*/
func (lbc *RemoteDecitionConnection)checkPackage(){
	/*
	以下的处理环节存在一个问题，即无法解决第一个包为 <gmlb>"cmd":10010,"mgs":"abc"<gmle><gmlb>"cmd":10010,"mgs":"abc"
	第二个包为<gmlb>"cmd":10010,"mgs":"abc"<gmle>的丢包情况。 这种情况下只会处理第一个包，而第二个包会被丢弃，视作丢包
	*/
	_,isOk := <- lbc.readChan;
	if false == isOk{
		return;
	}
	str := string(lbc.readBuffer);
	waitDelStr := "";
	j := strings.Index(str,pkgHead);//取包头位置
	pos := 0;
	//解析并做粘包处理
	for lbc.isConnected == true && j >= 0{
		//删除包头前无用的字节
		if j > 0{
			waitDelStr = str[0:j];//取被删除的字符传
			pos = len([]byte(waitDelStr));//取被删除的字符传所占字节的长度
			lbc.readBuffer = lbc.readBuffer[pos:];//更新 “读缓存”,从“读缓存”的开头删除<gmlb>标识之前的无用字节
			str = str[j:];//更新当前要处理的 "全量字符传"
			j = 0;//更新当前J的位置
		}
		//取包尾标识，如果有则解析包，没有则跳出循环
		j = strings.Index(str,pkgFoot);
		if j > 0{
			j += len(pkgFoot);
			waitDelStr = str[0:j];//取被删除的字符传
			pos = len([]byte(waitDelStr));//取被删除的字符传所占字节的长度
			lbc.readBuffer = lbc.readBuffer[pos:];//更新 “读缓存”
			str = str[j:];//更新当前要处理的 "全量字符传"
			j = strings.Index(str,pkgHead);//更新当前J的位置到下一个包头的位置
			lbc.execPackage(waitDelStr);//执行数据包
			continue;
		}
		break;
	}
	lbc.readChan <- 1
}

/*处理客户端发来的数据包*/
func (lbc *RemoteDecitionConnection)execPackage(pkgStr string){
	jsonStr := strings.Replace(pkgStr,pkgHead,"",-1);
	jsonStr = strings.Replace(jsonStr,pkgFoot,"",-1);
	lbc._execPackage([]byte(jsonStr));//交给 数据包处理者进行处理
}


/**
数据包处理者
*/

func (client *RemoteDecitionConnection)_execPackage(jsonByte []byte){
	var jsonObj map[string]interface{};
	err := json.Unmarshal(jsonByte,&jsonObj);
	if err == nil{
		//取cmd，并决策执行
		cmd := jsonObj["cmd"];
		if temp,ok := cmd.(float64);ok == true{
			command := uint32(temp);
			fmt.Println("sid:",client.SID," 收到数据包的cmd:",command);
			switch command{
			case model.C_REQ_S_HEARTBEAT:
				//返回服务端心跳
				if resObj := CreateProtocal(model.S_RES_C_HEARTBEAT);resObj != nil{
					client.Write(resObj);
				}
				break;
			case model.C_REQ_S_JOINROOM:
				c2s_login(client,jsonByte);
				break;
			case model.C_REQ_S_LEAVEROOM:
				c2s_logout(client,jsonByte);
				break;
			case model.C_REQ_S_STRATEGYCHANGED:
				c2s_StrategyChanaged(client,jsonByte);
				break;
			default:
				break;
			}
		}
	}else{
		log.Println("sid:",client.SID," 数据包解析错误:",err.Error());
	}
	
}

/**
处理进入教室
*/
func c2s_login(client *RemoteDecitionConnection,jsonByte []byte){
	var req model.JoinRoom_c2s;
	err := json.Unmarshal(jsonByte,&req);
	if err == nil{
		uid := req.Uid;
		preClient := OwnedConnectUIDMap_GetValue(uid);//根据UID获取当前用户已经进入教室的socket连接，正常情况下应为nil
		if preClient == nil{
			//进入教室
			login(client,req);
		}else{
			if preClient == client{
				//重复登录的情况
					//先调用登出
					logout(client,uid);
					//后登录
					login(client,req);
			}else{
				//不同的socket，之前的socket已经存在于教室，则将其踢出
				DestroySocket(preClient,func(){
					login(client,req);
				})
			}
		}
	}
}

/**进入教室*/
func login(client *RemoteDecitionConnection,req model.JoinRoom_c2s){
	result := false;
	tempRes := CreateProtocal(model.S_RES_C_JOINROOM);
	if tempRes == nil{
		return ;
	}

	res,ok := tempRes.(*model.JoinRoom_s2c);
	if ok{
		if req.Uid <= 0{
			res.Code = 262;
			res.FaildMsg = "进入room失败,uid无效";
		}else{
			//将各种ID绑定到当前的socket上
			client.UID = req.Uid;
			//封装请求端的回执信息,并发送
			res.Code = 0;
			res.FaildMsg = "";
			client.Write(res);
			result = true;

			if client.CurrentStrategyConfigPath != "" && client.CurrentConditionConfigPath != ""{
				//发送策略变更协议，便于客户端更新策略
				strategyRes := &model.StrategyChanged_s2c_notify{ConditionPath:client.CurrentConditionConfigPath,StrategyPath:client.CurrentStrategyConfigPath,Msg:"{}"}
				//通知所有客户端
				client.Write(strategyRes);
			}
			
		}
	}
	if result == true{
		//全新的用户进入教室
		NewUserClientlogin(client.SID,req.Uid,client);
	}
}


func c2s_StrategyChanaged(client *RemoteDecitionConnection,jsonByte []byte){
	var req model.StrategyChanged_c2s;
	err := json.Unmarshal(jsonByte,&req);
	if err == nil{
		client.CurrentStrategyConfigPath = req.ConditionPath;
		client.CurrentConditionConfigPath = req.StrategyPath;
		res := &model.StrategyChanged_s2c_notify{ConditionPath:req.ConditionPath,StrategyPath:req.StrategyPath,Msg:req.Msg}
		//通知所有客户端
		for _,sock := range ownedConnect{
			sock.Write(res)
		}
	}
}

/**
处理离开教室
*/
func c2s_logout(client *RemoteDecitionConnection,jsonByte []byte){
	var req model.LeaveRoom_c2s;
	err := json.Unmarshal(jsonByte,&req);
	if err == nil{
		logout(client,req.Uid);
	}
}

/**
离开教室
*/
func logout(client *RemoteDecitionConnection,uid int64){
	UnOwnedConnect_SetValue(client.SID,client);
	OwnedConnect_SetValue(client.SID,nil);
	OwnedConnectUIDMap_SetValue(uid,nil);
	temp := CreateProtocal(model.S_RES_C_LEAVEROOM);
	if temp != nil{
		if res,ok := temp.(*model.LeaveRoom_s2c);ok == true{
			res.Uid = client.UID;
			client.Write(res);
		}
	}
	client.UID = -1;
}
package luboService;
import(
	"net"
	"time"
	"log"
	model "../models"
	"encoding/json"
	_"fmt"
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
type LuBoClientConnection struct{
	isConnected bool/*是否处于连接状态*/
	writeChan chan int;/*socket写入操作，锁*/
	waitSendMSGBuffer []interface{};/*即将发送给客户端的消息的队列*/
	readChan chan int;/*数据读取锁*/
	readBuffer []byte;/*数据读取队列*/
	SID int64 /*socket id*/
	UID int64 /*用户ID默认为-1*/
	RID int64 /*教室ID 默认为-1*/
	Sock net.Conn;/*真正的socket操作对象*/
	TimeoutSecond time.Duration;/*超时时长*/
	OnTimeout func (*LuBoClientConnection);/*当当前socket超时时触发*/
	OnError func (*LuBoClientConnection);/*当当前socket 发生除了Timeout错误以外，如EOF时触发*/
	OnSocketCloseComplete func(...interface{});//当当前socket连接被close完毕后触发的处理函数。
	OnSocketCloseCompleteArgs []interface{};//OnSocketCloseComplete函数的参数
}

/**
创建新的socket管理工具
*/
func NewLuBoClientConn(sid int64,conn net.Conn)(*LuBoClientConnection){
	client := &LuBoClientConnection{SID:sid,UID:-1,RID:-1,Sock:conn,TimeoutSecond:socketTimeoutSecond,isConnected:true};
	client.writeChan = make(chan int,1);//初始化 socket写入操作锁
	client.readChan = make(chan int,1);
	go client.runLoopRead();//开socket read队列
	//封装并回执给客户端“客户端接入”信令
	pro := CreateProtocal(model.S_RES_C_HY);
	client.Write(pro);
	return client;
}

/**
关闭并释放socket,支持在关闭socket之前，依然向原有socket中写入一条消息
*/
func (lbc *LuBoClientConnection)DestroySocket(arg interface{}){
	if lbc.isConnected == false{
		return;
	}
	lbc.isConnected = false;
	lbc.OnTimeout = nil;
	lbc.OnError = nil;
	lbc.UID = -1;
	lbc.RID = -1;
	go lbc.writeLastMsgAndCloseSock(arg);//送最后一条消息后，关闭socket
	
}

/*发送最后一条数据给客户端，之后关闭socket*/
func (lbc *LuBoClientConnection)writeLastMsgAndCloseSock(arg interface{}){
	sock := lbc.Sock
	if sock == nil{
		return;
	}
	lbc.writeChan <- 1;
	//执行写入
	sock.SetDeadline(time.Now().Add(socketTimeoutSecond));//延长超时时间
	data,err := json.Marshal(arg);
	if err != nil{
		log.Println("sock:",lbc.SID,"数据转换出错:",err.Error())
		<- lbc.writeChan
		lbc.closeSocket();
		return;
	}
	log.Println("sock:",lbc.SID," 准备写入socket的数据:",string(data));
	n,err := sock.Write(data);
	_ = n;
	if err != nil{
		<- lbc.writeChan
		lbc.closeSocket();
		return;
	}
	//log.Println("sock:",lbc.SID," 已发送的数据长度",n);
	<- lbc.writeChan
	lbc.closeSocket();
}

/*
发送数据给客户端
*/
func (lbc *LuBoClientConnection)Write(arg interface{}){
	lbc.writeChan <- 1
	lbc.waitSendMSGBuffer = append(lbc.waitSendMSGBuffer,arg);//向 “等待发送的消息队列中”添加一条消息
	<- lbc.writeChan
	go lbc.writeToSocket();//调用消息发送函数
}

/*检查是否需要关闭客户端socket连接*/
func (lbc *LuBoClientConnection)closeSocket(){
	lbc.SID = -1;
	lbc.Sock.Close();
	lbc.Sock = nil;
	close(lbc.writeChan)//关闭channel
	close(lbc.readChan)
	if lbc.OnSocketCloseComplete != nil{
		if lbc.OnSocketCloseCompleteArgs != nil{
			lbc.OnSocketCloseComplete(lbc.OnSocketCloseCompleteArgs ...);
		}else{
			lbc.OnSocketCloseComplete();
		}
		lbc.OnSocketCloseComplete = nil;
	}
}

func (lbc *LuBoClientConnection)writeToSocket(){
	sock := lbc.Sock
	if sock == nil{
		return;
	}
	lbc.writeChan <- 1;
	j := len(lbc.waitSendMSGBuffer);
	if j <= 0{
		<- lbc.writeChan;
		return;
	}
	//取一条被写数据
	arg := lbc.waitSendMSGBuffer[0];
	lbc.waitSendMSGBuffer = lbc.waitSendMSGBuffer[1:j];
	//执行写入
	sock.SetDeadline(time.Now().Add(socketTimeoutSecond));//延长超时时间
	data,err := json.Marshal(arg);
	if err != nil{
		log.Println("sock:",lbc.SID,"数据转换出错:",err.Error())
		<- lbc.writeChan
		return;
	}
	log.Println("sock:",lbc.SID," 准备写入socket的数据:",string(data));
	n,err := sock.Write(data);
	_ = n;
	if err != nil{
		<- lbc.writeChan
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
	//log.Println("sock:",lbc.SID," 已发送的数据长度",n);
	<- lbc.writeChan
}

/**
从socket中循环读取数据
*/
func (lbc *LuBoClientConnection)runLoopRead(){
	buffer := make([]byte,1024);//每次读取1024个字节的数据
	sock := lbc.Sock
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
			lbc.readChan <- 1;
			lbc.readBuffer = append(lbc.readBuffer,data...);
			<- lbc.readChan
			go lbc.checkPackage(data);
		}
	}
}

/*在readbuffer中检索有用的数据包*/
func (lbc *LuBoClientConnection)checkPackage(data []byte){
	/*
	以下的处理环节存在一个问题，即无法解决第一个包为 <gmlb>"cmd":10010,"mgs":"abc"<gmle><gmlb>"cmd":10010,"mgs":"abc"
	第二个包为<gmlb>"cmd":10010,"mgs":"abc"<gmle>的丢包情况。 这种情况下只会处理第一个包，而第二个包会被丢弃，视作丢包
	*/
	lbc.readChan <- 1;
	str := string(lbc.readBuffer);
	//fmt.Println(str);
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
	<- lbc.readChan;
}

/*处理客户端发来的数据包*/
func (lbc *LuBoClientConnection)execPackage(pkgStr string){
	jsonStr := strings.Replace(pkgStr,pkgHead,"{",-1);
	jsonStr = strings.Replace(jsonStr,pkgFoot,"}",-1);
	ExecPackage(lbc,[]byte(jsonStr));//交给 数据包处理者进行处理
}
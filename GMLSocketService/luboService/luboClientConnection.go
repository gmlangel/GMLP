package luboService;
import(
	"net"
	"time"
	"log"
	model "../models"
	"encoding/json"
	_"fmt"
)

const(
	 socketTimeoutSecond time.Duration= time.Second * 60;
)
/**
用于管理单个客户端的Socket连接
*/
type LuBoClientConnection struct{
	willCloseSocket bool;//是否需要延迟关闭socket
	isConnected bool/*是否处于连接状态*/
	writeChan chan int;/*socket写入操作，锁*/
	waitSendMSGBuffer []interface{};/*即将发送给客户端的消息的队列*/
	SID int64 /*socket id*/
	UID int64 /*用户ID默认为-1*/
	RID int64 /*教室ID 默认为-1*/
	Sock net.Conn;/*真正的socket操作对象*/
	TimeoutSecond time.Duration;/*超时时长*/
	OnTimeout func (*LuBoClientConnection);/*当当前socket超时时触发*/
	OnError func (*LuBoClientConnection);/*当当前socket 发生除了Timeout错误以外，如EOF时触发*/
	OnSocketCloseComplete func(v ...interface{});//当当前socket连接被close完毕后触发的处理函数。
}

/**
创建新的socket管理工具
*/
func NewLuBoClientConn(sid int64,conn net.Conn)(*LuBoClientConnection){
	client := &LuBoClientConnection{willCloseSocket:false,SID:sid,UID:-1,RID:-1,Sock:conn,TimeoutSecond:socketTimeoutSecond,isConnected:true};
	client.writeChan = make(chan int,1);//初始化 socket写入操作锁
	go client.runLoopRead();//开socket read队列
	//封装并回执给客户端“客户端接入”信令
	pro := CreateProtocal(model.S_RES_C_HY);
	client.Write(pro);
	return client;
}

/**
关闭并释放socket,支持在关闭socket之前，依然向原有socket中写入一条消息
*/
func (lbc *LuBoClientConnection)DestroySocket(str string){
	if lbc.isConnected == false{
		return;
	}
	lbc.isConnected = false;
	lbc.OnTimeout = nil;
	lbc.OnError = nil;
	lbc.OnSocketCloseComplete = nil;
	lbc.UID = -1;
	lbc.RID = -1;
	lbc.willCloseSocket = true;
	lbc.Write(str);
	
}
func (lbc *LuBoClientConnection)Write(arg interface{}){
	lbc.writeChan <- 1
	lbc.waitSendMSGBuffer = append(lbc.waitSendMSGBuffer,arg);//向 “等待发送的消息队列中”添加一条消息
	<- lbc.writeChan
	go lbc.writeToSocket();//调用消息发送函数
}

/*检查是否需要关闭客户端socket连接*/
func (lbc *LuBoClientConnection)checkNeedCloseSocket(){
	if lbc.willCloseSocket == true{
		lbc.SID = -1;
		lbc.Sock.Close();
		lbc.Sock = nil;
		close(lbc.writeChan)//关闭channel
	}
}

func (lbc *LuBoClientConnection)writeToSocket(){
	lbc.writeChan <- 1;
	j := len(lbc.waitSendMSGBuffer);
	if j <= 0{
		<- lbc.writeChan;
		lbc.checkNeedCloseSocket();
		return;
	}
	//取一条被写数据
	arg := lbc.waitSendMSGBuffer[0];
	lbc.waitSendMSGBuffer = lbc.waitSendMSGBuffer[1:j];
	//执行写入
	sock := lbc.Sock
	sock.SetDeadline(time.Now().Add(socketTimeoutSecond));//延长超时时间
	data,err := json.Marshal(arg);
	if err != nil{
		log.Println("sock:",lbc.SID,"数据转换出错:",err.Error())
		<- lbc.writeChan
		lbc.checkNeedCloseSocket();
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
		lbc.checkNeedCloseSocket();
		return;
	}
	//log.Println("sock:",lbc.SID," 已发送的数据长度",n);
	<- lbc.writeChan
	lbc.checkNeedCloseSocket();
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
			//如果有数据则进行解析
			
			
		}
	}
}
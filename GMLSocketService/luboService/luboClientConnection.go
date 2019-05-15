package luboService;
import(
	"net"
	"time"
	"log"
)

const(
	 socketTimeoutSecond time.Duration= time.Second * 60;
)
/**
用于管理单个客户端的Socket连接
*/
type LuBoClientConnection struct{
	isConnected bool/*是否处于连接状态*/
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
	client := &LuBoClientConnection{SID:sid,UID:-1,RID:-1,Sock:conn,TimeoutSecond:socketTimeoutSecond,isConnected:true};
	go client.runLoopRead();//开socket read队列
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
	lbc.write(str);
	lbc.Sock.Close();
	lbc.OnTimeout = nil;
	lbc.OnError = nil;
	lbc.OnSocketCloseComplete = nil;
	lbc.SID = -1;
	lbc.UID = -1;
	lbc.RID = -1;
	lbc.Sock = nil;
}

func (lbc *LuBoClientConnection)write(str string){

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
					lbc.OnTimeout(lbc);
				}else if lbc.OnError != nil{
					lbc.OnError(lbc);
				}
			}else if lbc.OnError != nil{
				lbc.OnError(lbc);
			}
			break;//跳出消息循环
		}
		data := buffer[:n];//最终读取出的数据
		if len(data) > 0{
			//延长心跳超时时间
			lbc.Sock.SetDeadline(time.Now().Add(socketTimeoutSecond));
			//如果有数据则进行解析
			
			
		}
	}
}
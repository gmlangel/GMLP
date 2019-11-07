package front

import(
	"net"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	models "../models"
	"log"
)

var(
	socketTimeoutSecond time.Duration= time.Second * 60;
	heartBeatTimeSecond time.Duration= time.Second * 50
	 pkgHead = "<gmlb>";//包头
	 pkgFoot = "<gmle>";//包尾
)
func NewRDSocket()*RDSocket{
	sockStruct := &RDSocket{}
	return sockStruct;
}

type RDSocket struct{
	sock net.Conn
	isConnected bool;//sock是否已经连接成功
	readBuffer []byte;
	waitSendMSGBuffer []interface{};
	linkType string;//socket连接类型
	url string;//socket连接地址
}

func(rd *RDSocket)Start(linkType string,url string){
	if rd.isConnected == true{
		return;
	}
	var err error;
	rd.linkType = linkType;
	rd.url = url;
	rd.sock, err = net.Dial(linkType, url);
    if err != nil {
		fmt.Println("dial error:", err.Error())
		rd.isConnected = false;
		go rd.relink()
		return;
	}
	rd.isConnected = true;
	fmt.Println(":connect to server ok")
	go rd.startHeartBeat();//开启心跳
	rd.runLoopRead();//开启消息循环
}

//心跳机制
func(rd *RDSocket)startHeartBeat(){
	for rd.isConnected == true{
		//发送心跳数据到远端
		req := models.HeartBeat_c2s{Cmd:0x00FF0001};
		rd.Write(req);
		time.Sleep(heartBeatTimeSecond);
	}
}

func(rd *RDSocket)close(){
	if rd.isConnected == true{
		//如果之前处于连接状态，则关闭。  如果不加这个判断回出现上一次未连接成功，本次关闭，崩溃。
		rd.sock.Close();//关闭之前的sock
	}
	
	rd.isConnected = false;
	rd.readBuffer = []byte{};
	rd.waitSendMSGBuffer = []interface{}{};
}

//重连机制
func(rd *RDSocket)relink(){
	//延迟一定时间后，重新连接
	time.Sleep(time.Second * 10);
	fmt.Println("socket重连中");
	rd.Start(rd.linkType,rd.url)
}

/**
从socket中循环读取数据
*/
func (rd *RDSocket)runLoopRead(){
	buffer := make([]byte,1024);//每次读取1024个字节的数据
	sock := rd.sock
	if nil == sock{
		return;
	}
	sock.SetDeadline(time.Now().Add(socketTimeoutSecond));//延长超时时间
	for rd.isConnected{
		n,err := sock.Read(buffer);//调用这行代码后，当前read 携程会默认阻塞，直到读取到数据或者timeout或者err之后才会执行下面的代码
		if err != nil{
			log.Println("socket:,数据读取错误",sock.RemoteAddr().String(), " connection error: ", err.Error());
			rd.close();
			go rd.relink();//延时重连
			// if operr,ok :=err.(*net.OpError);ok == true{
			// 	if operr.Timeout() == true && rd.OnTimeout != nil{
			// 		rd.relink();//重连
			// 	}else if rd.OnError != nil{
			// 		rd.relink();//演示重连
			// 	}
			// }else if rd.OnError != nil{
			// 	rd.OnError(lbc);//socket错误时， 通知管理器，进行处理
			// }
			break;//跳出消息循环
		}
		data := buffer[:n];//最终读取出的数据
		if len(data) > 0{
			//延长心跳超时时间
			sock.SetDeadline(time.Now().Add(socketTimeoutSecond));
			rd.readBuffer = append(rd.readBuffer,data...);
			rd.checkPackage();
		}
	}
}

/*在readbuffer中检索有用的数据包*/
func (rd *RDSocket)checkPackage(){
	/*
	以下的处理环节存在一个问题，即无法解决第一个包为 <gmlb>"cmd":10010,"mgs":"abc"<gmle><gmlb>"cmd":10010,"mgs":"abc"
	第二个包为<gmlb>"cmd":10010,"mgs":"abc"<gmle>的丢包情况。 这种情况下只会处理第一个包，而第二个包会被丢弃，视作丢包
	*/
	
	str := string(rd.readBuffer);
	waitDelStr := "";
	j := strings.Index(str,pkgHead);//取包头位置
	pos := 0;
	//解析并做粘包处理
	for rd.isConnected == true && j >= 0{
		//删除包头前无用的字节
		if j > 0{
			waitDelStr = str[0:j];//取被删除的字符传
			pos = len([]byte(waitDelStr));//取被删除的字符传所占字节的长度
			rd.readBuffer = rd.readBuffer[pos:];//更新 “读缓存”,从“读缓存”的开头删除<gmlb>标识之前的无用字节
			str = str[j:];//更新当前要处理的 "全量字符传"
			j = 0;//更新当前J的位置
		}
		//取包尾标识，如果有则解析包，没有则跳出循环
		j = strings.Index(str,pkgFoot);
		if j > 0{
			j += len(pkgFoot);
			waitDelStr = str[0:j];//取被删除的字符传
			pos = len([]byte(waitDelStr));//取被删除的字符传所占字节的长度
			rd.readBuffer = rd.readBuffer[pos:];//更新 “读缓存”
			str = str[j:];//更新当前要处理的 "全量字符传"
			j = strings.Index(str,pkgHead);//更新当前J的位置到下一个包头的位置
			rd.execPackage(waitDelStr);//执行数据包
			continue;
		}
		break;
	}
}

/*处理客户端发来的数据包*/
func (rd *RDSocket)execPackage(pkgStr string){
	jsonStr := strings.Replace(pkgStr,pkgHead,"",-1);
	jsonStr = strings.Replace(jsonStr,pkgFoot,"",-1);
	rd._execPackage([]byte(jsonStr));//交给 数据包处理者进行处理
}


/**
数据包处理者
*/

func (rd *RDSocket)_execPackage(jsonByte []byte){
	var jsonObj map[string]interface{};
	err := json.Unmarshal(jsonByte,&jsonObj);
	if err == nil{
		//取cmd，并决策执行
		cmd := jsonObj["cmd"];
		if temp,ok := cmd.(float64);ok == true{
			command := uint32(temp);
			fmt.Println(" 收到数据包的cmd:",command);
			switch command{
			case 0x00FF0002:
				//服务端返回心跳
				break;
			default:
				break;
			}
		}
	}else{
		fmt.Println(" 数据包解析错误:",err.Error());
	}
	
}


/*
发送数据给客户端
*/
func (rd *RDSocket)Write(arg interface{}){
	rd.waitSendMSGBuffer = append(rd.waitSendMSGBuffer,arg);//向 “等待发送的消息队列中”添加一条消息
	rd.writeToSocket();//调用消息发送函数
}

func (rd *RDSocket)writeToSocket(){
	sock := rd.sock
	if sock == nil{
		return;
	}
	j := len(rd.waitSendMSGBuffer);
	if j <= 0{
		return;
	}
	currentIdx := 0;
	//取一条被写数据
	for i,arg := range rd.waitSendMSGBuffer{
		currentIdx = i+1;
		//执行写入
		sock.SetDeadline(time.Now().Add(socketTimeoutSecond));//延长超时时间
		data,err := json.Marshal(arg);
		tj := len(data);
		if err != nil || tj <= 2{
			fmt.Println("数据转换出错:",err.Error())
			break;;
		}
		data = append([]byte(pkgHead),data...);
		data = append(data,[]byte(pkgFoot)...);
		log.Println(" 准备写入socket的数据:",string(data));
		n,err := sock.Write(data);
		_ = n;
		if err != nil{
			// if operr,ok :=err.(*net.OpError);ok == true{
			// 	if operr.Timeout() == true && rd.OnTimeout != nil{
			// 		rd.OnTimeout(lbc);//socket超时， 通知管理器，进行处理
			// 	}else if rd.OnError != nil{
			// 		rd.OnError(lbc);//socket错误时， 通知管理器，进行处理
			// 	}
			// }else if rd.OnError != nil{
			// 	rd.OnError(lbc);//socket错误时， 通知管理器，进行处理
			// }
			fmt.Println("socket:,数据写入错误",sock.RemoteAddr().String(), " connection error: ", err.Error());
			rd.close();
			go rd.relink();//延时重连
			break;
		}
	}

	rd.waitSendMSGBuffer = rd.waitSendMSGBuffer[currentIdx:j];
	
}
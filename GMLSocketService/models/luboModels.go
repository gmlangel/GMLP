package models;

type LoBoServerConfig struct{
	Host string;/*地址*/
	Port uint32;/*端口*/
	ServerType string;/*tcp or  udp*/
}

//协议ID定义-----------------------------
//MARK:Socket 通讯 数据包命令ID定义
const(
	/**服务器返回'欢迎'信令*/
    S_RES_C_HY = 0x00FF0000;
    //心跳服务
    C_REQ_S_HEARTBEAT = 0x00FF0001;
    S_RES_C_HEARTBEAT = 0x00FF0002;
    
    //登录服务
    C_REQ_S_LOGIN = 0x00FF0003;
    S_RES_C_LOGIN = 0x00FF0004;
    
    //登出服务
    C_REQ_S_LOGOUT = 0x00FF0005;
    S_RES_C_LOGOUT = 0x00FF0006;
    
    //掉线通知
    S_NOTIFY_C_OFFLINE = 0x00FF0007;
    
    //获取用户信息
    C_REQ_S_GETUSERINFO = 0x00FF0008;
    S_RES_C_GETUSERINFO = 0x00FF0009;
    
    //更新用户信息
    C_REQ_S_UPDATEUSERINFO = 0x00FF000A;
    S_RES_C_UPDATEUSERINFO = 0x00FF000B;
    
    //创建教室
    C_REQ_S_CREATEROOM = 0x00FF000C;
    S_RES_C_CREATEROOM = 0x00FF000D;
    
    //获取用户创建的教室信息
    C_REQ_S_GETROOMSINFOBYUSER = 0x00FF000E;
    S_RES_C_GETROOMSINFOBYUSER = 0x00FF000F;
    
    //删除教室
    C_REQ_S_DELROOM = 0x00FF0011;
    S_RES_C_DELROOM = 0x00FF0012;
    
    //教室状态变更通知
    S_NOTIFY_C_ROOMSTATECHANGE = 0x00FF0013;
    
    //进入教室
    C_REQ_S_JOINROOM = 0x00FF0014;
    S_RES_C_JOINROOM = 0x00FF0015;
    
    //离开教室
	C_REQ_S_LEAVEROOM = 0x00FF0016;
	S_RES_C_LEAVEROOM = 0x00FF0026;
    
    //其他人状态信息变更通知
    S_NOTIFY_C_OTHERUSERSTATECHANGE = 0x00FF0017;
    
    //发送文本消息
    C_REQ_S_SENDCHAT = 0x00FF0018;
    
    //文本消息通知
    S_NOTIFY_C_CHAT = 0x00FF0019;
    
    //发送管理员命令
    C_REQ_S_SENDADMINCMD = 0x00FF001A;
    
    //管理员命令通知
    S_NOTIFY_C_ADMINCMD = 0x00FF001B;
    
    //上报答题结果
	C_REQ_S_UPLOADANSWERCMD = 0x00FF001C;
	S_RES_C_UPLOADANSWERCMD = 0x00FF002C;
    
    //收到服务器下发的教学脚本
    S_NOTIFY_C_TEACHSCRIPTCMD = 0x00FF001D;
    
    //请求课程报告
    C_REQ_S_USERLESSONRESULT = 0x00FF001E;
    
    //收到服务器下发的用户课程报告
    S_RES_C_USERLESSONRESULT = 0x00FF001F;
    
    //服务器下发通知 教学脚本加载完毕
	S_NOTIFY_C_TEACHSCRIPTLOADEND = 0x00FF0020;
	
	//客户端数据上报
    C_REQ_S_UPLOADDATA = 0x00FF003A;
    S_RES_C_UPLOADDATA = 0x00FF003B;
)

/*新客户端socket介入后的 欢迎消息反馈*/
type ClientIn_s2c struct{
	Cmd uint32 `json:"cmd"`
	Des string `json:"des"`
}

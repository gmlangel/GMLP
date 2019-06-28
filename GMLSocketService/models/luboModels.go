package models;
import(
    "time"
)


const(
    RoomState_NotStart Enum_RoomState = "nostart";//课程状态   未开始
    RoomState_Started Enum_RoomState = "started";//课程状态   已开始
    RoomState_End Enum_RoomState = "end";//课程状态   已结束
    TeachScriptTimeInterval = time.Second * time.Duration(1);//教材脚本下发的时间间隔
)

type Enum_RoomState string;

type LoBoServerConfig struct{
	Host string;/*地址*/
	Port uint32;/*端口*/
	ServerType string;/*tcp or  udp*/
}


type RoomInfo struct{
    //通用---begin
    Rid int64;
    RoomState Enum_RoomState;//课程状态
    StartTimeInterval int64;//课程开始时间beginTime
    TeachingTmaterialScriptID int64;//教学脚本ID
    UserArr []CurrentUser;//当前频道中的人的信息数组
    UserIdArr []int64;//用户ID数组
    AnswerUIDQueue []int64;//用户答题序列数组
    TongyongCMDArr []map[string]interface{};//通用教学命令
    MainFrames []MediaMainFrame;//主媒体播放轴
    CurrentAnswerState string;//当前问题的用户答题结果 faild success timeouterr
    CurrentProcess int ;//当前正在处理哪个流程 。  0 代表 正在处理媒体播放流程   1代表正在处理关键帧处理流程
    //通用---end

    //媒体播放相关---begin
    MCurrent *MediaData;//当前正在执行的媒体脚本
    MAllowNew bool;//是否允许下发新的媒体脚本
    MCompleteTime int64;//当前媒体脚本的播放结束时间，应与MCurrentTimeInterval进行比对
    MCurrentTimeInterval int64;//某一段媒体脚本已经执行了的时间，用于进行各种时间比对及计算
    MCurrentMainFrameIdx int64;//当前播放视频对应的关键帧数组中 当前正在播放的关键帧的索引
    //媒体播放相关---end

    //教学脚本相关---begin
    SCurrentTimeInterval int64;//某一段教学脚本已经执行了的时间，用于进行各种时间比对及计算
    SAllowNew bool;//是否允许下发新的教学脚本
    SCurrent *ScriptStepData;//当前正在执行的教学脚本
    SCurrentQuestionId int64;//当前等待应答的问题的ID
    SCurrentQuesionTimeOut int64;//关键帧对应的脚本执行所需的超时时间他和SCurrentTimeInterval可以直接进行计算
    SWaitAnswerUids []int64;//等待做答的用户ID数组,只有swaitAnswerUids长度为0，才可以执行SCurrent对应的下一个教学脚本或者重置MAllowNew为true
    //教学脚本相关---end
    
}

/**媒体脚本*/
type MediaData struct{
    Id int64 `json:"id"`;//唯一标识
    Type string `json:"type"`;//脚本类型
    Next int64 `json:"n"`;//下一个执行的脚本,-1代表不存在
    Pre int64 `json:"p"`;//上一个执行的脚本,-1代表不存在
    Value map[string]interface{} `json:"value"`;//详细数据
    MainFrames []MediaMainFrame `json:"mainFrames"`;//媒体关键帧数组
}

type MediaMainFrame struct{
    MediaTime int64 `json:"mt"`;//当前脚本关键帧时间戳对应的媒体关键帧时间戳，用于退出教室后重新进入教室续播
    StepTime int64 `json:"st"`;//脚本关键帧时间戳
    StepId int64 `json:"sid"`;//对应的教学脚本id
}

/*教学脚本*/
type ScriptStepData struct{
    Id int64 `json:"id"`;//唯一标识
    Type string `json:"type"`;//脚本类型
    Next int64 `json:"n"`;//下一个执行的脚本,-1代表不存在
    Pre int64 `json:"p"`;//上一个执行的脚本,-1代表不存在
    Value map[string]interface{} `json:"value"`;//详细数据
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
    C_REQ_S_UPLOADREPORTDATA = 0x00FF003A;
    S_RES_C_UPLOADREPORTDATA = 0x00FF003B;
)

/*新客户端socket介入后的 欢迎消息反馈*/
type ClientIn_s2c struct{
	Cmd uint32 `json:"cmd"`;
	Description string `json:"des"`;//描述
}

/*客户端心跳*/
type HeartBeat_c2s struct{
    Cmd uint32 `json:"cmd"`
    Seq uint32 `json:"seq"`;//数据包的序号，可以为0
    LocalTimeinterval uint32 `json:"lt"`;//客户端发送请求时的UTC时间的秒值
}

type HeartBeat_s2c struct{
    Cmd uint32 `json:"cmd"`
    Seq uint32 `json:"seq"`;//数据包的序号，可以为0
    C_Seq uint32 `json:"c_seq"`;//数据包的序号，可以为0
    Servertime uint32 `json:"st"`;//服务器端UTC时间的秒值
}

/*客户端进入教室*/
type JoinRoom_c2s struct{
    Cmd uint32 `json:"cmd"`
    Seq uint32 `json:"seq"`;//数据包的序号，可以为0
    Rid int64 `json:"rid"`;//教室ID
    TeachScriptID int64 `json:"tts"`;//课程教学脚本的ID
    StartTimeinterval int64 `json:"sti"`;//课程开始时间的UTC时间戳秒值
    Uid int64 `json:"uid"`;//用户ID
    NickName string `json:"nn"`;//用户昵称
}

type JoinRoom_s2c struct{
    Cmd uint32 `json:"cmd"`
    C_Seq uint32 `json:"c_seq"`;//数据包的序号，可以为0
    Rid int64 `json:"rid"`;//教室ID
    Code uint32 `json:"code"`;//进入教室是否成功0 = 成功 ,262 = 进入room失败,uid无效,263 = 进入room失败,roomId小于0,无效
    FaildMsg string `json:"fe"`;//报错信息
    UserArr []CurrentUser `json:"ua"`;//用户列表
}

type CurrentUser struct{
    Uid int64 `json:"uid"`;//用户ID
    NickName string `json:"nn"`;//用户昵称
    Type bool `json:"type"`;//用户进出教室状态 True = 进入教室,False = 离开教室
}

/**
教学脚本加载完毕,服务器下发课程教学脚本中的resource config部分的数据
*/
type TeachScriptLoadEnd_s2c_notify struct{
    Cmd uint32 `json:"cmd"`
    Code uint32 `json:"code"`;//服务端课程脚本加载完毕 0 = 成功
    FaildMsg string `json:"fe"`;//报错信息
    ScriptConfigData ScriptConfigDataMap `json:"scriptConfigData"`;//课程脚本相关配置
}

type ScriptConfigDataMap struct{
    CourseId uint32 `json:"courseId"`;//教材ID
    Width uint32 `json:"width"`;//教材原始宽度
    Height uint32 `json:"height"`;//教材原始高度
    Resource map[string]interface{} `json:"resource"`;//课程脚本中的resource config部分
}

/**
服务器下发（重新进入教室时也会下发教室内缓存的
*/
type PushTeachScriptCache_s2c_notify struct{
    Cmd uint32 `json:"cmd"`
    Rid int64 `json:"rid"`;//教室ID
    Code uint32 `json:"code"`;//暂时无意义 0 = 成功
    FaildMsg string `json:"fe"`;//报错信息
    PlayTimeInterval int64 `json:"playTimeInterval"`;//本消息中的最后一条教学脚本已经执行了的时间，秒值
    Datas []map[string]interface{} `json:"datas"`;//教学脚本数组
    AnswerUIDQueue []int64;//学员uid列表。用于1对多课程场景，暂时无用。
}

/*上报答题结果*/
type UploadAnswer_c2s struct{
    Cmd uint32 `json:"cmd"`
    Uid int64 `json:"uid"`;//用户ID
    LocalTimeinterval uint32 `json:"lt"`;//客户端发送请求时的UTC时间的秒值
    Id int64 `json:"id"`;//题号,对应0x00FF001D 消息datas消息中的id
    Data Answer `json:"data"`;//答案内容
}

type Answer struct{
    Tplate string `json:"tplate"`;//知识标签,因教材中目前还没有“知识标签”的预埋，所以暂时先传空字符串
    ReAnswerCount uint32 `json:"reAnswer"`;//学员重复作答次数。
    Score uint32 `json:"score"`;//本题得分,0=错误   100=正确，其它得分用于语音模板等特殊模板
}

type Answer_c2s struct{
    Id int64 `json:"id"`;//题号,对应0x00FF001D 消息datas消息中的id
    Data Answer `json:"data"`;//答案内容
}

type UploadAnswer_s2c struct{
    Cmd uint32 `json:"cmd"`
    Code uint32 `json:"code"`;//暂时无意义 0 = 成功
    FaildMsg string `json:"fe"`;//报错信息
}

/*客户端请求课程学习报告*/
type UserLessonResult_c2s struct{
    Cmd uint32 `json:"cmd"`
    Rid int64 `json:"rid"`;//教室ID
    Uid int64 `json:"uid"`;//用户ID
}

type UserLessonResult_s2c struct{
    Cmd uint32 `json:"cmd"`
    Rid int64 `json:"rid"`;//教室ID
    Code uint32 `json:"code"`;//暂时无意义 0 = 成功
    FaildMsg string `json:"fe"`;//报错信息
    Seq uint32 `json:"seq"`;//数据包的序号，可以为0
    Datas []Answer_c2s `json:"datas"`;//答案内容数组
}

/*客户端数据上报*/
type DataReport_c2s struct{
    Cmd uint32 `json:"cmd"`
    Rid int64 `json:"rid"`;//教室ID
    Uid int64 `json:"uid"`;//用户ID
    C_Code uint32 `json:"c_code"`;//状态码 1 = 听不到老师声音,2 = 看不到教材,3 = 教材画面卡住，不动
    Msg string `json:"msg"`;//具体数据
}

type DataReport_s2c struct{
    Cmd uint32 `json:"cmd"`
    Rid int64 `json:"rid"`;//教室ID
}

/*客户端请求离开教室*/
type LeaveRoom_c2s struct{
    Cmd uint32 `json:"cmd"`
    Rid int64 `json:"rid"`;//教室ID
    Uid int64 `json:"uid"`;//用户ID
}

type LeaveRoom_s2c struct{
    Cmd uint32 `json:"cmd"`
    Rid int64 `json:"rid"`;//教室ID
    Uid int64 `json:"uid"`;//用户ID
    Code uint32 `json:"code"`;// 0 = 成功
}

/*掉线通知*/
type OfflineNotify_s2c struct{
    Cmd uint32 `json:"cmd"`
    Seq uint32 `json:"seq"`;//数据包的序号，可以为0
    Code uint32 `json:"code"`;
    Reason string `json:"reason"`;//掉线原因
}
package luboService;
import(
	model "../models"
	"time"
)

func CreateProtocal(cmd uint32)interface{}{
	var result interface{} = nil;
	switch cmd{
	case model.S_RES_C_HY:
		result = createClientIn_s2c();
		break;
	case model.S_RES_C_HEARTBEAT:
		result = createHeartBeat_s2c()
		break;
	case model.S_RES_C_JOINROOM:
		result = createJoinRoom_s2c();
		break;
	case model.S_NOTIFY_C_TEACHSCRIPTLOADEND:
		result = createTeachScriptLoadEnd_s2c_notify();
		break;
	case model.S_NOTIFY_C_TEACHSCRIPTCMD:
		result = createPushTeachScriptCache_s2c_notify();
		break;
	case model.S_RES_C_USERLESSONRESULT:
		result = createUserLessonResult_s2c();
		break;
	case model.S_RES_C_UPLOADREPORTDATA:
		result = createDataReport_s2c();
		break;
	case model.S_RES_C_LEAVEROOM:
		result = createLeaveRoom_s2c();
		break;
	case model.S_NOTIFY_C_OFFLINE:
		result = createOfflineNotify_s2c();
		break;
	default:break;
	}
	return result;
}

func createClientIn_s2c()interface{}{
	result := &model.ClientIn_s2c{Cmd:model.S_RES_C_HY,Description:"欢迎加入AI录播服务器"};
	return result;
}


func createHeartBeat_s2c()interface{}{
	result := &model.HeartBeat_s2c{Cmd:model.S_RES_C_HEARTBEAT,Seq:0,C_Seq:0,Servertime:uint32(time.Now().Unix())};
	return result;
}

func createJoinRoom_s2c()interface{}{
	result := &model.JoinRoom_s2c{Cmd:model.S_RES_C_JOINROOM,Rid:0,Code:0,C_Seq:0,FaildMsg:""};
	return result;
}

func createTeachScriptLoadEnd_s2c_notify()interface{}{
	result := &model.TeachScriptLoadEnd_s2c_notify{Cmd:model.S_NOTIFY_C_TEACHSCRIPTLOADEND,Code:0,FaildMsg:""}
	return result;
}

func createPushTeachScriptCache_s2c_notify()interface{}{
	result := &model.PushTeachScriptCache_s2c_notify{Cmd:model.S_NOTIFY_C_TEACHSCRIPTCMD,Code:0,FaildMsg:"",Rid:0,PlayTimeInterval:0};
	return result;
}

func createUserLessonResult_s2c()interface{}{
	result := &model.UserLessonResult_s2c{Cmd:model.S_RES_C_USERLESSONRESULT,Code:0,FaildMsg:"",Rid:0,Seq:0}
	return result;
}

func createDataReport_s2c()interface{}{
	result := &model.DataReport_s2c{Cmd:model.S_RES_C_UPLOADREPORTDATA,Rid:0};
	return result;
}

func createLeaveRoom_s2c()interface{}{
	result := &model.LeaveRoom_s2c{Cmd:model.S_RES_C_LEAVEROOM,Rid:0,Uid:0,Code:0};
	return result;
}

func createOfflineNotify_s2c()interface{}{
	result := &model.OfflineNotify_s2c{Cmd:model.S_NOTIFY_C_OFFLINE,Seq:0,Code:0,Reason:""}
	return result;
}
package remoteDecitionService;
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
	result := &model.ClientIn_s2c{Cmd:model.S_RES_C_HY,Description:"欢迎加入RemoteDecition服务器"};
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

func createLeaveRoom_s2c()interface{}{
	result := &model.LeaveRoom_s2c{Cmd:model.S_RES_C_LEAVEROOM,Rid:0,Uid:0,Code:0};
	return result;
}

func createOfflineNotify_s2c()interface{}{
	result := &model.OfflineNotify_s2c{Cmd:model.S_NOTIFY_C_OFFLINE,Seq:0,Code:0,Reason:""}
	return result;
}
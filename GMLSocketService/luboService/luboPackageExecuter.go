package luboService;
import(
	"fmt"
	model "../models"
	"encoding/json"
)
/**
数据包处理者
*/

func ExecPackage(client *LuBoClientConnection,jsonByte []byte){
	var jsonObj map[string]interface{};
	err := json.Unmarshal(jsonByte,&jsonObj);
	if err == nil{
		fmt.Println(fmt.Sprintf("sid:%d 收到数据包:%v",client.SID,jsonObj));
		//取cmd，并决策执行
		cmd := jsonObj["cmd"];
		if temp,ok := cmd.(float64);ok == true{
			command := uint32(temp);
			fmt.Println("数据包的cmd:",command);
			switch command{
			case model.C_REQ_S_HEARTBEAT:
				//返回服务端心跳
				if resObj := CreateProtocal(model.S_RES_C_HEARTBEAT);resObj != nil{
					client.Write(resObj);
				}
				break;
			case model.C_REQ_S_JOINROOM:
				c2s_JoinRoom(client,jsonByte);
				break;
			case model.C_REQ_S_LEAVEROOM:break;
			case model.C_REQ_S_SENDCHAT:break;
			case model.C_REQ_S_SENDADMINCMD:break;
			case model.C_REQ_S_UPLOADANSWERCMD:break;
			case model.C_REQ_S_USERLESSONRESULT:break;
			case model.C_REQ_S_UPLOADDATA:break;
			default:
				break;
			}
		}
	}else{
		fmt.Println("sid:",client.SID," 数据包解析错误:",err.Error());
	}
	
}

/**
处理进入教室
*/
func c2s_JoinRoom(client *LuBoClientConnection,jsonByte []byte){
	var req model.JoinRoom_c2s;
	err := json.Unmarshal(jsonByte,&req);
	if err == nil{
		preClient := OwnedConnectUIDMap_GetValue(req.Uid);//根据UID获取当前用户已经进入教室的socket连接，正常情况下应为nil
		if preClient == nil{
			//进入教室
			if joinRoom(client,req) == true{
				//全新的用户进入教室
				NewUserClientJoinRoom(client.SID,req.Uid,client);
			}
		}else{
			// if preClient == client{
			// 	//同一个socket，已经进入过教室，又重复的进教室
			// 	var rid = dataObj.rid || -1;
			// 	var roominfo = roomMap[rid];
			// 	if(roominfo)
			// 	{
			// 		//先调用离开教室
			// 		leaveRoom(sid,roominfo,uid);
			// 		//后调用进入教室
			// 		joinroom(sid,dataObj);
			// 	}
			// 	else{
			// 		joinroom(sid,dataObj);
			// 	}
			// }else{
			// 	//不同的socket，之前的socket已经存在于教室，则将其踢出
			// 	closePreUserSocket
			// }
		}
		// req.Rid;
		// req.TeachScript;
		// req.StartTimeinterval;
		// req.Uid;
		// req.NickName;
	}
}

/**进入教室*/
func joinRoom(client *LuBoClientConnection,req model.JoinRoom_c2s)bool{
	result := false;
	tempRes := CreateProtocal(model.S_RES_C_JOINROOM);
	if tempRes == nil{
		return result;
	}

	res,ok := tempRes.(*model.JoinRoom_s2c);
	if ok{
		if req.Uid <= 0{
			res.Code = 262;
			res.FaildMsg = "进入room失败,uid无效";
		}else if req.Rid < 0{
			res.Code = 263;
			res.FaildMsg = "进入room失败,roomId小于0,无效";
		}else{
			roomInfo := RoomInfoMap_GetValue(req.Rid);
			if roomInfo == nil{
				//如果教室信息不存在，则创建新的教室信息
				roomInfo = &model.RoomInfo{};
				roomInfo.Rid = req.Rid;
				roomInfo.RoomState = model.RoomState_NotStart;
				roomInfo.TeachingTmaterialScriptID = req.TeachScript;
				roomInfo.CurrentQuestionId = -1;
				roomInfo.AllowNewScript = true;
				RoomInfoMap_SetValue(req.Rid,roomInfo);//存入教室信息集合
			}

			lessonResultKey := fmt.Sprintf("%d_%d",req.Uid,req.Rid);
			if lessonResult := LessonResultMap_GetValue(lessonResultKey);lessonResult == nil{
				//初始化用户相关的课程报告集合
				LessonResultMap_SetValue(lessonResultKey,[]model.Answer_c2s{});
			}

			//将各种ID绑定到当前的socket上
			client.UID = req.Uid;
			client.RID = req.Rid;

			//加入到教室用户的列表
			user := model.CurrentUser{Uid:req.Uid,NickName:req.NickName,Type:true};
			roomInfo.UserArr = append(roomInfo.UserArr,user);
			roomInfo.UserIdArr = append(roomInfo.UserIdArr,user.Uid);
			roomInfo.AnswerUIDQueue = append(roomInfo.AnswerUIDQueue,user.Uid);

			//封装请求端的回执信息,并发送
			res.Code = 0;
			res.FaildMsg = "";
			res.Rid = req.Rid;
			res.UserArr = roomInfo.UserArr;
			client.Write(res);
			result = true;
			//向教室内的其它用户发送 用户状态变更通知
			//向该用户推送教室内缓存的文本消息通知
			//向该用户推送管理员操作命令通知

			//如果教材脚本加载完毕，则下推教材脚本
			if tsObj := TeachScriptMap_GetValue(req.TeachScript);tsObj != nil{
				pushTeachingTmaterialScriptLoadEndNotify(client,tsObj);
				//向该用户推送正在执行的教学命令
				sendTeachScriptNotify([]int64{req.Uid},req.Rid,roomInfo.TongyongCMDArr,roomInfo.CurrentTimeInterval,roomInfo.AnswerUIDQueue)
			}else{
				//加载教学脚本
			}
		}
	}
	
	return result;
}


/**
下推教材脚本中的资源相关配置
*/
func pushTeachingTmaterialScriptLoadEndNotify(client *LuBoClientConnection,tsObj map[string]interface{}){
	tempResult := CreateProtocal(model.S_NOTIFY_C_TEACHSCRIPTLOADEND);
	if tempResult == nil{
		return;
	}

	tsRes,ok:= tempResult.(*model.TeachScriptLoadEnd_s2c_notify);
	if ok{
		tsRes.Code = 0;
		tsRes.FaildMsg = "";
		courseID := uint32(0);
		if courseIDObj := tsObj["courseId"];courseIDObj != nil{
			courseID = courseIDObj.(uint32);
		}
		resource := map[string]interface{}{};
		if resourceObj := tsObj["resource"];resourceObj != nil{
			resource = resourceObj.(map[string]interface{});
		}
		tsRes.ScriptConfigData = model.ScriptConfigDataMap{CourseId:courseID,Resource:resource};
		//下推教材脚本的资源相关配置
		client.Write(tsRes);
	}
}

/**
 * 下发教学脚本
 * */
func sendTeachScriptNotify(uidArr []int64,rid int64,tongyongCMDArr []map[string]interface{},playTimeInterval int64,answerUIDQueue []int64){
	tempRes := CreateProtocal(model.S_NOTIFY_C_TEACHSCRIPTCMD);
	if tempRes == nil{
		return;
	}
	res,ok := tempRes.(*model.PushTeachScriptCache_s2c_notify);
	if ok{
		res.Code = 0;
		res.Rid = rid;
		res.Datas = tongyongCMDArr;
		res.AnswerUIDQueue = answerUIDQueue;
		res.PlayTimeInterval = playTimeInterval;
		for _,v:= range uidArr{
			sock := OwnedConnectUIDMap_GetValue(v)
			if sock != nil{
				sock.Write(res);
			}
		}
	}
}
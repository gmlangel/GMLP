package luboService;
import(
	"fmt"
	model "../models"
	"encoding/json"
	"net/http"
	"time"
	"io/ioutil"
	"log"
)

const (
	TeachScriptUrlFormat = "https://www.juliaol.cn/%d.cof?timeInterval=%d"
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
			case model.C_REQ_S_LEAVEROOM:
				c2s_LeaveRoom(client,jsonByte);
				break;
			case model.C_REQ_S_SENDCHAT:break;
			case model.C_REQ_S_SENDADMINCMD:break;
			case model.C_REQ_S_UPLOADANSWERCMD:
				c2s_UploadAnswerCMD(client,jsonByte);
				break;
			case model.C_REQ_S_USERLESSONRESULT:
				c2s_GetLessonResult(client,jsonByte);
				break;
			case model.C_REQ_S_UPLOADREPORTDATA:
				c2s_UploadReportData(client,jsonByte);
				break;
			default:
				break;
			}
		}
	}else{
		fmt.Println("sid:",client.SID," 数据包解析错误:",err.Error());
	}
	
}

func c2s_GetLessonResult(client *LuBoClientConnection,jsonByte []byte){
	var req model.UserLessonResult_c2s;
	err := json.Unmarshal(jsonByte,&req);
	var res *model.UserLessonResult_s2c = nil;
	tmp := CreateProtocal(model.S_RES_C_UPLOADREPORTDATA)
	if tmp != nil{
		if t,ok := tmp.(*model.UserLessonResult_s2c);ok == true{
			res = t;
			res.Code = 1
			res.FaildMsg = "未找到课程报告"
		}
	}
	if err == nil{
		rid := req.Rid;
		uid := req.Uid;
		resultKey := fmt.Sprintf("%d_%d",uid,rid);
		if lessArr := LessonResultMap_GetValue(resultKey);lessArr != nil && len(lessArr) > 0{
			if res != nil{
				res.Rid = rid;
				res.Datas = lessArr;
				res.Code = 0
				res.FaildMsg = ""
			}
		}
	}else{
		res.Code = 2
		res.FaildMsg = "课程报告的请求数据的格式有问题，Unmarshal失败"
	}

	if res != nil{
		client.Write(res);
	}
}

/*
处理 用户上报的数据
*/
func c2s_UploadReportData(client *LuBoClientConnection,jsonByte []byte){
	var req model.DataReport_c2s;
	err := json.Unmarshal(jsonByte,&req);
	if err == nil{
		tmp := CreateProtocal(model.S_RES_C_UPLOADREPORTDATA);
		if tmp != nil{
			if res,ok := tmp.(*model.DataReport_s2c);ok == true{
				res.Rid = client.RID;
				client.Write(res)
			}
		}
	}
}

/**
处理用户上报的“做题答案”
*/
func c2s_UploadAnswerCMD(client *LuBoClientConnection,jsonByte []byte){
	var req model.UploadAnswer_c2s;
	err := json.Unmarshal(jsonByte,&req);
	if err == nil{
		roomInfo := RoomInfoMap_GetValue(client.RID);
		if roomInfo != nil{
			if roomInfo.CurrentQuestionId != req.Id{
				return;//如果学生上报的答案，不是当前的问题的答案，则不作数
			}
			tempUid := req.Uid;
			tempArr := roomInfo.WaitAnswerUids
			for i,v := range tempArr{
				if v == tempUid{
					//从等待答题的用户列表中移除改用户
					roomInfo.WaitAnswerUids = append(roomInfo.WaitAnswerUids[0:i],roomInfo.WaitAnswerUids[i+1:]...);
					//每一个用户提交答案后进行判断，脚本执行时间不足5秒的，补充至5秒
					if roomInfo.CompleteTime - roomInfo.CurrentTimeInterval < 5{
						roomInfo.CompleteTime = roomInfo.CurrentTimeInterval + 5
					}
					//记录用户相关的课程报告
					resultKey := fmt.Sprintf("%d_%d",client.UID,client.RID);
					resultArr := LessonResultMap_GetValue(resultKey)
					ans_c2s := model.Answer_c2s{Id:req.Id,Data:req.Data};
					if resultArr == nil{
						resultArr = append([]model.Answer_c2s{},ans_c2s)
					}else{
						resultArr = append(resultArr,ans_c2s);
					}
					LessonResultMap_SetValue(resultKey,resultArr);
					break;
				}
			}
			//通过判断是否所有的用户都已经答题完毕，5秒后更新allowNewScript（“是否下发下一个教学脚本”）的状态，  5秒的时间是留给客户端播放奖励声音和动画
			if len(roomInfo.WaitAnswerUids) == 0{
				roomInfo.CompleteTime = roomInfo.CurrentTimeInterval + 5
			}

			//发送客户端回执
			res := &model.UploadAnswer_s2c{Cmd:model.S_RES_C_UPLOADANSWERCMD,Code:0,FaildMsg:""};
			client.Write(res);
		}
	}
}

/**
处理离开教室
*/
func c2s_LeaveRoom(client *LuBoClientConnection,jsonByte []byte){
	var req model.LeaveRoom_c2s;
	err := json.Unmarshal(jsonByte,&req);
	if err == nil{
		roomInfo := RoomInfoMap_GetValue(req.Rid);
		if roomInfo != nil{
			leaveRoom(client,req.Uid,roomInfo)
		}
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
			joinRoom(client,req);
		}else{
			if preClient == client{
				//同一个socket，已经进入过教室，又重复的进教室
				roominfo := RoomInfoMap_GetValue(req.Rid);
				if roominfo != nil{
					//先调用离开教室
					leaveRoom(client,req.Uid,roominfo);
					//后调用进入教室
					joinRoom(client,req);
				}else{
					joinRoom(client,req);
				}
			}else{
				//不同的socket，之前的socket已经存在于教室，则将其踢出
				DestroySocket(client,func(){
					joinRoom(client,req);
				})
			}
		}
	}
}

/**进入教室*/
func joinRoom(client *LuBoClientConnection,req model.JoinRoom_c2s){
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
				roomInfo.TeachingTmaterialScriptID = req.TeachScriptID;
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
			if tsObj := TeachScriptMap_GetValue(req.TeachScriptID);tsObj != nil{
				pushTeachingTmaterialScriptLoadEndNotify(client,tsObj);
				//向该用户推送正在执行的教学命令
				sendTeachScriptNotify([]int64{req.Uid},req.Rid,roomInfo.TongyongCMDArr,roomInfo.CurrentTimeInterval,roomInfo.AnswerUIDQueue)
			}else{
				//加载教学脚本
				go loadTeachingTmaterialScript(client,req,roomInfo);
			}
		}
	}
	if result == true{
		//全新的用户进入教室
		NewUserClientJoinRoom(client.SID,req.Uid,client);
	}
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
		isOk := false;
		if resourceObj := tsObj["resource"];resourceObj != nil{
			resource,ok:= resourceObj.(map[string]interface{});
			if ok == true{
				tsRes.ScriptConfigData = model.ScriptConfigDataMap{CourseId:courseID,Resource:resource};
				isOk = true;
			}
		}
		if isOk == false{
			tsRes.Code = 1;
			tsRes.FaildMsg = "数据格式转换失败"
		}
		
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

/**
加载教学脚本
*/
func loadTeachingTmaterialScript(client *LuBoClientConnection,req model.JoinRoom_c2s,roomInfo *model.RoomInfo){
	teachScriptID := req.TeachScriptID;
	url := fmt.Sprintf(TeachScriptUrlFormat,teachScriptID,time.Now().Unix());
	resp,err :=  http.Get(url)//请求 教材脚本资源
    if err != nil {
		// handle error
		log.Println("教材ID:",teachScriptID," 请求教材脚本资源出错:",err.Error())
		return;
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		// handle error
		log.Println("教材ID:",teachScriptID," 教材脚本资源内容有问题:",err.Error())
		return;
	}
	
	var teachScriptObj map[string]interface{};
	err = json.Unmarshal(body,&teachScriptObj);
	if err != nil{
		log.Println("教材ID:",teachScriptID," 教材脚本资源转JSON出错:",err.Error());
		return;
	}
	fmt.Println("教材ID:",teachScriptID," 加载完毕");
	TeachScriptMap_SetValue(teachScriptID,teachScriptObj);//将教学脚本添加至脚本集合
	pushTeachingTmaterialScriptLoadEndNotify(client,teachScriptObj);
	//向该用户推送正在执行的教学命令
	sendTeachScriptNotify([]int64{req.Uid},req.Rid,roomInfo.TongyongCMDArr,roomInfo.CurrentTimeInterval,roomInfo.AnswerUIDQueue)
}

/**
离开教室
*/
func leaveRoom(client *LuBoClientConnection,uid int64,roomInfo *model.RoomInfo){
	UnOwnedConnect_SetValue(client.SID,client);
	OwnedConnect_SetValue(client.SID,nil);
	OwnedConnectUIDMap_SetValue(uid,nil);
	temp := CreateProtocal(model.S_RES_C_LEAVEROOM);
	if temp != nil{
		if res,ok := temp.(*model.LeaveRoom_s2c);ok == true{
			res.Rid = client.RID;
			res.Uid = client.UID;
			client.Write(res);
		}
	}
	client.RID = -1;
	client.UID = -1;
}
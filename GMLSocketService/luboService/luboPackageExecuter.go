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
		//取cmd，并决策执行
		cmd := jsonObj["cmd"];
		if temp,ok := cmd.(float64);ok == true{
			command := uint32(temp);
			fmt.Println("sid:",client.SID," 收到数据包的cmd:",command);
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
		log.Println("sid:",client.SID," 数据包解析错误:",err.Error());
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
		uid := req.Uid;
		preClient := OwnedConnectUIDMap_GetValue(uid);//根据UID获取当前用户已经进入教室的socket连接，正常情况下应为nil
		if preClient == nil{
			//进入教室
			joinRoom(client,req);
		}else{
			if preClient == client{
				//同一个socket，已经进入过教室，又重复的进教室
				roominfo := RoomInfoMap_GetValue(req.Rid);
				if roominfo != nil{
					//先调用离开教室
					leaveRoom(client,uid,roominfo);
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
				roomInfo.TongyongCMDArr = []map[string]interface{}{};
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
			client.RoomInfo = roomInfo;

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
				sendTeachScriptToUser(client,req.Rid,roomInfo.TongyongCMDArr,roomInfo.CurrentTimeInterval,roomInfo.AnswerUIDQueue)
				client.TeachScriptStepDataArr = getObjArray(tsObj["stepData"],nil);
				go loopSendTeachScript(client);//定时下发教学脚本
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
		tempCouseID := getInt64(tsObj["courseId"],-1)
		if tempCouseID > -1{
			courseID := uint32(tempCouseID);
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
		}else{
			tsRes.Code = 1;
			tsRes.FaildMsg = "数据格式转换失败"
		}
	
		//下推教材脚本的资源相关配置
		client.Write(tsRes);
	}
}

// /**
//  * 下发教学脚本
//  * */
// func sendTeachScriptNotify(uidArr []int64,rid int64,tongyongCMDArr []map[string]interface{},playTimeInterval int64,answerUIDQueue []int64){
// 	tempRes := CreateProtocal(model.S_NOTIFY_C_TEACHSCRIPTCMD);
// 	if tempRes == nil{
// 		return;
// 	}
// 	res,ok := tempRes.(*model.PushTeachScriptCache_s2c_notify);
// 	if ok{
// 		res.Code = 0;
// 		res.Rid = rid;
// 		res.Datas = tongyongCMDArr;
// 		res.AnswerUIDQueue = answerUIDQueue;
// 		res.PlayTimeInterval = playTimeInterval;
// 		for _,v:= range uidArr{
// 			sock := OwnedConnectUIDMap_GetValue(v)
// 			if sock != nil{
// 				sock.Write(res);
// 			}
// 		}
// 	}
// }


/**
 * 下发教学脚本
 * */
 func sendTeachScriptToUser(sock *LuBoClientConnection,rid int64,tongyongCMDArr []map[string]interface{},playTimeInterval int64,answerUIDQueue []int64){
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
		sock.Write(res);
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
	sendTeachScriptToUser(client,req.Rid,roomInfo.TongyongCMDArr,roomInfo.CurrentTimeInterval,roomInfo.AnswerUIDQueue)

	client.TeachScriptStepDataArr = getObjArray(teachScriptObj["stepData"],nil);
	go loopSendTeachScript(client);//定时下发教学脚本
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
	ClearUIDByRoomInfo(roomInfo,uid);
}

//清空各种数组中现存的UID
func ClearUIDByRoomInfo(roomInfo *model.RoomInfo,uid int64){
	if nil == roomInfo{
		return;
	}
	
	uidArr := roomInfo.UserIdArr;
	userArr := roomInfo.UserArr;
	ansArr := roomInfo.AnswerUIDQueue;
	waitArr := roomInfo.WaitAnswerUids;
	for i,v := range uidArr{
		if v == uid{
			roomInfo.UserIdArr = append(roomInfo.UserIdArr[0:i],roomInfo.UserIdArr[i+1:]...);
			break;
		}
	}

	for i,v := range userArr{
		if v.Uid == uid{
			roomInfo.UserArr = append(roomInfo.UserArr[0:i],roomInfo.UserArr[i+1:]...);
			break;
		}
	}

	for i,v := range ansArr{
		if v == uid{
			roomInfo.AnswerUIDQueue = append(roomInfo.AnswerUIDQueue[0:i],roomInfo.AnswerUIDQueue[i+1:]...);
			break;
		}
	}

	for i,v := range waitArr{
		if v == uid{
			roomInfo.WaitAnswerUids = append(roomInfo.WaitAnswerUids[0:i],roomInfo.WaitAnswerUids[i+1:]...);
			break;
		}
	}
}

//定时下发教学脚本
func loopSendTeachScript(client *LuBoClientConnection){
	client.GTimerInterval = time.Now().Unix();//获取当前服务器时间
	var curTime int64 = 0;
	for client.SID != -1{
		time.Sleep(model.TeachScriptTimeInterval);//每隔一定时间，计算要下发的教材脚本
		rid := client.RID;
		roomInfo := client.RoomInfo;
		stepDataArr := client.TeachScriptStepDataArr;
		if nil == roomInfo || nil == stepDataArr{
			break;
		}
		if roomInfo.RoomState == model.RoomState_End{
			break;//如果课程已经停止，则停止下发数据
		}
		curTime = time.Now().Unix();
		if roomInfo.RoomState == model.RoomState_Started{
			//课程已开始， 计时下发指定教材
			roomInfo.CurrentTimeInterval += curTime - client.GTimerInterval;
			client.GTimerInterval = curTime;//更新上一次处理脚本时的时间记录.
			if roomInfo.CurrentTimeInterval >= roomInfo.CompleteTime{
				//已经达到超时时间，为了不影响之后的脚本运行，则应该直接执行下个脚本
				roomInfo.CurrentTimeInterval = 0;
				roomInfo.AllowNewScript = true;
			}

			if roomInfo.AllowNewScript == false{
				continue;
			}
			cmdArr := []map[string]interface{}{};//要下发的教学脚本数组
			j := int64(len(stepDataArr));
			for roomInfo.CurrentStepIdx < j{
				needBread := false;//是否需要跳出循环
				scriptItem := stepDataArr[roomInfo.CurrentStepIdx];//获取一条教学命令
				roomInfo.CurrentStepIdx += 1;//更新教学命令的 索引游标
				roomInfo.CurrentQuestionId = getInt64(scriptItem["id"],-1);//设置当前正在提问的问题ID
				//将服务端脚本转换为客户端可以执行的脚本命令
				clientScriptItem := map[string]interface{}{"suid":0,"st":curTime,"data":scriptItem};
				cmdArr = append(cmdArr,clientScriptItem);//将脚本塞入 下发列表
				sType := getString(scriptItem["type"],"");
				switch sType{
					case "changePage":
						//移除之前的批处理教学命令缓存,添加新的教学命令缓存
						roomInfo.TongyongCMDArr = []map[string]interface{}{clientScriptItem};
						break;
					case "onWall":
						//添加新的教学命令缓存
						roomInfo.TongyongCMDArr = append(roomInfo.TongyongCMDArr,clientScriptItem);
						break;
					case "delay":
						//延迟一定时间后，下发下一条命令
						roomInfo.CompleteTime = getInt64(getMap(scriptItem["value"],map[string]interface{}{})["timeLength"],0);
						roomInfo.AllowNewScript = false;
						cmdArr = cmdArr[0:len(cmdArr)-1];//从下发命令集合中删除delay命令
						needBread = true;
						break;
					case "classEnd":
						roomInfo.TongyongCMDArr = roomInfo.TongyongCMDArr[0:1];//除第一条换页命令外，移除其余的命令
						roomInfo.TongyongCMDArr = append(roomInfo.TongyongCMDArr,clientScriptItem);//添加新的教学命令缓存
						roomInfo.RoomState = model.RoomState_End;//更新教室状态
						//测试用， 重置教室，反复使用教室
						RoomInfoMap_SetValue(rid,nil);
						break;
					case "templateCMD":
						roomInfo.WaitAnswerUids = roomInfo.UserIdArr;//设置应答序列
						//设置超时等待时间和等待回答响应的用户数组
						roomInfo.CompleteTime = getInt64(getMap(scriptItem["value"],map[string]interface{}{})["timeout"],30);
						roomInfo.TongyongCMDArr = append(roomInfo.TongyongCMDArr,clientScriptItem);//添加新的教学d命令到缓存
						roomInfo.AllowNewScript = false;
						needBread = true;
						break;
					case "audio":
						roomInfo.WaitAnswerUids = roomInfo.UserIdArr;//设置应答序列
						//设置应答超时时间
						item := getMap(scriptItem["value"],nil);
						if nil != item{
							roomInfo.CompleteTime = getInt64(item["endSecond"],1) - getInt64(item["beginSecond"],1) + 3;
						}else{
							roomInfo.CompleteTime = 5;
						}
						roomInfo.TongyongCMDArr = append(roomInfo.TongyongCMDArr,clientScriptItem);//添加新的教学d命令到缓存
						roomInfo.AllowNewScript = false;
						needBread = true;
						break;
					case "video":
						roomInfo.WaitAnswerUids = roomInfo.UserIdArr;//设置应答序列
						//设置应答超时时间
						item := getMap(scriptItem["value"],nil);
						if nil != item{
							roomInfo.CompleteTime = getInt64(item["endSecond"],1) - getInt64(item["beginSecond"],1) + 3;
						}else{
							roomInfo.CompleteTime = 5;
						}
						roomInfo.TongyongCMDArr = append(roomInfo.TongyongCMDArr,clientScriptItem);//添加新的教学d命令到缓存
						roomInfo.AllowNewScript = false;
						needBread = true;
						break;
					default:break;
				} 

				if needBread == true{
					break;//跳出循环
				}
			}

			if len(cmdArr) > 0{
				//下发教学命令到客户端
				sendTeachScriptToUser(client,rid,cmdArr,roomInfo.CurrentTimeInterval,roomInfo.AnswerUIDQueue);
				cmdArr = []map[string]interface{}{};//清空已发的命令集合
			}
		}else{
			//计时，实时更新课程状态
			if curTime > roomInfo.StartTimeInterval{
				roomInfo.RoomState = model.RoomState_Started
			}
		}
		client.GTimerInterval = curTime;//更新上一次处理脚本时的时间记录.
	}
}

/*object转字符传*/
func getString(val interface{},def string)string{
	if nil == val{
		return def
	}
	result,ok := val.(string);
	if ok{
		return result;
	}else{
		return def;
	}
}

/*object转int64*/
func getInt64(val interface{},def int64)int64{
	if nil == val{
		return def;
	}
	result,ok := val.(float64);
	if ok{
		return int64(result);
	}else{
		return def;
	}
}

/*object转bool*/
func getBool(val interface{},def bool)bool{
	if nil == val{
		return def;
	}
	result,ok := val.(bool);
	if ok{
		return result;
	}else{
		return def;
	}
}

/*object转map*/
func getMap(val interface{},def map[string]interface{})map[string]interface{}{
	if nil == val{
		return def;
	}
	result,ok := val.(map[string]interface{});
	if ok{
		return result;
	}else{
		return def;
	}
}

/*object转[]map[string]interface{}*/
func getObjArray(val interface{},def []map[string]interface{})[]map[string]interface{}{
	if nil == val{
		return def;
	}
	tem,ok := val.([]interface{});
	if ok{
		result := []map[string]interface{}{};
		for _,v := range tem{
			resultV := getMap(v,nil);
			if nil != resultV{
				result = append(result,resultV)
			}
		}
		return result;
	}else{
		return def;
	}
}
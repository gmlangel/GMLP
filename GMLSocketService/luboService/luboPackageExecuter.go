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
	TeachScriptUrlFormat = "https://www.juliaol.cn/%d.json?timeInterval=%d"
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
			case model.C_REQ_S_Test:
				c2s_Test(client,jsonByte);
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
快捷测试接口
*/
func c2s_Test(client *LuBoClientConnection,jsonByte []byte){
	var req model.ToTest_c2s;
	err := json.Unmarshal(jsonByte,&req);
	if err == nil{
		//拼接要下发的消息
		clientScriptItem := map[string]interface{}{"suid":0,"st":0,"playInterval":0,"data":&req.Msg};
		cmdArr := append([]map[string]interface{}{},clientScriptItem);//将脚本塞入 下发列表
		beControlClient := OwnedConnectUIDMap_GetValue(req.Uid);
		if nil != beControlClient{
			//停用客户端socket原有的教室loop教学脚本下发
			beControlClient.MediaStepDataArr = nil;
			//向指定客户端socket连接下发消息
			sendTeachScriptToUser(beControlClient,req.Rid,cmdArr,[]int64{});
		}
	}
}
	
/**判断当前的教学环节是否已经结束*/
func answerIsEnd(idx int,stepArr []model.ScriptStepData,roomInfo *model.RoomInfo)bool{
	if idx > -1 && idx < len(stepArr){
		stepData := stepArr[idx];
		if "templateCMD" == stepData.Type{
			//如果step != nil,且step的type == “templateCMD”，视为答题环节未结束，则return false
			return false;
		}else{
			if "star" == stepData.Type{
				roomInfo.WantAddCredit += int(getInt64(stepData.Value["count"],0));//记录即将获得的星星数量，用于断线重连
			}
			return answerIsEnd(int(stepData.Next),stepArr,roomInfo)
		}
	}else{
		//如果step == nil,视为答题完毕 则return true
		return true;
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
			if roomInfo.SCurrentQuestionId != req.Id{
				return;//如果学生上报的答案，不是当前的问题的答案，则不作数
			}

			_,dataSyncOk := <- client.dataSyncChan;
			if false == dataSyncOk{
				return;
			}
			roomInfo.SCurrentQuestionId = -1;//匹配成功后重置题号，比面重复匹配
			tempUid := req.Uid;
			j := len(roomInfo.SWaitAnswerUids);
			needRecordAnswerReport := false;//是否要记录用户上传上来的答案，用于生成报告
			for i:=0;i<j;{
				if roomInfo.SWaitAnswerUids[i] == tempUid{
					//从等待答题的用户列表中移除该用户
					roomInfo.SWaitAnswerUids = append(roomInfo.SWaitAnswerUids[0:i],roomInfo.SWaitAnswerUids[i+1:]...);
					needRecordAnswerReport = true;
					j -= 1;
				}else{
					i += 1;
				}
			}

			if true == needRecordAnswerReport{
				// //每一个用户提交答案后进行判断，脚本执行时间不足3秒的，补充至3秒
				// if roomInfo.SCurrentQuesionTimeOut - roomInfo.SCurrentTimeInterval < 3{
				// 	roomInfo.SCurrentQuesionTimeOut = roomInfo.SCurrentTimeInterval + 3//更新关键帧脚本的时间
				// }
				//记录用户相关的课程报告
				resultKey := fmt.Sprintf("%d_%d",client.UID,client.RID);
				resultArr := LessonResultMap_GetValue(resultKey)
				ans_c2s := model.Answer_c2s{Id:req.Id,Data:req.Data};
				score := req.Data.Score;
				if  score >= 60{
					roomInfo.CurrentAnswerState = "success";
				}else if score >= 0 && score < 60{
					roomInfo.CurrentAnswerState = "faild";
				}else{
					roomInfo.CurrentAnswerState = "timeouterr";
				}
				//计算当前问题是否已经答题完毕，从而确定教学环节是否已经完毕，避免重新进入教室后，播放重复的教学环节
				nextState := -1;
				if nil != roomInfo.SCurrent{
					nextState = int(getInt64(roomInfo.SCurrent.Value[roomInfo.CurrentAnswerState],-1));
				}
				if answerIsEnd(nextState,client.TeachScriptStepDataArr,roomInfo) == true{
					roomInfo.NeedPlayBack = false;//当答题错误后，并且当前问题的Next 链中没有再次答题的机会时，认定该问题已近答过，设置NeedPlayBack值为false，避免用户重新进入教室后，重复答题，
				}else{
					roomInfo.NeedPlayBack = true;
				}
				//将答题结果记录至课程报告数组
				if resultArr == nil{
					resultArr = append([]model.Answer_c2s{},ans_c2s)
				}else{
					resultArr = append(resultArr,ans_c2s);
				}
				LessonResultMap_SetValue(resultKey,resultArr);
			}

			//通过判断是否所有的用户都已经答题完毕，更新allowNewScript,一遍处理接下来的决策树
			if len(roomInfo.SWaitAnswerUids) == 0{
				roomInfo.SAllowNew = true;
			}
			client.dataSyncChan <- 1;

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
				DestroySocket(preClient,func(){
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
		nowTime := time.Now().Unix();
		if req.Uid <= 0{
			res.Code = 262;
			res.FaildMsg = "进入room失败,uid无效";
		}else if req.Rid < 0{
			res.Code = 263;
			res.FaildMsg = "进入room失败,roomId小于0,无效";
		}else if req.StartTimeinterval > nowTime{
			res.Code = 14;
			res.FaildMsg = "进入room失败,课程还未开始";
		}else{
			roomInfo := RoomInfoMap_GetValue(req.Rid);
			if roomInfo == nil{
				//如果教室信息不存在，则创建新的教室信息
				roomInfo = &model.RoomInfo{};
				roomInfo.Rid = req.Rid;
				roomInfo.RoomState = model.RoomState_NotStart;
				roomInfo.TeachingTmaterialScriptID = req.TeachScriptID;
				roomInfo.NeedPlayPreMedia = false;
				roomInfo.NeedPlayBack = false;
				// roomInfo.SCurrentQuestionId = -1;
				// roomInfo.MAllowNew = true;
				// roomInfo.SAllowNew = true;
				// roomInfo.SCurrent = nil;
				// roomInfo.MCurrent = nil;
				// roomInfo.SCurrentTimeInterval = 0;
				//roomInfo.CurrentProcess = 0;
				roomInfo.NeedJumpSleep = false;
				roomInfo.CurrentAnswerState = "";
				roomInfo.CreateTime = time.Now().Unix();
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
			ClearUIDByRoomInfo(roomInfo,user.Uid);//清楚已经存在的相同的UID数据
			//添加新的UID到各种数组中
			roomInfo.UserArr = append(roomInfo.UserArr,user);
			roomInfo.UserIdArr = append(roomInfo.UserIdArr,user.Uid);
			roomInfo.AnswerUIDQueue = append(roomInfo.AnswerUIDQueue,user.Uid);

			//封装请求端的回执信息,并发送
			res.Code = 0;
			res.FaildMsg = "";
			res.Rid = req.Rid;
			res.UserArr = roomInfo.UserArr;
			res.Credit = roomInfo.Credit + roomInfo.WantAddCredit;
			roomInfo.Credit = res.Credit;
			client.Write(res);
			result = true;
			//向教室内的其它用户发送 用户状态变更通知
			//向该用户推送教室内缓存的文本消息通知
			//向该用户推送管理员操作命令通知

			//每次重新进入教室后，重置教室内各脚本的执行进度，使每个脚本都从第0秒重新播放
			roomInfo.SCurrentTimeInterval = 0;
			roomInfo.MAllowNew = true;
			roomInfo.MainFrames = nil;
			roomInfo.MCurrentTimeInterval = 0;
			roomInfo.MCurrentMainFrameIdx = 0;
			roomInfo.MCompleteTime = 0;
			roomInfo.SAllowNew = true;
			roomInfo.SCurrent = nil;
			roomInfo.SCurrentTimeInterval = 0;
			roomInfo.SCurrentQuesionTimeOut = -1;
			roomInfo.CurrentProcess = 0;//重置播放状态。使其播放媒体命令
			roomInfo.NeedJumpSleep = false;
			roomInfo.SWaitAnswerUids = nil;
			roomInfo.SCurrentQuestionId = -1;
			if true == roomInfo.NeedPlayBack{
				//当当前的问题流程还未答完，则从上一次退出教室时的关键帧对应的视频节点开始播放
				roomInfo.MCurrentTimeInterval = roomInfo.MMainFramePrePlayInterval;
				roomInfo.MMainFramePrePlayInterval = 0;
				roomInfo.NeedPlayPreMedia = true;
			}else{
				//当当前的问题流程回答完毕，不论之后的奖励和视频是否播完，则播放时间置0，播下一个视频

				//这里有个问题，如果答题完毕，但是星星还未发，目前的逻辑会出现少星星的问题
				roomInfo.MCurrentTimeInterval = 0;
				roomInfo.NeedPlayPreMedia = false;
				roomInfo.MPreMainFrameIdx = 0;
				roomInfo.MMainFramePrePlayInterval = 0;
			}

			//计算课程开课时间和结束时间
			roomInfo.StartTimeinterval = req.StartTimeinterval;
			roomInfo.EndTimeinterval = req.EndTimeinterval + 1800;//延长课程结束时间30分钟，防止用户在正常上课期间频繁掉线，导致课程无法准点结束
			//如果教材脚本加载完毕，则下推教材脚本
			if tsObj := TeachScriptMap_GetValue(req.TeachScriptID);tsObj != nil{
				pushTeachingTmaterialScriptLoadEndNotify(client,tsObj);
				if nil != roomInfo.PreChangePageCMD{
					sendTeachScriptToUser(client,req.Rid,[]map[string]interface{}{roomInfo.PreChangePageCMD},roomInfo.AnswerUIDQueue);
				}
				client.TeachScriptStepDataArr = objArrToStepDataArr(getObjArray(tsObj["stepData"],nil));
				client.MediaStepDataArr = objArrToMediaDataArr(getObjArray(tsObj["mediaData"], nil))
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
		tempWidth := uint32(getInt64(tsObj["width"],0));
		tempHeigth := uint32(getInt64(tsObj["height"],0));
		if tempCouseID > -1{
			courseID := uint32(tempCouseID);
			isOk := false;
			if resourceObj := tsObj["resource"];resourceObj != nil{
				resource,ok:= resourceObj.(map[string]interface{});
				if ok == true{
					tsRes.ScriptConfigData = model.ScriptConfigDataMap{CourseId:courseID,Resource:resource,Width:tempWidth,Height:tempHeigth};
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


/**
 * 下发教学脚本
 * */
 func sendTeachScriptToUser(sock *LuBoClientConnection,rid int64,cmdArr []map[string]interface{},answerUIDQueue []int64){
	tempRes := CreateProtocal(model.S_NOTIFY_C_TEACHSCRIPTCMD);
	if tempRes == nil || len(cmdArr) == 0{
		return;
	}
	res,ok := tempRes.(*model.PushTeachScriptCache_s2c_notify);
	if ok{
		res.Code = 0;
		res.Rid = rid;
		res.Datas = cmdArr;
		res.AnswerUIDQueue = answerUIDQueue;
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
	if nil != roomInfo.PreChangePageCMD{
		sendTeachScriptToUser(client,req.Rid,[]map[string]interface{}{roomInfo.PreChangePageCMD},roomInfo.AnswerUIDQueue);
	}
	client.TeachScriptStepDataArr = objArrToStepDataArr(getObjArray(teachScriptObj["stepData"],nil));
	client.MediaStepDataArr = objArrToMediaDataArr(getObjArray(teachScriptObj["mediaData"],nil));
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
	waitArr := roomInfo.SWaitAnswerUids;
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
			roomInfo.SWaitAnswerUids = append(roomInfo.SWaitAnswerUids[0:i],roomInfo.SWaitAnswerUids[i+1:]...);
			break;
		}
	}
}

//定时下发教学脚本
func loopSendTeachScript(client *LuBoClientConnection){
	_,chanOK := <- client.runLoopExecChan;
	if false == chanOK{
		return;
	}
	client.GTimerInterval = time.Now().Unix();//获取当前服务器时间
	var curTime int64 = 0;
	var clientScriptItem map[string]interface{} = nil;
	for client.RID != -1 && client.SID != -1{
		rid := client.RID;
		roomInfo := client.RoomInfo;
		stepDataArr := client.TeachScriptStepDataArr;
		mediaDataArr := client.MediaStepDataArr;
		if nil == roomInfo || nil == stepDataArr || nil == mediaDataArr{
			break;
		}
		if true == roomInfo.NeedJumpSleep{
			//跳过延时，营造一种视频或者脚本连续播放的效果
			roomInfo.NeedJumpSleep = false;
		}else{
			time.Sleep(model.TeachScriptTimeInterval);//每隔一定时间，计算要下发的教材脚本
		}
		
		curTime = time.Now().Unix();
		if curTime > roomInfo.EndTimeinterval || roomInfo.RoomState == model.RoomState_End{
			//如果当前时间大于课程结束时间、或者课程处于结束状态，则直接下发classend命令，告知客户端，并break出循环
			roomInfo.RoomState = model.RoomState_End
			tempScript := map[string]interface{}{"id":len(mediaDataArr),"type":"classEnd","value":map[string]interface{}{}};
			clientScriptItem = map[string]interface{}{"suid":0,"playInterval":0,"st":curTime,"data":tempScript};
			cmdArr := []map[string]interface{}{clientScriptItem};//将脚本塞入 下发列表
			sendTeachScriptToUser(client,rid,cmdArr,roomInfo.AnswerUIDQueue);
			break;//如果课程已经停止，则停止下发数据
		}
		if roomInfo.RoomState == model.RoomState_Started{
			cmdArr := []map[string]interface{}{};//要下发的教学脚本数组
			offsetTime := curTime - client.GTimerInterval;
			//课中处理逻辑
			if 1 == roomInfo.CurrentProcess{
				//处理教学脚本
				//如果允许执行下一个教学脚本，则执行
				if true == roomInfo.SAllowNew{
					if nil != roomInfo.MainFrames{
						//尝试生成要执行的关键帧脚本数组
						tempCmdArr,frameStepIdx,_ := execStepDataByMainFrames(roomInfo.MainFrames,stepDataArr,roomInfo,roomInfo.MCurrentMainFrameIdx,curTime);
						roomInfo.MCurrentMainFrameIdx = frameStepIdx;//更新关键帧执行进度
						if len(tempCmdArr) > 0{
							cmdArr = append(cmdArr,tempCmdArr...);//将脚本塞入 下发列表
						}
						if len(cmdArr) > 0{
							//遍历到了关键正命令，则顺序执行关键帧
							//下发教学命令到客户端
							sendTeachScriptToUser(client,rid,cmdArr,roomInfo.AnswerUIDQueue);
							cmdArr = []map[string]interface{}{};//清空已发的命令集合
						}else{
							if int(roomInfo.MCurrentMainFrameIdx) >= len(roomInfo.MainFrames){
								//没有可执行的关键帧时，继续播放媒体脚本
								roomInfo.SCurrent = nil;
								roomInfo.CurrentProcess = 0;
								roomInfo.NeedJumpSleep = true;
							}
						}
					}else{
						//如果没有对应的关键帧序列，则继续播放媒体脚本
						roomInfo.SCurrent = nil;
						roomInfo.CurrentProcess = 0;
						roomInfo.NeedJumpSleep = true;
					}
				}
				//更新脚本计时时间
				roomInfo.MCurrentTimeInterval += offsetTime;//无论是否正在执行教学脚本， 都要更新媒体脚本的计时，防止教学脚本执行完毕后，要过很久才能迎来新的媒体脚本播放时机
				roomInfo.SCurrentTimeInterval += offsetTime;
				client.GTimerInterval = curTime;//更新上一次处理脚本时的时间记录.
				if roomInfo.SCurrentQuesionTimeOut > -1 && roomInfo.SCurrentTimeInterval >= roomInfo.SCurrentQuesionTimeOut{
					//关键帧脚本已经达到超时时间，为了不影响之后的脚本运行，则应该直接执行下个关键帧
					roomInfo.SCurrentQuesionTimeOut = -1;
					_,dataSyncOk := <- client.dataSyncChan;
					if false == dataSyncOk{
						break;
					}
					roomInfo.SAllowNew = true;
					roomInfo.CurrentAnswerState = "timeouterr";//设置答题结果为'超时'
					//重置题号，避免触发了超时响应后，又进入答案上报计算逻辑，
					roomInfo.SCurrentQuestionId = -1;
					if nil != roomInfo.SCurrent && "templateCMD" == roomInfo.SCurrent.Type{
						nextState := -1;
						nextState = int(getInt64(roomInfo.SCurrent.Value[roomInfo.CurrentAnswerState],-1));
						if answerIsEnd(nextState,client.TeachScriptStepDataArr,roomInfo) == true{
							roomInfo.NeedPlayBack = false;//当答题超时后，并且当前问题的Next 链中没有再次答题的机会时，认定该问题已近答过，设置NeedPlayBack值为false，避免用户重新进入教室后，重复答题，
						}else{
							roomInfo.NeedPlayBack = true;
						}
					}
					client.dataSyncChan <- 1;
				}
			}else{
				//处理媒体脚本
				//遍历媒体播放数组
				if true == roomInfo.MAllowNew{
					idx := int64(0);
					j := int64(len(mediaDataArr));
					if nil != roomInfo.MCurrent{
						if true == roomInfo.NeedPlayPreMedia{
							roomInfo.NeedPlayPreMedia = false;
							//如果想要播放退出教室前的上一次的媒体视频，则走这个流程
							idx = roomInfo.MCurrent.Id;
						}else{
							//播放MCurrent的下一个视频
							idx = roomInfo.MCurrent.Next;//如果存在当前播放列，则取当前播放项的下一条进行播放
							roomInfo.MPreMainFrameIdx = 0;
							roomInfo.MMainFramePrePlayInterval = 0;
						}
					}
					
					if idx > -1 && idx < j{
						mediaItem := &mediaDataArr[idx];//获取一条教学命令
						roomInfo.NeedPlayBack = true;//每次成功获取到一条媒体播放命令后，都要设置NeedPlayBack为true，这样可以让用户每次进入教室后，播放上一次未完成的视频
						roomInfo.MCurrentMainFrameIdx = roomInfo.MPreMainFrameIdx;
						roomInfo.MPreMainFrameIdx = 0;
						roomInfo.MCurrent = mediaItem;
						//将服务端脚本转换为客户端可以执行的脚本命令
						clientScriptItem = map[string]interface{}{"suid":0,"st":curTime,"playInterval":roomInfo.MCurrentTimeInterval,"data":mediaConverScript(mediaItem)};
						cmdArr = append(cmdArr,clientScriptItem);//将脚本塞入 下发列表
						itemValue := mediaItem.Value;
						sType := mediaItem.Type;
						switch sType{
							case "video":
								roomInfo.MainFrames = mediaItem.MainFrames;
								roomInfo.MAllowNew = false;
								roomInfo.MCompleteTime = getInt64(itemValue["endSecond"],0) - getInt64(itemValue["beginSecond"],0) + 2;//设置脚本超时时间
								break;
							default:break;
						} 
					}else{
						////已播放到课程结尾，课程结束
						roomInfo.RoomState = model.RoomState_End;//设置课程结束
						roomInfo.NeedPlayPreMedia = false;
						roomInfo.NeedPlayBack = false;

						roomInfo.MAllowNew = true;
						roomInfo.MCurrent = nil;
						roomInfo.MainFrames = nil;
						roomInfo.MCurrentTimeInterval = 0;
						roomInfo.MCurrentMainFrameIdx = 0;
						roomInfo.MPreMainFrameIdx = 0;
						roomInfo.MMainFramePlayInterval = 0;
						roomInfo.MMainFramePrePlayInterval = 0;

						roomInfo.SAllowNew = true;
						roomInfo.SCurrent = nil;
						roomInfo.SCurrentTimeInterval = 0;
						roomInfo.SCurrentQuesionTimeOut = -1;
						roomInfo.CurrentProcess = 0;//重置播放状态。使其播放媒体命令
						roomInfo.NeedJumpSleep = false;
						roomInfo.Credit = 0;
						continue;
					}
				}

				//根据mainFrames时间轴，解析stepData,并获取要执行的命令数组
				tempCmdArr,frameStepIdx,_ := execStepDataByMainFrames(roomInfo.MainFrames,stepDataArr,roomInfo,roomInfo.MCurrentMainFrameIdx,curTime)
				roomInfo.MCurrentMainFrameIdx = frameStepIdx;
				if len(tempCmdArr) > 0{
					cmdArr = append(cmdArr,tempCmdArr...);//将脚本塞入 下发列表
				}
				if len(cmdArr) > 0{
					//下发教学命令到客户端
					sendTeachScriptToUser(client,rid,cmdArr,roomInfo.AnswerUIDQueue);
					cmdArr = []map[string]interface{}{};//清空已发的命令集合
				}

				roomInfo.MCurrentTimeInterval += offsetTime;//无论是否正在执行教学脚本， 都要更新媒体脚本的计时，防止教学脚本执行完毕后，要过很久才能迎来新的媒体脚本播放时机
				client.GTimerInterval = curTime;//更新上一次处理脚本时的时间记录.
				if roomInfo.MCurrentTimeInterval >= roomInfo.MCompleteTime{
					//已经达到当前媒体脚本的播放结束时间，执行下一个媒体脚本
					roomInfo.MCurrentTimeInterval = 0;
					roomInfo.MAllowNew = true;
				}
			}
		}else{
			//计时，实时更新课程状态
			if curTime > roomInfo.StartTimeinterval && curTime < roomInfo.EndTimeinterval{
				roomInfo.RoomState = model.RoomState_Started
			}
		}
		client.GTimerInterval = curTime;//更新上一次处理脚本时的时间记录.
	}
	client.runLoopExecChan <- 1;
}

func tempFunc(stepItem * model.ScriptStepData,stData []model.ScriptStepData,rinfo *model.RoomInfo,source []map[string]interface{},curTime int64)(result []map[string]interface{},hasChangePage bool){
	tempI := 0;
	tresult,hasChangePageCount,templateItem := foreachScriptItem(stepItem,stData,rinfo,source,curTime,tempI);//通过递归获取最终要执行教学脚本数组
	if hasChangePageCount > 0{
		hasChangePage = true;
	}
	if nil != templateItem{
		//如果存在特殊教学脚本，则进行特殊处理
		itemType := templateItem.Type;
		itemValue := templateItem.Value;
		//重新计算脚本结束时间
		rinfo.SWaitAnswerUids = rinfo.UserIdArr;//设置应答序列
		var timeLength int64;
		if itemType == "templateCMD"{
			timeLength = getInt64(itemValue["timeout"],1) + 2;
		}else{
			timeLength = getInt64(itemValue["endSecond"],0) - getInt64(itemValue["beginSecond"],0);
		}
		rinfo.SAllowNew = false;//禁用关键帧脚本的执行
		rinfo.SCurrentQuesionTimeOut = rinfo.SCurrentTimeInterval + timeLength;
	}
	return tresult,hasChangePage;
}

func execStepDataByMainFrames(mainFrames []model.MediaMainFrame,stData []model.ScriptStepData,rinfo *model.RoomInfo,currentFrameStepIdx int64,curTime int64)(result []map[string]interface{},idx int64,hasChangePage bool){
	hasChangePageCount := 0;
	hasChangePage = false;
	stDataLength := len(stData);
	var templateItem *model.ScriptStepData = nil;//教学脚本中 需要特殊处理的项
	if nil != rinfo.SCurrent{
		callBackFuncID := getInt64(rinfo.SCurrent.Value[rinfo.CurrentAnswerState],-1);//取当前处理脚本对应用户操作的正确处理命令
		if callBackFuncID <= -1{
			callBackFuncID = rinfo.SCurrent.Next;
		}
		
		if callBackFuncID > -1 && int(callBackFuncID) < stDataLength{
			//如果存在已经执行完的教学脚本， 则尝试执行这个教学脚本的下一个脚本
			item := &stData[callBackFuncID];
			result,hasChangePage = tempFunc(item,stData,rinfo,result,curTime);
			//返回要执行的脚本和不变的媒体关键帧播放进度，以及是否有翻页命令存在
			return result,currentFrameStepIdx,hasChangePage;
		}else{
			//如果这个脚本没有next脚本，则执行下一批脚本
			rinfo.SCurrent = nil;
		}
	}

	currentPlayTime := rinfo.MCurrentTimeInterval;
	j := int64(len(mainFrames));
	if currentFrameStepIdx < j{
		arr := mainFrames[currentFrameStepIdx:j];//取出还未执行的媒体关键帧
		//遍历关键帧时间 <= currentPlayTime 出来进行播放
		for _,frame := range arr{
			ct := frame.StepTime
			if ct < 0{
				ct = currentPlayTime + 1;
			}
			if ct <= currentPlayTime{
				rinfo.MCurrentMainFrameIdx = currentFrameStepIdx;//记录当前媒体的播放关键帧，用于和之后的MPreMainFrameIdx配合使用记录changepage命令对应的关键帧，用于断线重连或者续播
				currentFrameStepIdx += 1;//关键帧向后移动一位
				//找到满足条件的关键帧，则取出对应的stepData命令
				sid := frame.StepId;
				//将媒体播放命令的已播放时间，更新至当前关键帧对应的媒体播放时间，与MMainFramePrePlayInterval配合使用 用于断线重连后的续播
				rinfo.MMainFramePlayInterval = frame.MediaTime;
				//处理可执行的教学脚本
				if sid > -1 && int(sid) < stDataLength{
					item := stData[int(sid)];
					//log.Println("处理,id=",sid," type=",item.Type," ct=",ct," currentPlayTime=",currentPlayTime);
					result,hasChangePageCount,templateItem = foreachScriptItem(&item,stData,rinfo,result,curTime,hasChangePageCount);//通过递归获取最终要执行教学脚本数组
					if hasChangePageCount > 0{
						hasChangePage = true;
					}
					if nil != templateItem{
						//如果存在特殊教学脚本，则进行特殊处理
						itemType := templateItem.Type;
						itemValue := templateItem.Value;
						//重新计算脚本结束时间
						rinfo.SWaitAnswerUids = rinfo.UserIdArr;//设置应答序列
						var timeLength int64;
						if itemType == "templateCMD"{
							timeLength = getInt64(itemValue["timeout"],1) + 2;
						}else{
							timeLength = getInt64(itemValue["endSecond"],0) - getInt64(itemValue["beginSecond"],0);
						}
						rinfo.SAllowNew = false;//禁用关键帧脚本的执行
						rinfo.SCurrentQuesionTimeOut = rinfo.SCurrentTimeInterval + timeLength;
						break;
					}
				}
			}
		}
	}
	return result,currentFrameStepIdx,hasChangePage;
}

func foreachScriptItem(item *model.ScriptStepData,stData []model.ScriptStepData,rinfo *model.RoomInfo,source []map[string]interface{},curTime int64,hasChangePageCount int)(result []map[string]interface{},pageCount int,template *model.ScriptStepData){
	itemType := item.Type;
	rinfo.SCurrentTimeInterval = 0;//重置教学脚本计时时间
	rinfo.SCurrent = item;//设置当前正在执行的教学脚本
	rinfo.CurrentProcess = 1;
	rinfo.CurrentAnswerState = "";
	subItem := map[string]interface{}{"suid":0,"playInterval":0,"st":curTime,"data":item};
	result = append(source,subItem);//添加到返回数组
	if itemType == "templateCMD" || itemType == "video" || itemType == "audio"{
		rinfo.SCurrentQuestionId = item.Id;//设置题号,用于答题匹配
		//如遇关键脚本，则直接返回数据
		return result,hasChangePageCount,item;			
	}else{
		if itemType == "changePage"{
			hasChangePageCount += 1;//记录有无 翻页命令存在
			rinfo.MPreMainFrameIdx = rinfo.MCurrentMainFrameIdx;//记录关键的翻页命令对应的 媒体播放关键帧，用于退出教室后的重播或续播逻辑
			//将媒体播放命令的已播放时间，更新至当前关键帧对应的媒体播放时间， 用于断线重连后的续播
			rinfo.MMainFramePrePlayInterval = rinfo.MMainFramePlayInterval;
			rinfo.PreChangePageCMD = subItem;
		}else if itemType == "star"{
			starCount := int(getInt64(item.Value["count"],0));//递增星星数量
			rinfo.Credit += starCount;
			rinfo.WantAddCredit -= starCount;
			if rinfo.WantAddCredit < 0 {
				rinfo.WantAddCredit = 0;
			}
		}
		if item.Next > -1 && int(item.Next) < len(stData){
			//如果当前脚本存在下一个脚本，则递归下一个脚本
			newItem := stData[item.Next];
			return foreachScriptItem(&newItem,stData,rinfo,result,curTime,hasChangePageCount);
		}else{
			return result,hasChangePageCount,nil;
		}
	}
}

/**
将媒体命令转换为 教学命令
*/
func mediaConverScript(mediaItem *model.MediaData)map[string]interface{}{
	result := map[string]interface{}{"id":mediaItem.Id,"type":mediaItem.Type,"value":mediaItem.Value}
	return result;
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

func objArrToStepDataArr(arr []map[string]interface{})[]model.ScriptStepData{
	var result []model.ScriptStepData;
	if nil != arr{
		for _,v := range arr{
			item := model.ScriptStepData{};
			item.Id = getInt64(v["id"],-1);
			item.Type = getString(v["type"],"");
			item.Next = getInt64(v["n"],-1);
			item.Pre = getInt64(v["p"],-1);
			item.Value = getMap(v["value"],nil);
			result = append(result,item);
		}
	}
	return result;
}

/*[]object转[]MediaData*/
func objArrToMediaDataArr(arr []map[string]interface{})[]model.MediaData{
	var result []model.MediaData;
	if nil != arr{
		for _,v := range arr{
			item := model.MediaData{};
			item.Id = getInt64(v["id"],-1);
			item.Type = getString(v["type"],"");
			item.Next = getInt64(v["n"],-1);
			item.Pre = getInt64(v["p"],-1);
			item.Value = getMap(v["value"],nil);
			item.MainFrames = objArrToMainFrames(getObjArray(v["mainFrames"],nil));
			result = append(result,item);
		}
	}
	return result;
}

func objArrToMainFrames(arr []map[string]interface{})[]model.MediaMainFrame{
	var result []model.MediaMainFrame;
	if nil != arr{
		for _,v := range arr{
			item := model.MediaMainFrame{};
			item.MediaTime = getInt64(v["mt"],-1);
			item.StepTime = getInt64(v["st"],-1);
			item.StepId = getInt64(v["sid"],-1);
			result = append(result,item);
		}
	}
	return result;
}

/*object转[]map[string]interface{}*/
func getfloat64Array(val interface{},def []float64)[]float64{
	if nil == val{
		return def;
	}
	tem,ok := val.([]interface{});
	if ok{
		result := []float64{};
		for _,v := range tem{
			resultV,ok:= v.(float64);
			if ok == true{
				result = append(result,resultV)
			}
		}
		return result;
	}else{
		return def;
	}
}
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
		// req.Rid;
		// req.TeachScript;
		// req.StartTimeinterval;
		// req.Uid;
		// req.NickName;
	}
}
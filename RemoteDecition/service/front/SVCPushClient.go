/**
svc 推送服务客户端
*/
package front

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	// 发送约课信息变动通知
	appoint_course int = 1

	// 发布通用用户push消息
	pub_comm_user_push_msg int = 2

	// 编辑已发布的通用用户push消息
	edit_comm_user_push_msg int = 3

	// 删除已发布的通用用户push消息
	del_comm_user_push_msg int = 4

	// 发送通用用户push消息给客户端
	send_comm_user_push_msg int = 5
)

var (
	svc_push_http_host     string = "172.16.70.13" //http push地址
	svc_push_http_port     uint16 = 10751          //http push 端口
	svc_push_get_token_url string = "/access_token"
	svc_push_send_msg      string = "/push/inst/b"
)

type SVCPushClient struct {
	app_id     string
	app_secret string
	token      string
}

//获取请求权限
func (svc *SVCPushClient) getAuth(callback func(string, string)) {
	if nil != callback {
		callback("DmhersGpsmq4mxDQnLSHE9", "DmhersGqNcRk6xyMLEVa3A") //目前是从svc服务那里手动指定的，之后会改成HTTP请求获得
	}

}

//获取token
func (svc *SVCPushClient) getToken(callback func(string)) {
	reqStr := fmt.Sprintf("app_id=%s&app_secret=%s&grant_type=authorization_code&code=FplxlWBeZQUYbYS6WxSgIA&redirect_uri=", svc.app_id, svc.app_secret)
	_databody := strings.NewReader(reqStr)
	_url := fmt.Sprintf("http://%s:%d%s", svc_push_http_host, svc_push_http_port, svc_push_get_token_url)
	go svc.sendHttpRequest(_url, _databody, "application/x-www-form-urlencoded", "", func(bts []byte) {
		resObj := map[string]interface{}{}
		err := json.Unmarshal(bts, &resObj)
		if nil == err {
			if _token, isOK := resObj["access_token"].(string); isOK == true && nil != callback {
				callback(_token)
			}
		}
	})

}

func (svc *SVCPushClient) PushMsg(msgObj interface{}) {
	svc.getAuth(func(_appid string, _appsecret string) {
		svc.app_id = _appid
		svc.app_secret = _appsecret
		svc.getToken(func(_token string) {
			svc.token = _token
			msgBts, err := json.Marshal(msgObj)
			if nil == err {
				msg := string(msgBts)
				reqBody := map[string]interface{}{
					"msg_type": 3,
					"content":  msg,
					"ext":      ""}
				byts, err := json.Marshal(reqBody)
				if nil == err {
					log.Println(string(byts))
					_databody := strings.NewReader(string(byts))
					_idUUID, e := uuid.NewV4()
					if nil == e {
						_id := fmt.Sprintf("%v", _idUUID)
						_id = base64.StdEncoding.EncodeToString([]byte(_id))
						_url := fmt.Sprintf("http://%s:%d%s?id=%s", svc_push_http_host, svc_push_http_port, svc_push_send_msg, _id)
						go svc.sendHttpRequest(_url, _databody, "application/json", svc.token, nil)
					}

				}
			}
		})
	})
}

func (svc *SVCPushClient) sendHttpRequest(url string, dataBody io.Reader, contentType string, token string, callback func([]byte)) {
	http_req, reqErr := http.NewRequest("POST", url, dataBody)
	http_req.Header.Set("Content-Type", contentType)
	if token != "" {
		http_req.Header.Set("Authorization", "Bearer "+token)
	}
	if nil != reqErr {
		log.Println(reqErr.Error())
		return
	}
	http_client := &http.Client{Timeout: time.Second * 30} //设置默认超时为30秒
	res, err := http_client.Do(http_req)
	if nil != err {
		log.Println(err.Error())
	} else {
		resBytes, err := ioutil.ReadAll(res.Body)
		if nil != err {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("收到svc回调:%s", string(resBytes)))
			if nil != callback {
				callback(resBytes)
			}
		}
	}
}

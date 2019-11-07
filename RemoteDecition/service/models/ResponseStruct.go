package models


type CurrentResponse struct{
	Code string `json:"code"`
	Msg string `json:"msg"`
}

type DataResponse struct{
	Code string `json:"code"`
	Msg string `json:"msg"`
	Data []map[string]interface{} `json:"data"`
}


type ModelRDRole struct{
	AuthGroup string `json:"authGroup"`
	ID uint16 `json:"id"`
	RoleName string `json:"roleName"`
	RoleDes string `json:"roleDes"`
	Remark string `json:"remark"`
}

/*客户端心跳*/
type HeartBeat_c2s struct{
    Cmd uint32 `json:"cmd"`
    Seq uint32 `json:"seq"`;//数据包的序号，可以为0
    LocalTimeinterval uint32 `json:"lt"`;//客户端发送请求时的UTC时间的秒值
}
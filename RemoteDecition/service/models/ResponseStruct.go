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
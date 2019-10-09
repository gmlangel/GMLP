package models


type CurrentResponse struct{
	Code string
	Msg string
}

type DataResponse struct{
	Code string
	Msg string
	Data []map[string]interface{}
}


type ModelRDRole struct{
	AuthGroup string `json:"authGroup"`
	ID uint16 `json:"id"`
	RoleName string `json:"roleName"`
	RoleDes string `json:"roleDes"`
	Remark string `json:"remark"`
}
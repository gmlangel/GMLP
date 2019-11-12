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


type ForceStrategyBeUse_c2s struct{
	Cmd uint32 `json:"cmd"`
	StrategyPath string `json:"strategyPath"`
	ConditionPath string `json:"conditionPath"`
	StrategyID uint64 `json:"strategyId"`
}

/*策略变更协议*/
type StrategyChanged_c2s struct{
	Cmd uint32 `json:"cmd"`
	StrategyPath string `json:"strategyPath"`
}

//条件变更协议
type ConditionChanged_c2s struct{
	Cmd uint32 `json:"cmd"`
	ConditionPath string `json:"conditionPath"`
}

type StrategyInfo struct{
	Id uint64 `json:"id"`
	Sid uint64 `json:"sid"`
	Cgroup string `json:"conditionGroup"`
	ValuePath string `json:"strategyPath"`
	Enabled uint64 `json:"enabled"`
	ExpireDate uint64 `json:"expireDate"`
	Name string `json:"name"`
	LastUpdate string `json:"lastUpdate"`
}

type ConditionInfo struct{
	Id uint64 `json:"id"`
	Typeid uint64 `json:"typeID"`
	TypeName string `json:"typeName"`
	Value string `json:"value"`
	Operator string `json:"operator"`
	Probability float64 `json:"probability"`
	LastUpdate string `json:"lastUpdate"`
}

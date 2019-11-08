package models

//MARK:Socket 通讯 数据包命令ID定义
const(
	//策略更新
	C_REQ_S_STRATEGYCHANGED = 0x00FF003C;
	S_NOTIFY_C_STRATEGYCHANGED = 0x00FF003E;

	//条件更新
	C_REQ_S_CONDITIONCHANGED = 0x00FF003F;
	S_NOTIFY_C_CONDITIONCHANGED = 0x00FF0040;
)

/*策略变更协议*/
type StrategyChanged_c2s struct{
	Cmd uint32 `json:"cmd"`
	StrategyPath string `json:"strategyPath"`
	IdArr []uint64 `json:"idArr"`
	Type string `json:"type"`
}

//条件变更协议
type ConditionChanged_c2s struct{
	Cmd uint32 `json:"cmd"`
	ConditionPath string `json:"conditionPath"`
	IdArr []uint64 `json:"idArr"`
	Type string `json:"type"`
}

type StrategyChanged_s2c_notify struct{
	Cmd uint32 `json:"cmd"`
	StrategyPath string `json:"strategyPath"`
	IdArr []uint64 `json:"idArr"`
	Type string `json:"type"`
}

type ConditionChanged_s2c_notify struct{
	Cmd uint32 `json:"cmd"`
	ConditionPath string `json:"conditionPath"`
	IdArr []uint64 `json:"idArr"`
	Type string `json:"type"`
}
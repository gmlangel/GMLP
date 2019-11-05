package models

//MARK:Socket 通讯 数据包命令ID定义
const(
	//策略更新
	C_REQ_S_STRATEGYCHANGED = 0x00FF003C;
	S_NOTIFY_C_STRATEGYCHANGED = 0x00FF003E;
)

type StrategyChanged_c2s struct{
	ConditionPath string `json:"conditionPath"`
	StrategyPath string `json:"strategyPath"`
	Msg string `json:"msg"`
}

type StrategyChanged_s2c_notify struct{
	ConditionPath string `json:"conditionPath"`
	StrategyPath string `json:"strategyPath"`
	Msg string `json:"msg"`
}
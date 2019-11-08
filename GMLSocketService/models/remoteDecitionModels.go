package models

//MARK:Socket 通讯 数据包命令ID定义
const(
	//策略更新
	C_REQ_S_STRATEGYCHANGED = 0x00FF003C;
	S_NOTIFY_C_STRATEGYCHANGED = 0x00FF003E;

	//条件更新
	C_REQ_S_CONDITIONCHANGED = 0x00FF003F;
	S_NOTIFY_C_CONDITIONCHANGED = 0x00FF0040;

	//强制策略即时生效
	C_REQ_ForceStrategyBeUseage = 0x00FF0041;
	S_NOTIFY_C_ForceStrategyBeUseage = 0x00FF0042;
)

type ForceStrategyBeUse_c2s struct{
	Cmd uint32 `json:"cmd"`
	StrategyPath string `json:"strategyPath"`
	ConditionPath string `json:"conditionPath"`
	StrategyID uint64 `json:"strategyId"`
}

type ForceStrategyBeUse_s2c_notify struct{
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

type StrategyChanged_s2c_notify struct{
	Cmd uint32 `json:"cmd"`
	StrategyPath string `json:"strategyPath"`
}

type ConditionChanged_s2c_notify struct{
	Cmd uint32 `json:"cmd"`
	ConditionPath string `json:"conditionPath"`
}
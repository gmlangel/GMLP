package tools;

import(
	"encoding/json"
)
/**
struct转字符串
*/
func StructToJSONStr(argStruct interface{})string{
	bytes,jsonerr := json.Marshal(argStruct);
	if jsonerr == nil{
		return string(bytes);
	}else{
		return "{}";
	}
}
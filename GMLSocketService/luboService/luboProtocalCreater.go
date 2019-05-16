package luboService;
import(
	model "../models"
)

func CreateProtocal(cmd uint32)interface{}{
	var result interface{} = nil;
	switch cmd{
	case model.S_RES_C_HY:
		result = createClientIn_s2c();
		break;
	default:break;
	}
	return result;
}

func createClientIn_s2c()interface{}{
	result := &model.ClientIn_s2c{Cmd:model.S_RES_C_HY,Des:"欢迎加入AI录播服务器"};
	return result;
}


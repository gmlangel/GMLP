package front
import(
	"github.com/kataras/iris"
	"fmt"
	m "../models"
	"encoding/json"
	"strconv"
)


type AllService struct{
	SQL m.SQLInterface;
} 
/**
添加后台用户
*/
func (ser *AllService)AddUser(ctx iris.Context){
	fmt.Println("添加成功");
	ln := ctx.URLParam("ln");
	pwd := ctx.URLParam("pwd");
	authType,err := ctx.URLParamInt("authType");
	response := m.CurrentResponse{};
	if "" != ln && "" != pwd && -1 < authType && nil == err{
		_,err := ser.SQL.Exec(fmt.Sprintf("insert into RDUser(`signName`,`signPWD`,`roleID`) Values('%s','%s',%d)",ln,pwd,authType))
		if nil != err{
			response.Code = "-1";
			response.Msg = fmt.Sprintf("数据库插入数据出错，%v",err);
		}else{
			response.Code = "0";
			response.Msg = "账号创建成功";
		}
		res,err := json.Marshal(response)
		if nil != err{
			res = []byte("");
		}else{
			ctx.Write(res)
		}
	}else{
		response.Code = "-1";
		response.Msg = "添加用户失败，请检查参数的正确性";
		res,err := json.Marshal(response)
		if nil != err{
			res = []byte("");
		}
		ctx.Write(res);
	}
}

/**
获取所有的可选角色类型
*/
func(ser *AllService)GetAllRoleType(ctx iris.Context){
	result,err := ser.SQL.Query("select * from RDRole");
	response := &m.DataResponse{};
	if nil != err{
		response.Code = "-1";
		response.Msg = fmt.Sprintf("获取角色类型失败，%v",err);
	}else{
		response.Code = "0";
		//遍历角色数组
		j:= len(result);
		tmpArr := []map[string]interface{}{};
		for i:=0;i<j;i++{
			item := result[i];
			decodedItem := map[string]interface{}{};
			for k,v := range(item){
				if "id" == k{
					decodedItem[k],err= strconv.ParseUint(string(v),0,16);
				}else{
					decodedItem[k] = string(v);
				}
			}
			tmpArr = append(tmpArr,decodedItem)
		}
		response.Data = tmpArr;
	}
	res,err := json.Marshal(response)
	if nil != err{
		ctx.Write([]byte(""))
	}else{
		ctx.Write(res);
	}
}

/**
获取所有权限的说明
*/
func (ser *AllService)GetAllAuth(ctx iris.Context){
	queryResult,err := ser.SQL.Query("select * from RDAuth");
	res := &m.DataResponse{};
	if nil != err{
		res.Code = "-1";
		res.Msg = fmt.Sprintf("获取权限失败,%s",err);
	}else{
		//遍历所有权限
		authArr := []map[string]interface{}{};
		for _,v := range(queryResult){
			item := map[string]interface{}{};
			for k,tv := range(v){
				if "id" == k{
					item[k],err= strconv.ParseUint(string(tv),0,16);
				}else{
					item[k] = string(tv);
				}
			}
			authArr = append(authArr,item);
		}
		res.Code = "0";
		res.Data = authArr;
	}
	resBytes,err := json.Marshal(res);
	if nil != err{
		ctx.Write([]byte(""))
	}else{
		ctx.Write(resBytes)
	}
	
}

/**
获取所有的后台管理账号（支持分页查询），默认为查询前20条数据
*/
func (ser *AllService)GetAllUsers(ctx iris.Context){
	sp,err:= ctx.URLParamInt("startPoint");
	rc,err2:= ctx.URLParamInt("readCount");
	queryStr := "";
	if nil == err && nil == err2{
		//分页查询
		queryStr = fmt.Sprintf("select * from RDUser limit %d,%d",sp,rc)
	}else{
		//查询前20
		queryStr = fmt.Sprintf("select * from RDUser limit %d,%d",0,20)
	}
	queryResult,err := ser.SQL.Query(queryStr);
	res := &m.DataResponse{};
	if nil != err{
		res.Code = "-1";
		res.Msg = fmt.Sprintf("获取用户信息失败,%s",err);
	}else{
		//遍历用户列表
		userArr := []map[string]interface{}{};
		for _,v := range(queryResult){
			item := map[string]interface{}{};
			for k,tv := range(v){
				if "id" == k{
					item[k],err= strconv.ParseUint(string(tv),0,32)
				}else if "roleID" == k{
					item[k],err= strconv.ParseUint(string(tv),0,16)
				}else{
					item[k] = string(tv);
				}
			}
			userArr = append(userArr,item);
		}
		res.Code = "0";
		res.Data = userArr;
	}
	resBytes,err := json.Marshal(res);
	if nil != err{
		ctx.Write([]byte(""))
	}else{
		ctx.Write(resBytes);
	}
}
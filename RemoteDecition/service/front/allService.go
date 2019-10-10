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

/**
删除后台管理账号
*/
func (ser *AllService)DeleteUser(ctx iris.Context){
	userID,err:= strconv.ParseUint(ctx.URLParam("uid"),0,32);
	res := &m.CurrentResponse{}
	if nil != err{
		res.Code = "-1";
		res.Msg = "删除失败,参数uid有问题"
	}else{
		result,err := ser.SQL.Exec(fmt.Sprintf("Delete from RDUser where `id`= %d",userID))
		if nil != err{
			res.Code = "-1";
			res.Msg = fmt.Sprintf("删除失败,%v",err)
		}else if count,err := result.RowsAffected();err == nil&&count > 0{
			res.Code = "0";
			res.Msg = "删除成功"
		}else{
			res.Code = "-1";
			res.Msg = "删除失败,用户不存在"
		}
	}
	resBytes,err := json.Marshal(res);
	if nil != err{
		ctx.Write([]byte(""))
	}else{
		ctx.Write(resBytes);
	}
}

/**
更新用户信息
*/
func (ser *AllService)UpdateUserInfo(ctx iris.Context){
	signName := ctx.URLParam("sn");
	signPWD := ctx.URLParam("sp");
	roleID,err:= ctx.URLParamInt("rid");
	uid,err2:=ctx.URLParamInt("uid");
	res := &m.CurrentResponse{};
	if "" != signName && "" != signPWD && nil == err && nil == err2{
		//开始更新数据库
		result,err := ser.SQL.Exec(fmt.Sprintf("update RDUser set `signName` = '%s' , `signPWD` = '%s',`roleID` = %d  where `id` = %d",signName,signPWD,roleID,uid))
		if nil != err{
			res.Code = "-1";
			res.Msg = fmt.Sprintf("用户信息更新失败,%s",err);
		}else if count,e:=result.RowsAffected();nil == e && count > 0{
			res.Code = "0";
			res.Msg = "用户信息更新成功"
		}else{
			res.Code = "-1";
			res.Msg = "用户信息更新失败,用户不存在";
		}
	}else{
		res.Code = "-1";
		res.Msg = "更新用户数据失败，参数有问题";
	}
	resBytes,err := json.Marshal(res);
	if nil != err{
		ctx.Write([]byte(""));
	}else{
		ctx.Write(resBytes);
	}
}

/**
查询所有策略条件信息,支持分页查询，默认返回最近的20条信息
*/
func (ser *AllService)GetAllConditionInfo(ctx iris.Context){
	startPoint,err1 := ctx.URLParamInt("startPoint");
	readCount,err2 := ctx.URLParamInt("readCount");
	queryStr := "";
	res := &m.DataResponse{};
	if nil == err1 && nil == err2{
		//开始分页查询
		queryStr = fmt.Sprintf("select `Condition`.`id`,`ConditionType`.`zhName` ,`ConditionType`.`des` AS `typeDes` ,`Condition`.`value` ,`Condition`.`name` ,`Condition`.`probability` ,`Condition`.`des`    from   `Condition` LEFT JOIN   `ConditionType`   on   `Condition`.`typeID` = `ConditionType`.`id` ORDER BY `Condition`.`id` DESC LIMIT %d,%d;",startPoint,readCount)
	}else{
		//查询前20条
		queryStr = fmt.Sprintf("select `Condition`.`id`,`ConditionType`.`zhName` ,`ConditionType`.`des` AS `typeDes` ,`Condition`.`value` ,`Condition`.`name` ,`Condition`.`probability` ,`Condition`.`des`    from   `Condition` LEFT JOIN   `ConditionType`   on   `Condition`.`typeID` = `ConditionType`.`id` ORDER BY `Condition`.`id` DESC LIMIT %d,%d;",0,20)
	}
	result,err := ser.SQL.Query(queryStr);
	if nil != err{
		res.Code = "-1";
		res.Msg = fmt.Sprintf("获取数据错误，%v",err);
	}else{
		conditionArr := []map[string]interface{}{};
		for _,v := range(result){
			item := map[string]interface{}{};
			for k,nv := range(v){
				if "id" == k{
					item[k],err= strconv.ParseUint(string(nv),0,32);
				}else{
					item[k] = string(nv);
				}
			}
			conditionArr = append(conditionArr,item);
		}
		res.Code = "0";
		res.Data = conditionArr;
	}
	resBytes,err:=json.Marshal(res);
	if nil != err{
		ctx.Write([]byte(""));
	}else{
		ctx.Write(resBytes)
	}
}

/**
获取所有条件类型信息的接口
*/
func (ser *AllService)GetAllConditionTypeInfo(ctx iris.Context){
	result,err := ser.SQL.Query("select * from `ConditionType`")
	res := &m.DataResponse{};
	if nil != err{
		res.Code = "-1";
		res.Msg = fmt.Sprintf("数据读取错误,%v",err);
	}else{
		res.Code = "0";
		arr := []map[string]interface{}{};
		for _,v := range(result){
			item := map[string]interface{}{};
			for k,nv := range(v){
				if "id" == k{
					item[k],err= strconv.ParseUint(string(nv),0,32);
				}else{
					item[k] = string(nv);
				}
			}
			arr = append(arr,item)
		}
		res.Data = arr;
	}
	resBytes,err := json.Marshal(res);
	if nil != err{
		ctx.Write([]byte(""))
	}else{
		ctx.Write(resBytes);
	}
}

/**
添加条件的类型
*/
func (ser *AllService)AddConditionType(ctx iris.Context){
	zn := ctx.URLParam("zn");
	en := ctx.URLParam("en");
	des := ctx.URLParam("des");
	res := &m.CurrentResponse{};
	if "" != zn && "" != en{
		//执行插入
		_,err := ser.SQL.Exec(fmt.Sprintf("insert into `ConditionType`(`zhName`,`enName`,`des`) values('%s','%s','%s')",zn,en,des))
		if nil != err{
			res.Code = "-1";
			res.Msg = fmt.Sprintf("添加条件类型失败,%v",err);
		}else{
			res.Code = "0";
			res.Msg = "条件类型添加成功"
		}
	}else{
		res.Code = "-1";
		res.Msg = "添加条件类型失败，请检查其请求参数"
	}
	resBytes,err := json.Marshal(res);
	if nil != err{
		ctx.Write([]byte(""))
	}else{
		ctx.Write(resBytes);
	}
}

/**
添加条件信息
*/
func (ser *AllService)AddCondition(ctx iris.Context){
	cType,err1:= ctx.URLParamInt("cType");//条件类型
	name := ctx.URLParam("name");//条件名称
	val:= ctx.URLParam("value");//条件对应的值
	probability,err2 := strconv.ParseFloat(ctx.URLParam("probability"),32);//条件生效几率
	des := ctx.URLParam("des");//条件描述
	res := &m.CurrentResponse{};
	if nil == err1 && nil == err2 && "" != name && "" != val{
		//写入数据库
		_,err:= ser.SQL.Exec(fmt.Sprintf("insert into `Condition`(`typeID`,`value`,`name`,`probability`,`des`) values(%d,'%s','%s',%f,'%s')",cType,val,name,probability,des))
		if nil != err{
			res.Code = "-1";
			res.Msg = fmt.Sprintf("新增条件失败,%v",err)
		}else{
			res.Code = "0"
			res.Msg = "新增条件成功"
		}
	}else{
		res.Code = "-1";
		res.Msg = "条件添加失败，请检查参数"
	}

	resBytes,err:=json.Marshal(res);
	if nil != err{
		ctx.Write([]byte(""))
	}else{
		ctx.Write(resBytes);
	}
}
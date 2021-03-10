package front

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"sync"

	m "../models"
	"github.com/kataras/iris"
	uuid "github.com/satori/go.uuid"
)

type AllService struct {
	SQL                 m.SQLInterface
	Sock                *RDSocket
	SVCClient           *SVCPushClient
	StrategyConfigPath  string
	ConditionConfigPath string
	strategyArr         []m.StrategyInfo
	conditionArr        []m.ConditionInfo
	strategyChan        chan int
	conditionChan       chan int
	strategyChange      bool
	conditionChange     bool

	tmpConditionTypeInfoArr  []map[string]interface{}
	tmpConditionTypeInfoLock *sync.Mutex
}

//Socket长连接处理相关---------------------------------------------------------------------------------
/**
初始化数据*/
func (ser *AllService) InitData() {
	ser.tmpConditionTypeInfoArr = []map[string]interface{}{}
	ser.tmpConditionTypeInfoLock = &sync.Mutex{}
	ser.strategyChan = make(chan int, 1)
	ser.conditionChan = make(chan int, 1)
	ser.strategyChan <- 1
	ser.conditionChan <- 1
	//请求所有的策略信息，并生成配置文件
	ser._getAllStrategyInfo()
	//请求所有的条件信息，并生成配置文件
	ser._getAllConditionInfo()
	ser.Sock.onLinkComplete = func() {
		//断线重连后的重发机制
		ser.sendStrategyChangedToServer()
		ser.sendConditionChangedToServer()
	}
	//开启监听 策略和条件 的更新通知循环
	go ser.syncStrategyLoop()
	go ser.syncConditionLoop()
}

//开启监听 策略和条件 的更新通知循环，每隔120秒检测一次
func (ser *AllService) syncStrategyLoop() {
	//var item m.SyncLoopChan;
	for {
		_, isOk := <-ser.strategyChan
		if isOk == false {
			break
		}
		if true == ser.strategyChange {
			ser._saveStrategyConfig()
			ser.sendStrategyChangedToServer()
			ser.strategyChange = false
		}
		ser.strategyChan <- 1
		time.Sleep(time.Second * 120) //每隔120秒，检测一次
	}
}

func (ser *AllService) syncConditionLoop() {
	//var item m.SyncLoopChan;
	for {
		_, isOk2 := <-ser.conditionChan
		if isOk2 == false {
			break
		}
		if true == ser.conditionChange {
			ser._saveConditionConfig()
			ser.sendConditionChangedToServer()
			ser.conditionChange = false
		}
		ser.conditionChan <- 1
		time.Sleep(time.Second * 120) //每隔120秒，检测一次
	}
}

/**
通知长连接服务，策略已更新
@param gType string "sync" = 同步 ,"add" = 新增 ,"delete" = 删除 , "update" = 更新
*/
func (ser *AllService) sendStrategyChangedToServer() {
	//通知长连接服务器，策略文件与条件配置文件有更新
	req := &m.StrategyChanged_c2s{Cmd: 0x00FF003C, StrategyPath: ser.StrategyConfigPath}
	ser.Sock.Write(req)
	ser.SVCClient.PushMsg(req)
}

/**
通知长连接服务，条件已更新
@param gType string "sync" = 同步 ,"add" = 新增 ,"delete" = 删除 , "update" = 更新
*/
func (ser *AllService) sendConditionChangedToServer() {
	//通知长连接服务器，策略文件与条件配置文件有更新
	req := &m.ConditionChanged_c2s{Cmd: 0x00FF003F, ConditionPath: ser.ConditionConfigPath}
	ser.Sock.Write(req)
	ser.SVCClient.PushMsg(req)
}

func (ser *AllService) _getAllStrategyInfo() {
	queryStr := "select `Strategy`.`id`,`Strategy`.`sid`,`Strategy`.`conditionGroup`,`Strategy`.`valuePath`,`Strategy`.`enabled`,`Strategy`.`expireDate`,`Strategy`.`lastUpdate`,`Strategy`.`name`,`StrategyCategroy`.`name` as `sname` from `Strategy` left join `StrategyCategroy` on `Strategy`.`sid` = `StrategyCategroy`.`id`"
	res, err := ser.SQL.Query(queryStr)
	var item map[string][]byte
	ser.strategyArr = []m.StrategyInfo{}
	if nil == err {
		for _, v := range res {
			item = v
			si := m.StrategyInfo{}
			si.Id, err = strconv.ParseUint(string(item["id"]), 10, 32)
			si.Sid, err = strconv.ParseUint(string(item["sid"]), 10, 32)
			si.Cgroup = string(item["conditionGroup"])
			si.ValuePath = string(item["valuePath"])
			si.Enabled, err = strconv.ParseUint(string(item["enabled"]), 10, 32)
			si.ExpireDate, err = strconv.ParseUint(string(item["expireDate"]), 10, 32)
			si.Name = string(item["name"])
			si.Sname = string(item["sname"])
			si.LastUpdate = string(item["lastUpdate"])
			ser.strategyArr = append(ser.strategyArr, si)
		}
	}
	ser._saveStrategyConfig()         //更新本地文件
	ser.sendStrategyChangedToServer() //通知长连接
}

func (ser *AllService) _saveStrategyConfig() {
	//写入本地文件
	content, err1 := json.Marshal(ser.strategyArr)
	if nil == err1 {
		//删除旧的策略配置文件
		err3 := os.RemoveAll("./static/strategyConfig/")
		if nil != err3 {
			fmt.Println(err3.Error())
		}
		//更新策略的配置文件
		filepath, err2 := ser.writeToFile(string(content), "./static/strategyConfig/")
		if nil == err2 {
			ser.StrategyConfigPath = filepath
		}
	}
}

func (ser *AllService) _getAllConditionInfo() {
	queryStr := "select `Condition`.`id`,`Condition`.`lastUpdate`,`Condition`.`typeID`,`ConditionType`.`enName` ,`Condition`.`value` ,`Condition`.`probability`,`Condition`.`operator` from   `Condition` LEFT JOIN   `ConditionType`   on   `Condition`.`typeID` = `ConditionType`.`id`"
	res, err := ser.SQL.Query(queryStr)
	ser.conditionArr = []m.ConditionInfo{}
	if nil == err {
		for _, v := range res {
			cItem := m.ConditionInfo{}
			cItem.Id, _ = strconv.ParseUint(string(v["id"]), 10, 32)
			cItem.Typeid, _ = strconv.ParseUint(string(v["typeID"]), 10, 32)
			cItem.TypeName = string(v["enName"])
			cItem.Value = string(v["value"])
			cItem.Operator = string(v["operator"])
			cItem.Probability, _ = strconv.ParseFloat(string(v["probability"]), 32)
			cItem.LastUpdate = string(v["lastUpdate"])
			ser.conditionArr = append(ser.conditionArr, cItem)
		}
	}
	ser._saveConditionConfig()
	ser.sendConditionChangedToServer() //通知长连接，条件发生变更
}

func (ser *AllService) _saveConditionConfig() {
	content, err2 := json.Marshal(ser.conditionArr)
	if nil == err2 {
		//删除原有的配置，
		os.RemoveAll("./static/conditionConfig/")
		//存储新的配置
		fileName, err3 := ser.writeToFile(string(content), "./static/conditionConfig/")
		if nil == err3 {
			ser.ConditionConfigPath = fileName

		}
	}
}

//Webservice处理相关---------------------------------------------------------------------------------
/**
添加后台用户
*/
func (ser *AllService) AddUser(ctx iris.Context) {
	fmt.Println("添加成功")
	ln := ctx.URLParam("ln")
	pwd := ctx.URLParam("pwd")
	roleID, err := ctx.URLParamInt("roleID")
	response := m.CurrentResponse{}
	if "" != ln && "" != pwd && -1 < roleID && nil == err {
		_, err := ser.SQL.Exec(fmt.Sprintf("insert into RDUser(`signName`,`signPWD`,`roleID`) Values('%s','%s',%d)", ln, pwd, roleID))
		if nil != err {
			response.Code = "-1"
			response.Msg = fmt.Sprintf("数据库插入数据出错，%v", err)
		} else {
			response.Code = "0"
			response.Msg = "账号创建成功"
		}
		res, err := json.Marshal(response)
		if nil != err {
			res = []byte("")
		} else {
			ctx.Write(res)
		}
	} else {
		response.Code = "-1"
		response.Msg = "添加用户失败，请检查参数的正确性"
		res, err := json.Marshal(response)
		if nil != err {
			res = []byte("")
		}
		ctx.Write(res)
	}
}

/**
获取所有的可选角色类型
*/
func (ser *AllService) GetAllRoleType(ctx iris.Context) {
	result, err := ser.SQL.Query("select * from RDRole")
	response := &m.DataResponse{}
	if nil != err {
		response.Code = "-1"
		response.Msg = fmt.Sprintf("获取角色类型失败，%v", err)
	} else {
		response.Code = "0"
		//遍历角色数组
		j := len(result)
		tmpArr := []map[string]interface{}{}
		for i := 0; i < j; i++ {
			item := result[i]
			decodedItem := map[string]interface{}{}
			for k, v := range item {
				if "id" == k {
					decodedItem[k], err = strconv.ParseUint(string(v), 10, 16)
				} else {
					decodedItem[k] = string(v)
				}
			}
			tmpArr = append(tmpArr, decodedItem)
		}
		response.Data = tmpArr
	}
	res, err := json.Marshal(response)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(res)
	}
}

/**
获取所有权限的说明
*/
func (ser *AllService) GetAllAuth(ctx iris.Context) {
	queryResult, err := ser.SQL.Query("select * from RDAuth")
	res := &m.DataResponse{}
	if nil != err {
		res.Code = "-1"
		res.Msg = fmt.Sprintf("获取权限失败,%s", err)
	} else {
		//遍历所有权限
		authArr := []map[string]interface{}{}
		for _, v := range queryResult {
			item := map[string]interface{}{}
			for k, tv := range v {
				if "id" == k {
					item[k], err = strconv.ParseUint(string(tv), 10, 16)
				} else {
					item[k] = string(tv)
				}
			}
			authArr = append(authArr, item)
		}
		res.Code = "0"
		res.Data = authArr
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}

}

/**
获取后台管理账号数的总和
*/
func (ser *AllService) GetConditionCount(ctx iris.Context) {
	queryStr := "select COUNT(*) as cCount from `Condition`"
	queryResult, err := ser.SQL.Query(queryStr)
	res := &m.CurrentResponse{}
	if nil != err || len(queryResult) == 0 {
		res.Code = "-1"
		res.Msg = fmt.Sprintf("获取条件总数失败,%s", err)
	} else {
		res.Code = "0"
		count, _ := strconv.ParseUint(string(queryResult[0]["cCount"]), 10, 64)
		res.Msg = fmt.Sprintf("%d", count)
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
获取后台管理账号数的总和
*/
func (ser *AllService) GetUsersCount(ctx iris.Context) {
	queryStr := "select COUNT(*) as userCount from `RDUser`"
	queryResult, err := ser.SQL.Query(queryStr)
	res := &m.CurrentResponse{}
	if nil != err || len(queryResult) == 0 {
		res.Code = "-1"
		res.Msg = fmt.Sprintf("获取用户总数失败,%s", err)
	} else {
		//遍历用户列表
		res.Code = "0"
		count, _ := strconv.ParseUint(string(queryResult[0]["userCount"]), 10, 64)
		res.Msg = fmt.Sprintf("%d", count)
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
获取所有的后台管理账号（支持分页查询），默认为查询前20条数据
*/
func (ser *AllService) GetAllUsers(ctx iris.Context) {
	sp, err := ctx.URLParamInt("startPoint")
	rc, err2 := ctx.URLParamInt("readCount")
	queryStr := ""
	if nil == err && nil == err2 {
		//分页查询
		queryStr = fmt.Sprintf("select * from RDUser order by `id` desc limit %d,%d", sp, rc)
	} else {
		//查询前20
		queryStr = fmt.Sprintf("select * from RDUser order by `id` desc limit %d,%d", 0, 20)
	}
	queryResult, err := ser.SQL.Query(queryStr)
	res := &m.DataResponse{}
	if nil != err {
		res.Code = "-1"
		res.Msg = fmt.Sprintf("获取用户信息失败,%s", err)
	} else {
		//遍历用户列表
		userArr := []map[string]interface{}{}
		for _, v := range queryResult {
			item := map[string]interface{}{}
			for k, tv := range v {
				if "id" == k {
					item[k], err = strconv.ParseUint(string(tv), 10, 32)
				} else if "roleID" == k {
					item[k], err = strconv.ParseUint(string(tv), 10, 16)
				} else {
					item[k] = string(tv)
				}
			}
			userArr = append(userArr, item)
		}
		res.Code = "0"
		res.Data = userArr
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
删除后台管理账号
*/
func (ser *AllService) DeleteUser(ctx iris.Context) {
	userID, err := strconv.ParseUint(ctx.URLParam("uid"), 10, 32)
	res := &m.CurrentResponse{}
	if nil != err {
		res.Code = "-1"
		res.Msg = "删除失败,参数uid有问题"
	} else {
		result, err := ser.SQL.Exec(fmt.Sprintf("Delete from RDUser where `id`= %d", userID))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("删除失败,%v", err)
		} else if count, err := result.RowsAffected(); err == nil && count > 0 {
			res.Code = "0"
			res.Msg = "删除成功"
		} else {
			res.Code = "-1"
			res.Msg = "删除失败,用户不存在"
		}
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
更新用户信息
*/
func (ser *AllService) UpdateUserInfo(ctx iris.Context) {
	signName := ctx.URLParam("sn")
	signPWD := ctx.URLParam("sp")
	roleID, err := ctx.URLParamInt("rid")
	uid, err2 := ctx.URLParamInt("uid")
	res := &m.CurrentResponse{}
	if "" != signName && "" != signPWD && nil == err && nil == err2 {
		//开始更新数据库
		result, err := ser.SQL.Exec(fmt.Sprintf("update RDUser set `signName` = '%s' , `signPWD` = '%s',`roleID` = %d  where `id` = %d", signName, signPWD, roleID, uid))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("用户信息更新失败,%s", err)
		} else if count, e := result.RowsAffected(); nil == e && count > 0 {
			res.Code = "0"
			res.Msg = "用户信息更新成功"
		} else {
			res.Code = "-1"
			res.Msg = "用户信息更新失败,用户不存在或用户信息未被更改导致：提交更新失败"
		}
	} else {
		res.Code = "-1"
		res.Msg = "更新用户数据失败，参数有问题"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
查询所有策略条件信息,支持分页查询，默认返回全部信息
*/
func (ser *AllService) GetAllConditionInfo(ctx iris.Context) {
	startPoint, err1 := ctx.URLParamInt("startPoint")
	readCount, err2 := ctx.URLParamInt("readCount")
	queryStr := ""
	res := &m.DataResponse{}
	if nil == err1 && nil == err2 {
		//开始分页查询
		queryStr = fmt.Sprintf("select `Condition`.`id`,`Condition`.`typeID`,`ConditionType`.`zhName` ,`ConditionType`.`des` AS `typeDes` ,`Condition`.`value` ,`Condition`.`name` ,`Condition`.`probability`,`Condition`.`operator` ,`Condition`.`des`    from   `Condition` LEFT JOIN   `ConditionType`   on   `Condition`.`typeID` = `ConditionType`.`id` ORDER BY `Condition`.`id` DESC LIMIT %d,%d;", startPoint, readCount)
	} else {
		//查询前20条
		queryStr = fmt.Sprintf("select `Condition`.`id`,`Condition`.`typeID`,`ConditionType`.`zhName` ,`ConditionType`.`des` AS `typeDes` ,`Condition`.`value` ,`Condition`.`name` ,`Condition`.`probability`,`Condition`.`operator` ,`Condition`.`des`    from   `Condition` LEFT JOIN   `ConditionType`   on   `Condition`.`typeID` = `ConditionType`.`id` ORDER BY `Condition`.`id` DESC")
	}
	result, err := ser.SQL.Query(queryStr)
	if nil != err {
		res.Code = "-1"
		res.Msg = fmt.Sprintf("获取数据错误，%v", err)
	} else {
		conditionArr := []map[string]interface{}{}
		for _, v := range result {
			item := map[string]interface{}{}
			for k, nv := range v {
				if "id" == k || "typeID" == k {
					item[k], err = strconv.ParseUint(string(nv), 10, 32)
				} else {
					item[k] = string(nv)
				}
			}
			conditionArr = append(conditionArr, item)
		}
		res.Code = "0"
		res.Data = conditionArr
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
获取所有条件类型信息的接口
*/
func (ser *AllService) GetAllConditionTypeInfo(ctx iris.Context) {
	result, err := ser.SQL.Query("select * from `ConditionType` order by `id` DESC")
	res := &m.DataResponse{}
	if nil != err {
		res.Code = "-1"
		res.Msg = fmt.Sprintf("数据读取错误,%v", err)
	} else {
		res.Code = "0"
		arr := []map[string]interface{}{}
		for _, v := range result {
			item := map[string]interface{}{}
			for k, nv := range v {
				if "id" == k {
					item[k], err = strconv.ParseUint(string(nv), 10, 32)
				} else {
					item[k] = string(nv)
				}
			}
			arr = append(arr, item)
		}
		tmpConditionTypeInfoLock := ser.tmpConditionTypeInfoLock
		tmpConditionTypeInfoLock.Lock()
		ser.tmpConditionTypeInfoArr = arr
		tmpConditionTypeInfoLock.Unlock()
		res.Data = arr
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
获取条件类型的英文名称
*/
func (ser *AllService) getConditionEnName(ctypeId uint64) string {
	tmpConditionTypeInfoLock := ser.tmpConditionTypeInfoLock
	tmpConditionTypeInfoLock.Lock()
	defer tmpConditionTypeInfoLock.Unlock()
	arr := ser.tmpConditionTypeInfoArr
	result := ""
	for _, conditionTypeInfo := range arr {
		if id, isok := conditionTypeInfo["id"].(uint64); isok == true {
			if id == ctypeId {
				result, _ = conditionTypeInfo["enName"].(string)
				break
			}
		}
	}
	return result
}

/**
添加条件的类型
*/
func (ser *AllService) AddConditionType(ctx iris.Context) {
	zn := ctx.URLParam("zn")
	en := ctx.URLParam("en")
	des := ctx.URLParam("des")
	res := &m.CurrentResponse{}
	if "" != zn && "" != en {
		//执行插入
		_, err := ser.SQL.Exec(fmt.Sprintf("insert into `ConditionType`(`zhName`,`enName`,`des`) values('%s','%s','%s')", zn, en, des))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("添加条件类型失败,%v", err)
		} else {
			res.Code = "0"
			res.Msg = "条件类型添加成功"
		}
	} else {
		res.Code = "-1"
		res.Msg = "添加条件类型失败，请检查其请求参数"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
添加条件信息
*/
func (ser *AllService) AddCondition(ctx iris.Context) {
	cType, err1 := ctx.URLParamInt("cType")                                  //条件类型
	name := ctx.URLParam("name")                                             //条件名称
	val := ctx.URLParam("value")                                             //条件对应的值
	operator := ctx.URLParam("operator")                                     //操作符
	probability, err2 := strconv.ParseFloat(ctx.URLParam("probability"), 32) //条件生效几率
	des := ctx.URLParam("des")                                               //条件描述
	res := &m.CurrentResponse{}
	if nil == err1 && nil == err2 && "" != name && "" != val {
		//写入数据库
		gResult, err := ser.SQL.Exec(fmt.Sprintf("insert into `Condition`(`typeID`,`value`,`operator`,`name`,`probability`,`des`) values(%d,'%s','%s','%s',%f,'%s')", cType, val, operator, name, probability, des))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("新增条件失败,%v", err)
		} else {
			res.Code = "0"
			res.Msg = "新增条件成功"

			//更新条件配置数组
			tid, err := gResult.LastInsertId()
			go func() {
				_, isOk := <-ser.conditionChan
				if false == isOk {
					return
				}
				cItem := m.ConditionInfo{}
				if nil == err {
					cItem.Id = uint64(tid)
				} else {
					ser.conditionChan <- 1
					return
				}
				cItem.Typeid = uint64(cType)
				typeName := ser.getConditionEnName(cItem.Typeid)
				if typeName == "" {
					ser.conditionChan <- 1
					return
				}
				cItem.TypeName = typeName
				cItem.Value = val
				cItem.Operator = operator
				cItem.Probability = probability
				lastUpdate := fmt.Sprintf("%v", time.Now().Unix())
				cItem.LastUpdate = lastUpdate
				ser.conditionArr = append(ser.conditionArr, cItem)
				ser.conditionChange = true
				ser.conditionChan <- 1
			}()
		}
	} else {
		res.Code = "-1"
		res.Msg = "条件添加失败，请检查参数"
	}

	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
删除条件
*/
func (ser *AllService) DeleteCondition(ctx iris.Context) {
	id, err := strconv.ParseUint(ctx.URLParam("id"), 10, 32)
	res := &m.CurrentResponse{}
	if nil != err {
		res.Code = "-1"
		res.Msg = "删除条件失败,请检查参数"
	} else {
		result, err := ser.SQL.Exec(fmt.Sprintf("delete from `Condition` where `id` = %d", id))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("删除条件失败,%v", err)
		} else if count, e := result.RowsAffected(); nil == e && count > 0 {
			res.Code = "0"
			res.Msg = "条件删除成功"

			go func() {
				_, isOk := <-ser.conditionChan
				if false == isOk {
					return
				}
				//更新条件配置数组
				gj := len(ser.conditionArr)
				for gi, gv := range ser.conditionArr {
					if gv.Id == id {
						ser.conditionArr = append(ser.conditionArr[0:gi], ser.conditionArr[gi+1:gj]...)
						ser.conditionChange = true
						break
					}
				}
				ser.conditionChan <- 1
			}()

		} else {
			res.Code = "-1"
			res.Msg = "条件删除失败，参数id对应的条件不存在或者条件信息未变更导致更新失败"
		}
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
更新条件信息
*/
func (ser *AllService) UpdateConditionInfo(ctx iris.Context) {
	id, err1 := strconv.ParseUint(ctx.URLParam("id"), 10, 32)
	cType, err2 := strconv.ParseUint(ctx.URLParam("cType"), 10, 32)
	name := ctx.URLParam("name")
	val := ctx.URLParam("value")
	operator := ctx.URLParam("operator")
	probability, err3 := strconv.ParseFloat(ctx.URLParam("probability"), 32)
	des := ctx.URLParam("des")
	lastUpdate := fmt.Sprintf("%v", time.Now().Unix())
	res := &m.CurrentResponse{}
	if nil == err1 && nil == err2 && nil == err3 && "" != name && "" != val {
		result, err := ser.SQL.Exec(fmt.Sprintf("update `Condition` set `typeID`=%d,`value`='%s',`operator`='%s',`name`='%s',`probability`=%f,`des`='%s',`lastUpdate`='%s' where `id`=%d", cType, val, operator, name, probability, des, lastUpdate, id))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("条件信息更新失败,%v", err)
		} else if count, e := result.RowsAffected(); nil == e && count > 0 {
			res.Code = "0"
			res.Msg = "条件信息更新成功"

			go func() {
				_, isOk := <-ser.conditionChan
				if false == isOk {
					return
				}
				//更新条件配置文件
				queryStr := fmt.Sprintf("select `Condition`.`id`,`Condition`.`typeID`,`ConditionType`.`enName` ,`Condition`.`value` ,`Condition`.`probability`,`Condition`.`operator` from   `Condition` LEFT JOIN   `ConditionType`   on   `Condition`.`typeID` = `ConditionType`.`id` where `Condition`.`id` = %d", id)
				res, err := ser.SQL.Query(queryStr)
				if nil == err && len(res) > 0 {
					subV := res[0]
					cItem := m.ConditionInfo{}
					cItem.Id, _ = strconv.ParseUint(string(subV["id"]), 10, 32)
					cItem.Typeid, _ = strconv.ParseUint(string(subV["typeID"]), 10, 32)
					cItem.TypeName = string(subV["enName"])
					cItem.Value = string(subV["value"])
					cItem.Operator = string(subV["operator"])
					cItem.Probability, _ = strconv.ParseFloat(string(subV["probability"]), 32)
					cItem.LastUpdate = lastUpdate
					for gi, gv := range ser.conditionArr {
						if gv.Id == cItem.Id {
							ser.conditionArr[gi] = cItem
							ser.conditionChange = true
							break
						}
					}
				}
				ser.conditionChan <- 1
			}()

		} else {
			res.Code = "-1"
			res.Msg = "条件信息更新失败，参数id对应的条件不存在"
		}
	} else {
		res.Code = "-1"
		res.Msg = "更新条件信息失败，请检查请求参数"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
将内容写入文件
*/
func (ser *AllService) writeToFile(templateContent string, basePath string) (filePath string, err error) {
	//将模板文件写入静态服务器
	exist := true
	fileNameUUID, e := uuid.NewV4()
	if nil != e {
		return "", errors.New("生成失败,请重试")
	} else {
		//判断目录是否存在，否则就创建
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			//目录不存在，则创建
			mkdirErr := os.MkdirAll(basePath, 0774)
			if nil != mkdirErr {
				return "", errors.New(fmt.Sprintf("生成失败,%s路径创建失败", basePath))
			}
		}

		filePath = fmt.Sprintf("%s%v.json", basePath, fileNameUUID) //生成文件名
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			exist = false //判断文件是否存在
		}
		var f *os.File
		var fe error
		if exist {
			//如果文件存在则更新内容
			f, fe = os.OpenFile(filePath, os.O_WRONLY, 0774)
		} else {
			//创建文件，写入内容
			f, fe = os.Create(filePath)
		}
		if nil != fe {
			return "", errors.New(fmt.Sprintf("生成失败，原因是文件写入失败,%v", fe))
		} else {
			defer f.Close()                          //关闭文件
			_, fe = f.Write([]byte(templateContent)) //写入数据
			if nil != fe {
				return "", errors.New(fmt.Sprintf("生成失败，原因是文件写入失败,%v", fe))
			} else {
				//文件写入成功后， 更新数据库信息
				return filePath, nil
			}
		}
	}
}

/**
新增策略组
*/
func (ser *AllService) AddStrategyCategroy(ctx iris.Context) {
	name := ctx.PostValue("name")
	des := ctx.PostValue("des")
	templateContent := ctx.PostValue("templateContent")
	res := &m.CurrentResponse{}
	if "" != name && "" != templateContent {
		//校验templateContent是否是JSON内容
		var tmpJson []map[string]interface{}
		jsonErr := json.Unmarshal([]byte(templateContent), &tmpJson)
		if nil != jsonErr {
			res.Code = "-1"
			res.Msg = "策略模板生成失败,原因：参数templateContent对应的json内容格式无效"
		} else {
			//将模板文件写入静态服务器
			filePath, e := ser.writeToFile(templateContent, "./static/")
			if nil != e {
				res.Code = "-1"
				res.Msg = fmt.Sprintf("策略模板%v", e)
			} else {
				//文件写入成功后， 更新数据库信息
				_, err := ser.SQL.Exec(fmt.Sprintf("insert into `StrategyCategroy`(`name`,`des`,`baseTemplatePath`) values('%s','%s','%s')", name, des, filePath))
				if nil != err {
					res.Code = "-1"
					res.Msg = fmt.Sprintf("新建策略组失败，%v", err)
				} else {
					res.Code = "0"
					res.Msg = "新建策略组成功"
				}
			}
		}
	} else {
		res.Code = "-1"
		res.Msg = "新增策略组失败，请检查请求参数"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
获取所有的策略类型
*/
func (ser *AllService) GetAllStrategyCategroyInfo(ctx iris.Context) {
	queryStr := "select * from `StrategyCategroy`"
	res := &m.DataResponse{}
	resultList, err := ser.SQL.Query(queryStr)
	if err == nil {
		res.Code = "0"
		dataArr := []map[string]interface{}{}
		for _, v := range resultList {
			item := map[string]interface{}{}
			for k, nv := range v {
				if k == "id" {
					item[k], err = strconv.ParseUint(string(nv), 10, 32)
				} else {
					item[k] = string(nv)
				}
			}
			dataArr = append(dataArr, item)
		}
		res.Data = dataArr
	} else {
		res.Code = "-1"
		res.Msg = fmt.Sprintf("查询策略类别信息失败,%s", err.Error())
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
编辑策略组信息
*/
func (ser *AllService) UpdateStrategyCategroy(ctx iris.Context) {
	// id,err1 := strconv.ParseUint(ctx.URLParam("id"),10,32);
	// name := ctx.URLParam("name");
	// des := ctx.URLParam("des");
	// templateContent := ctx.URLParam("templateContent");
	id, err1 := strconv.ParseUint(ctx.PostValue("id"), 10, 32)
	des := ctx.PostValue("des")
	templateContent := ctx.PostValue("templateContent")
	res := &m.CurrentResponse{}
	if nil == err1 && "" != des && "" != templateContent {
		//校验模板格式
		var jsonObj []map[string]interface{}
		jsonErr := json.Unmarshal([]byte(templateContent), &jsonObj)
		if nil != jsonErr {
			res.Code = "-1"
			res.Msg = "templateContent内容不是JSON"
		} else {
			//删除之前的模板文件
			tmpSel, err := ser.SQL.Query(fmt.Sprintf("select `baseTemplatePath` from `StrategyCategroy` where `id` = %d", id))
			if nil != err {
				res.Code = "-1"
				res.Msg = "更新策略组信息失败，未能找到id对应的策略组"
			} else {
				//生成策略组的模板文件
				filePath, e := ser.writeToFile(templateContent, "./static/")
				if nil != e {
					res.Code = "-1"
					res.Msg = fmt.Sprintf("策略模板%v", e)
				} else {
					//写入数据库
					_, e := ser.SQL.Exec(fmt.Sprintf("update `StrategyCategroy` set `des`='%s',`baseTemplatePath`='%s' where `id`=%d", des, filePath, id))
					if nil != e {
						res.Code = "-1"
						res.Msg = fmt.Sprintf("更新策略组信息失败,%s", e.Error())
					} else {
						res.Code = "0"
						res.Msg = "更新策略组信息，成功"
						//删除旧的无用的策略模板文件
						for _, v := range tmpSel {
							baseTemplatePath := string(v["baseTemplatePath"])
							os.Remove(baseTemplatePath) //删除文件
						}
					}
				}
			}
		}
	} else {
		res.Code = "-1"
		res.Msg = "更新策略组信息失败，请检查参数"
	}
	resBytes, e := json.Marshal(res)
	if nil != e {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
删除策略组
*/
func (ser *AllService) DeleteStrategyCategroy(ctx iris.Context) {
	id, err := strconv.ParseUint(ctx.URLParam("id"), 10, 32)
	res := &m.CurrentResponse{}
	if nil != err {
		res.Code = "-1"
		res.Msg = "删除策略组失败，请检查参数"
	} else {
		filePathArr := []string{} //文件路径集合,用于删除
		//查询策略分类表中的记录
		result, e := ser.SQL.Query(fmt.Sprintf("select `baseTemplatePath` from `StrategyCategroy` where `id` = %d", id))
		if nil == e {
			for _, v := range result {
				filePathArr = append(filePathArr, string(v["baseTemplatePath"])) //填充文件路径组
			}
		}
		//遍历 策略记录表中的记录
		result, e = ser.SQL.Query(fmt.Sprintf("select `valuePath` from `Strategy` where `sid` = %d", id))
		if nil == e {
			for _, v := range result {
				filePathArr = append(filePathArr, string(v["valuePath"]))
			}
		}
		//删除策略类别表中的记录
		_, e = ser.SQL.Exec(fmt.Sprintf("delete from `StrategyCategroy` where `id` = %d", id))
		if nil != e {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("策略模板删除失败，%s", e.Error())
		} else {
			//删除策略表中的所有记录
			_, e = ser.SQL.Exec(fmt.Sprintf("delete from `Strategy` where `sid` = %d", id))
			if nil != e {
				res.Code = "-1"
				res.Msg = fmt.Sprintf("策略模板删除失败，%s", e.Error())
			} else {
				res.Code = "0"
				res.Msg = "策略组信息删除成功"
			}
			//删除所有无用的文件
			for _, v := range filePathArr {
				os.Remove(v)
			}
		}
	}

	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
新增策略
*/
func (ser *AllService) AddStrategy(ctx iris.Context) {
	// sid,err1 := strconv.ParseUint(ctx.URLParam("sid"),10,32);//策略组ID
	// strategyContext := ctx.URLParam("strategyContext");//策略内容
	// expire,err2 := strconv.ParseUint(ctx.URLParam("expire"),10,32);//过期时间戳
	// isEnabled,err3:= strconv.ParseUint(ctx.URLParam("enabled"),10,32);//是否为开启状态
	// name := ctx.URLParam("name");//策略名称
	sid, err1 := strconv.ParseUint(ctx.PostValue("sid"), 10, 32)           //策略组ID
	strategyContext := ctx.PostValue("strategyContext")                    //策略内容
	expire, err2 := strconv.ParseUint(ctx.PostValue("expire"), 10, 32)     //过期时间戳
	isEnabled, err3 := strconv.ParseUint(ctx.PostValue("enabled"), 10, 32) //是否为开启状态
	name := ctx.PostValue("name")                                          //策略名称
	res := &m.CurrentResponse{}
	if nil == err3 && nil == err1 && nil == err2 && "" != strategyContext && "" != name {
		//校验strategyContext是否为json
		var jsonObj map[string]interface{}
		err := json.Unmarshal([]byte(strategyContext), &jsonObj)
		if nil != err {
			res.Code = "-1"
			res.Msg = "添加策略失败,strategyContext对应的内容不是json"
		} else {
			//将策略内容写入文件
			filePath, err := ser.writeToFile(strategyContext, "./static/strategy/")
			if nil != err {
				res.Code = "-1"
				res.Msg = fmt.Sprintf("策略文件创建失败%v", err)
			} else {
				//写入数据库
				sqlRes, err := ser.SQL.Exec(fmt.Sprintf("insert into `Strategy`(`name`,`expireDate`,`enabled`,`valuePath`,`sid`) values('%v',%d,%d,'%v',%d)", name, expire, isEnabled, filePath, sid))
				if nil != err {
					res.Code = "-1"
					res.Msg = fmt.Sprintf("创建策略失败,%s", err.Error())
				} else if count, err := sqlRes.RowsAffected(); nil == err && count > 0 {
					res.Code = "0"
					res.Msg = "创建策略成功"
					//更新策略配置
					go func() {
						_, isOk := <-ser.strategyChan
						if false == isOk {
							return
						}
						tmpID, err := sqlRes.LastInsertId()
						if nil == err {
							querlStr := fmt.Sprintf("select `Strategy`.`id`,`Strategy`.`sid`,`Strategy`.`conditionGroup`,`Strategy`.`valuePath`,`Strategy`.`enabled`,`Strategy`.`expireDate`,`Strategy`.`lastUpdate`,`Strategy`.`name`,`StrategyCategroy`.`name` as `sname` from `Strategy` left join `StrategyCategroy` on `Strategy`.`sid` = `StrategyCategroy`.`id` where `Strategy`.`id` = %d", tmpID)
							tmpres, err := ser.SQL.Query(querlStr)
							if nil == err && len(tmpres) > 0 {
								item := tmpres[0]
								si := m.StrategyInfo{}
								si.Id, err = strconv.ParseUint(string(item["id"]), 10, 32)
								si.Sid, err = strconv.ParseUint(string(item["sid"]), 10, 32)
								si.Cgroup = string(item["conditionGroup"])
								si.ValuePath = string(item["valuePath"])
								si.Enabled, err = strconv.ParseUint(string(item["enabled"]), 10, 32)
								si.ExpireDate, err = strconv.ParseUint(string(item["expireDate"]), 10, 32)
								si.Name = string(item["name"])
								si.Sname = string(item["sname"])
								lastUpdate := fmt.Sprintf("%v", time.Now().Unix())
								si.LastUpdate = lastUpdate
								ser.strategyArr = append(ser.strategyArr, si)
								ser.strategyChange = true
							}
						}
						ser.strategyChan <- 1
					}()
				} else {
					res.Code = "-1"
					res.Msg = "创建策略失败,数据库插入失败"
				}
			}
		}
	} else {
		res.Code = "-1"
		res.Msg = "添加策略失败，请检查参数"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
为策略编辑匹配条件
*/
func (ser *AllService) EditConditionForStrategy(ctx iris.Context) {
	id, err1 := strconv.ParseUint(ctx.URLParam("id"), 10, 32) //策略id
	conditionStr := ctx.URLParam("conditionGroup")            //条件id的组合字符传
	lastUpdate := fmt.Sprintf("%v", time.Now().Unix())
	res := &m.CurrentResponse{}
	if nil == err1 {
		sqlRes, err := ser.SQL.Exec(fmt.Sprintf("update `Strategy` set `conditionGroup`='%s',`lastUpdate`='%s' where `id`=%d", conditionStr, lastUpdate, id))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("策略条件变更失败,%s", err.Error())
		} else if count, err := sqlRes.RowsAffected(); nil == err && count > 0 {
			res.Code = "0"
			res.Msg = "策略条件变更成功"
			//更新策略配置
			go func() {
				_, isOk := <-ser.strategyChan
				if false == isOk {
					return
				}
				for gi, gv := range ser.strategyArr {
					if gv.Id == id {
						gv.Cgroup = conditionStr
						gv.LastUpdate = lastUpdate
						ser.strategyArr[gi] = gv
						ser.strategyChange = true
						break
					}
				}
				ser.strategyChan <- 1
			}()
		} else {
			res.Code = "0"
			res.Msg = fmt.Sprintf("策略条件变更成功")
		}
	} else {
		res.Code = "-1"
		res.Msg = "为策略附加匹配条件，失败。请检查参数"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
根据策略ID，查询该策略对应的匹配条件
*/
func (ser *AllService) GetConditionInfoByStrategyID(ctx iris.Context) {
	id, err1 := strconv.ParseUint(ctx.URLParam("id"), 10, 32)
	res := &m.DataResponse{}
	if nil == err1 {
		//开始查询
		result, err := ser.SQL.Query(fmt.Sprintf("select `conditionGroup` from `Strategy` where `id` = %d", id))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("查询失败,%s", err.Error())
		} else {
			cidArr := []string{}
			for _, v := range result {
				strArr := strings.Split(string(v["conditionGroup"]), ",")
				for _, sv := range strArr {
					if _, ce := strconv.ParseUint(sv, 10, 32); ce == nil {
						cidArr = append(cidArr, sv)
					}
				}
			}
			cidStr := strings.Join(cidArr, ",")
			result, err = ser.SQL.Query(fmt.Sprintf("select `Condition`.`id` as `ConditionID`,`Condition`.`value`,`Condition`.`operator`,`Condition`.`name` as `ConditionName`,`Condition`.`probability`,`ConditionType`.`enName`,`ConditionType`.`zhName` from `Condition` left join `ConditionType`  on  `Condition`.`typeID` = `ConditionType`.`id` where `Condition`.`id` in(%s)", cidStr))
			if nil != err {
				res.Code = "-1"
				res.Msg = fmt.Sprintf("查询村略对应的条件，失败。%s", err.Error())
			} else {
				resultData := []map[string]interface{}{}
				//遍历查询结果
				for _, v := range result {
					tmp := map[string]interface{}{}
					for key, nv := range v {
						if key == "id" {
							tmp[key], err = strconv.ParseUint(string(nv), 10, 32)
						} else if "probability" == key {
							tmp[key], err = strconv.ParseFloat(string(nv), 32)
						} else {
							tmp[key] = string(nv)
						}
					}
					resultData = append(resultData, tmp)
				}
				res.Data = resultData
				res.Code = "0"
			}
		}
	} else {
		res.Code = "-1"
		res.Msg = "查询策略对应的条件，失败。请检查请求参数"
	}

	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
获取指定策略组内的所有策略信息
*/
func (ser *AllService) GetStrategyByStrategyCategroyID(ctx iris.Context) {
	sid, err := strconv.ParseUint(ctx.URLParam("sid"), 10, 32) //策略组id
	res := &m.DataResponse{}
	if nil == err {
		result, err := ser.SQL.Query(fmt.Sprintf("select * from `Strategy` where `sid` = %d", sid))
		if nil != err {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("查询策略信息失败,%s", err.Error())
		} else {
			arr := []map[string]interface{}{}
			//遍历结果
			for _, v := range result {
				item := map[string]interface{}{}
				for key, nv := range v {
					if "id" == key || "sid" == key || "expireDate" == key || "enabled" == key {
						item[key], err = strconv.ParseUint(string(nv), 10, 32)
					} else {
						item[key] = string(nv)
					}
				}
				arr = append(arr, item)
			}
			res.Code = "0"
			res.Data = arr
		}
	} else {
		res.Code = "-1"
		res.Msg = "查询策略信息失败，请检查请求参数"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
更新策略信息
*/
func (ser *AllService) UpdateStrategyInfo(ctx iris.Context) {
	// id,err1 := strconv.ParseUint(ctx.URLParam("id"),10,32);//策略ID
	// strategyContext := ctx.URLParam("strategyContext");//策略内容
	// expire,err2 := strconv.ParseUint(ctx.URLParam("expire"),10,32);//过期时间戳
	// isEnabled,err3:= strconv.ParseUint(ctx.URLParam("enabled"),10,32);//是否为开启状态
	// name := ctx.URLParam("name");//策略名称
	id, err1 := strconv.ParseUint(ctx.PostValue("id"), 10, 32)             //策略ID
	strategyContext := ctx.PostValue("strategyContext")                    //策略内容
	expire, err2 := strconv.ParseUint(ctx.PostValue("expire"), 10, 32)     //过期时间戳
	isEnabled, err3 := strconv.ParseUint(ctx.PostValue("enabled"), 10, 32) //是否为开启状态
	name := ctx.PostValue("name")                                          //策略名称
	lastUpdate := fmt.Sprintf("%v", time.Now().Unix())
	res := &m.CurrentResponse{}
	if nil == err3 && nil == err1 && nil == err2 && "" != strategyContext && "" != name {
		//校验strategyContext是否为json
		var jsonObj map[string]interface{}
		err := json.Unmarshal([]byte(strategyContext), &jsonObj)
		if nil != err {
			res.Code = "-1"
			res.Msg = "更新策略信息失败,strategyContext对应的内容不是json"
		} else {
			//根据id查询原有的策略信息
			result, err := ser.SQL.Query(fmt.Sprintf("select `valuePath` from `Strategy` where `id` = %d", id))
			if nil != err || len(result) <= 0 {
				res.Code = "-1"
				res.Msg = "未能找到策略ID对应的信息"
			} else {
				path := string(result[0]["valuePath"])
				//将策略内容写入新的文件
				filePath, err := ser.writeToFile(strategyContext, "./static/strategy/")
				if nil != err {
					res.Code = "-1"
					res.Msg = fmt.Sprintf("策略文件更新失败%v", err)
				} else {
					//写入数据库
					sqlRes, err := ser.SQL.Exec(fmt.Sprintf("update `Strategy` set `name`='%s',`expireDate`=%d,`enabled`=%d,`valuePath`='%s',`lastUpdate`='%s' where `id` = %d", name, expire, isEnabled, filePath, lastUpdate, id))
					if nil != err {
						res.Code = "-1"
						res.Msg = fmt.Sprintf("更新策略失败,%s", err.Error())
					} else if count, err := sqlRes.RowsAffected(); nil == err && count > 0 {
						res.Code = "0"
						res.Msg = "更新策略成功"
						if "" != path {
							os.Remove(path) //移除旧的策略
						}
						//更新策略配置
						go func() {
							_, isOk := <-ser.strategyChan
							if false == isOk {
								return
							}
							for gi, gv := range ser.strategyArr {
								if gv.Id == id {
									gv.ValuePath = filePath
									gv.Enabled = isEnabled
									gv.ExpireDate = expire
									gv.Name = name
									gv.LastUpdate = lastUpdate
									ser.strategyArr[gi] = gv
									ser.strategyChange = true
									break
								}
							}
							ser.strategyChan <- 1
						}()
					} else {
						res.Code = "0"
						res.Msg = "更新策略成功"
						if "" != path {
							os.Remove(path) //移除旧的策略
						}
					}
				}
			}
		}
	} else {
		res.Code = "-1"
		res.Msg = "更新策略信息失败，请检查参数"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/**
根据策略ID删除指定策略
*/
func (ser *AllService) DeleteStrategyByID(ctx iris.Context) {
	id, err := strconv.ParseUint(ctx.URLParam("id"), 10, 32)
	res := &m.CurrentResponse{}
	if nil == err {
		tmp, err := ser.SQL.Query(fmt.Sprintf("select `valuePath` from `Strategy` where `id` = %d", id))
		if nil != err || len(tmp) <= 0 {
			res.Code = "-1"
			res.Msg = fmt.Sprintf("删除策略失败，未能查询到id对应的指定策略")
		} else {
			oldPath := string(tmp[0]["valuePath"]) //获取旧的策略文件路径
			_, err := ser.SQL.Exec(fmt.Sprintf("delete from `Strategy` where `id` = %d", id))
			if nil != err {
				res.Code = "-1"
				res.Msg = fmt.Sprintf("删除策略失败，%s", err.Error())
			} else {
				res.Code = "0"
				res.Msg = "策略删除成功"
				if "" != oldPath {
					os.Remove(oldPath) //删除旧策略文件
				}
				go func() {
					_, isOk := <-ser.strategyChan
					if false == isOk {
						return
					}
					for gi, gv := range ser.strategyArr {
						if gv.Id == id {
							ser.strategyArr = append(ser.strategyArr[0:gi], ser.strategyArr[gi+1:len(ser.strategyArr)]...)
							ser.strategyChange = true
							break
						}
					}
					ser.strategyChan <- 1
				}()
			}
		}
	} else {
		res.Code = "-1"
		res.Msg = "删除策略失败，请检查参数"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

/*
 使策略强制即时生效，即所有在线用户即时更新指定ID对应的策略
 **/
func (ser *AllService) ForceStrategyBeUseage(ctx iris.Context) {
	id, err := strconv.ParseUint(ctx.URLParam("id"), 10, 32)
	res := &m.CurrentResponse{}
	if nil == err {
		res.Code = "0"
		res.Msg = "策略已即时生效"
		go func() {
			_, isOk := <-ser.strategyChan
			if isOk == false {
				return
			}

			_, isOk2 := <-ser.conditionChan
			if isOk2 == false {
				return
			}
			if true == ser.strategyChange {
				ser._saveStrategyConfig()
				ser.strategyChange = false
			}
			if true == ser.conditionChange {
				ser._saveConditionConfig()
				ser.conditionChange = false
			}
			//通知长连接服务，强制策略生效
			req := &m.ForceStrategyBeUse_c2s{Cmd: 0x00FF0041, StrategyPath: ser.StrategyConfigPath, ConditionPath: ser.ConditionConfigPath, StrategyID: id}
			ser.Sock.Write(req)
			ser.SVCClient.PushMsg(req)
			ser.strategyChan <- 1
			ser.conditionChan <- 1
		}()

	} else {
		res.Code = "-1"
		res.Msg = "策略未能即时生效"
	}
	resBytes, err := json.Marshal(res)
	if nil != err {
		ctx.Write([]byte(""))
	} else {
		ctx.Write(resBytes)
	}
}

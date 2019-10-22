var interfaceDomain = "http://localhost:8080/"

/**
 * 后台登录
 * http://localhost:8080/front/signByAccount?ln=guominglong&pwd=543b10.
 * @param ln string 用户名
 * @param pwd string 用户登录密码
*/
function signByAccount(ln,pwd,callbackFunc){
    $.ajax({
        url:interfaceDomain + "signByAccount",
        data:{
            ln:ln,
            pwd:pwd,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 登出后台
 * 
*/
function signOut(callbackFunc){
    $.ajax({
        url:interfaceDomain + "signOut",
        data:{
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}


/**
 * 添加后台管理账号
 * http://localhost:8080/AddUser?ln=gml&pwd=543b10.&roleID=2
 * @param ln string 用户名
 * @param pwd string 用户登录密码
 * @param roleID int 角色类型ID
 * @param callbackFunc function(data,err)回调函数
*/
function addUser(ln,pwd,roleID,callbackFunc){
    $.ajax({
        url:interfaceDomain + "AddUser",
        data:{
            ln:ln,
            pwd:pwd,
            roleID:roleID,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 获取所有的用户角色类型
 * http://localhost:8080/GetAllRoleType
*/
function getAllRoleType(callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetAllRoleType",
        data:{
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 获取所有的权限
 * http://localhost:8080/GetAllAuth
 * 
*/
function getAllAuth(callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetAllAuth",
        data:{
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 获取后台管理账号总数
 * http://localhost:8080/GetUsersCount
*/
function getUserCount(callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetUsersCount",
        data:{
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * http://localhost:8080/GetAllUsers
 * 分页查询后台管理用户账号的信息
 * @param startPoint int 开始读取数据的索引
 * @param readCount int 读取条数
*/
function getAllUsers(startPoint,readCount,callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetAllUsers",
        data:{
            startPoint:startPoint,
            readCount:readCount,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 删除后台管理账号
 * http://localhost:8080/DeleteUser?uid=1
 * @param uid int 用户id
*/
function deleteUser(uid,callbackFunc){
    $.ajax({
        url:interfaceDomain + "DeleteUser",
        data:{
            uid:uid,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 更新后台管理账号的信息
 * http://localhost:8080/UpdateUserInfo?sn=guominglong2&sp=123456&rid=1&uid=3
 * @param sn string 用户登录名
 * @param sp string 用户密码
 * @param rid int 用户角色类型ID
 * @param uid int 用户id
*/
function updateUserInfo(uid,sn,sp,rid,callbackFunc){
    $.ajax({
        url:interfaceDomain + "UpdateUserInfo",
        data:{
            uid:uid,
            sn:sn,
            sp:sp,
            rid:rid,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 分页获取策略条件信息
 * http://localhost:8080/GetAllConditionInfo?startPoint=0&readCount=10
 * @param startPoint int 起始位置
 * @param readCount int 获取数量
*/
function getAllConditionInfo(startPoint,readCount,callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetAllConditionInfo",
        data:{
            startPoint:startPoint,
            readCount:readCount,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 获取策略条件总数
 * http://localhost:8080/GetConditionCount
*/
function getConditionCount(callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetConditionCount",
        data:{
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 获取所有的策略条件类型
 * http://localhost:8080/GetAllConditionTypeInfo
 * 
*/
function getAllConditionTypeInfo(callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetAllConditionTypeInfo",
        data:{
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}


/**
 * 新增策略条件接口
 * http://localhost:8080/AddConditionType?zn=角色类型&en=RoleType
 * @param zn string 中文名
 * @param en string 英文名
 * @param des string 描述
*/
function addConditionType(zn,en,des,callbackFunc){
    $.ajax({
        url:interfaceDomain + "AddConditionType",
        data:{
            zn:zn,
            en:en,
            des:des,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 添加策略条件
 * http://localhost:8080/AddCondition?cType=1&name=国内用户&value=13&probability=0.5&des=这是描述
 * @param cType int 条件类型
 * @param name string 条件名称
 * @param value string 条件值
 * @param probability number 生效几率 0-1之间任意小数
 * @param des string 条件描述
*/
function addCondition(cType,name,value,probability,des,callbackFunc){
    $.ajax({
        url:interfaceDomain + "AddCondition",
        data:{
            cType:cType,
            name:name,
            value:value,
            probability:probability,
            des:des,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 更新策略条件信息
 * http://localhost:8080/UpdateConditionInfo?cType=2&name=%E5%9B%BD%E5%86%85%E7%94%A8%E6%88%B72&value=2&probability=2&id=8
 * @param cType int 条件类型
 * @param name string 条件名称
 * @param value string 条件值
 * @param probability number 生效几率 0-1之间任意小数
 * @param des string 条件描述
 * @param id int 条件ID
*/
function updateConditionInfo(cType,name,value,probability,des,id,callbackFunc){
    $.ajax({
        url:interfaceDomain + "UpdateConditionInfo",
        data:{
            cType:cType,
            name:name,
            value:value,
            probability:probability,
            des:des,
            id:id,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 删除策略条件
 * http://localhost:8080/DeleteCondition?id=3
 * @param id int 条件ID
*/
function deleteCondition(id,callbackFunc){
    $.ajax({
        url:interfaceDomain + "DeleteCondition",
        data:{
            id:id,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 添加策略类别
 * http://localhost:8080/AddStrategyCategroy?name=测试策略&des=这是测试这是测试&templateContent={"name":"郭明龙"}
 * @param name string 策略类别名称
 * @param des string 描述
 * @param templateContent string 策略模板内容
*/
function addStrategyCategroy(name,des,templateContent,callbackFunc){
    $.ajax({
        url:interfaceDomain + "AddStrategyCategroy",
        data:{
            name:name,
            des:des,
            templateContent:templateContent,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"POST",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 获取所有的 策略类别信息
*/
function getAllStrategyCategroy(callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetAllStrategyCategroyInfo",
        data:{
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 更新策略类别信息
 * http://localhost:8080/UpdateStrategyCategroy?name=测试策略4&des=这是444&templateContent={"name":"郭明龙4"}&id=3
 * @param name string 策略类别名称
 * @param des string 描述
 * @param templateContent string 策略模板内容
 * @param id int 策略类别id
*/
function updateStrategyCategroy(name,des,templateContent,id,callbackFunc){
    $.ajax({
        url:interfaceDomain + "UpdateStrategyCategroy",
        data:{
            name:name,
            des:des,
            templateContent:templateContent,
            id:id,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"POST",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 
 * 删除策略类别
 * http://localhost:8080/DeleteStrategyCategroy?id=3
 * @param id int 策略类别ID
*/
function deleteStrategyCategroy(id,callbackFunc){
    $.ajax({
        url:interfaceDomain + "DeleteStrategyCategroy",
        data:{
            id:id,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}


/**
 * 新增策略
 * http://localhost:8080/AddStrategy?sid=1&strategyContext={"name":"郭明龙"}&expire=123456&enabled=1&name=测试策略
 * @param sid int 策略类型ID
 * @param strategyContext string 策略内容
 * @param expire number 过期时间
 * @param enabled int 是否启用 0或者1
 * @param name string 策略名
 */
function addStrategy(sid,strategyContext,expire,enabled,name,callbackFunc){
    $.ajax({
        url:interfaceDomain + "AddStrategy",
        data:{
            sid:sid,
            strategyContext:strategyContext,
            expire:expire,
            enabled:enabled,
            name:name,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"POST",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 变更指定策略对应的条件
 * http://localhost:8080/EditConditionForStrategy?id=2&conditionGroup=1,2,4
 * @param id int 策略id
 * @param conditionGroup string 条件组
*/
function editConditionForStrategy(id,conditionGroup,callbackFunc){
    $.ajax({
        url:interfaceDomain + "EditConditionForStrategy",
        data:{
            id:id,
            conditionGroup:conditionGroup,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 获取指定策略对应的条件信息
 * http://localhost:8080/SelectConditionInfoByStrategyID?id=2
 * @param id int 策略ID
*/
function getConditionInfoByStrategyID(id,callbackFunc){
    $.ajax({
        url:interfaceDomain + "SelectConditionInfoByStrategyID",
        data:{
            id:id,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 获取指定策略类别下的所有策略
 * http://localhost:8080/GetStrategyByStrategyCategroyID?sid=2
 * @param sid int 策略类别ID
*/
function getStrategyByStrategyCategroyID(sid,callbackFunc){
    $.ajax({
        url:interfaceDomain + "GetStrategyByStrategyCategroyID",
        data:{
            sid:sid,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 更新指定的策略信息
 * http://localhost:8080/UpdateStrategyInfo?id=4&strategyContext={"name":"郭明龙"}&expire=111111&enabled=0&name=策略我的
 * @param id int 策略ID
 * @param strategyContext string 策略内容
 * @param expire number 过期时间
 * @param enabled int 是否启用 0或者1
 * @param name string 策略名
*/
function updateStrategyInfo(id,strategyContext,expire,enabled,name,callbackFunc){
    $.ajax({
        url:interfaceDomain + "UpdateStrategyInfo",
        data:{
            id:id,
            strategyContext:strategyContext,
            expire:expire,
            enabled:enabled,
            name:name,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"POST",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 删除指定ID对应的策略
 * http://localhost:8080/DeleteStrategyByID?id=4
 * @param id int 策略ID
*/
function deleteStrategyByID(id,callbackFunc){
    $.ajax({
        url:interfaceDomain + "DeleteStrategyByID",
        data:{
            id:id,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}

/**
 * 强制指定ID对应的策略即时生效
 * http://localhost:8080/ForceStrategyBeUseage?id=4
 * @param id int 策略ID
*/
function forceStrategyBeUseage(id,callbackFunc){
    $.ajax({
        url:interfaceDomain + "ForceStrategyBeUseage",
        data:{
            id:id,
            t:new Date().valueOf()
        },
        dataType:"json",
        method:"GET",
        success:function(data){
            callbackFunc(data,null);
        },
        error:function(err){
            callbackFunc(null,err);
        }
    })
}
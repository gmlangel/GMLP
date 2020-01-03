package proxy

import (
	"time"

	"../front"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

var (
	/*sql相关*/
	sqlType    = "mysql"
	sqlFullURL = "gmlmaster:123456@tcp(39.106.135.11:32306)/RemoteDecitionDB?charset=utf8"
	sqlpro     *SQLProxy
	rdsock     *front.RDSocket
	svcPushCli *front.SVCPushClient
	/*sessionManager相关*/
	sm *sessions.Sessions = sessions.New(
		sessions.Config{
			Cookie:  "cookieNameForSessionIDAndUserInfo",
			Expires: time.Duration(30) * time.Minute})
)

type WebServiceProxy struct {
	app *iris.Application
}

/*
webService管理
*/
func NewWebService(_app *iris.Application) *WebServiceProxy {
	pro := &WebServiceProxy{app: _app}
	pro.init()
	return pro
}

func (webs *WebServiceProxy) init() {
	//连接数据库
	sqlpro = NewSQL(sqlType, sqlFullURL)
	go sqlpro.Start()

	//连接远程长连接服务
	rdsock = front.NewRDSocket()
	go rdsock.Start("tcp", "0.0.0.0:63333")

	//创建svc push推送服务客户端
	svcPushCli = &front.SVCPushClient{}
	//测试用
	svcPushCli.PushMsg(map[string]string{"name": "ggg"})
}

func (webs *WebServiceProxy) Start() {
	//开启前端服务监听
	frontSignGroup := webs.app.Party("/front/sign") //前端需要验证签名的服务组
	f_loginProxy := front.LoginService{SqlPro: sqlpro, App: webs.app, SignGroup: &frontSignGroup, Sm: sm}
	frontSignGroup.Use(f_loginProxy.MW_CheckSinged) //添加登录校验
	f_loginProxy.Start()                            //启动login相关服务

	allser := &front.AllService{SQL: sqlpro, Sock: rdsock, SVCClient: svcPushCli}
	if sqlpro.IsConnected {
		allser.InitData() //直接初始化数据
	} else {
		//等待sql启动完成后，加载数据
		sqlpro.OnLinkComplete = func() {
			allser.InitData() //直接初始化数据
		}
	}

	webs.app.Get("AddUser", allser.AddUser)                                                 //新增后台管理账号
	webs.app.Get("GetAllRoleType", allser.GetAllRoleType)                                   //查询后台管理账号可选角色
	webs.app.Get("GetAllAuth", allser.GetAllAuth)                                           //查询后台角色的权限说明
	webs.app.Get("GetAllUsers", allser.GetAllUsers)                                         //查询后台角色的权限说明
	webs.app.Get("GetUsersCount", allser.GetUsersCount)                                     //获取后台用户总数
	webs.app.Get("DeleteUser", allser.DeleteUser)                                           //删除后台管理账号
	webs.app.Get("UpdateUserInfo", allser.UpdateUserInfo)                                   //更新管理账号的信息
	webs.app.Get("GetConditionCount", allser.GetConditionCount)                             //分页获取策略条件总数
	webs.app.Get("GetAllConditionInfo", allser.GetAllConditionInfo)                         //分页获取策略条件信息
	webs.app.Get("GetAllConditionTypeInfo", allser.GetAllConditionTypeInfo)                 //获取所有策略条件类型的信息
	webs.app.Get("AddConditionType", allser.AddConditionType)                               //添加条件类型
	webs.app.Get("AddCondition", allser.AddCondition)                                       //新增条件
	webs.app.Get("UpdateConditionInfo", allser.UpdateConditionInfo)                         //更新条件信息接口
	webs.app.Get("DeleteCondition", allser.DeleteCondition)                                 //删除策略条件接口
	webs.app.Post("AddStrategyCategroy", allser.AddStrategyCategroy)                        //新增策略组
	webs.app.Get("GetAllStrategyCategroyInfo", allser.GetAllStrategyCategroyInfo)           //获取所有的策略类别信息
	webs.app.Post("UpdateStrategyCategroy", allser.UpdateStrategyCategroy)                  //更新策略组
	webs.app.Get("DeleteStrategyCategroy", allser.DeleteStrategyCategroy)                   //删除策略组信息
	webs.app.Post("AddStrategy", allser.AddStrategy)                                        //新建策略
	webs.app.Get("EditConditionForStrategy", allser.EditConditionForStrategy)               //为策略添加匹配条件
	webs.app.Get("GetConditionInfoByStrategyID", allser.GetConditionInfoByStrategyID)       //查询指定策略对应的匹配条件
	webs.app.Get("GetStrategyByStrategyCategroyID", allser.GetStrategyByStrategyCategroyID) //根据策略组ID，获取对应的所有策略信息
	webs.app.Post("UpdateStrategyInfo", allser.UpdateStrategyInfo)                          //更新策略信息
	webs.app.Get("DeleteStrategyByID", allser.DeleteStrategyByID)                           //删除指定id对应的策略
	webs.app.Get("ForceStrategyBeUseage", allser.ForceStrategyBeUseage)                     //使策略强制即时生效，即所有在线用户即时更新指定ID对应的策略
	// //开启后端服务监听
	// //signServiceGroup2 := webs.app.Party("/rear/sign");
	// // rearProxy := rear.LoginService{};
	// // webs.app.Any("/",frontProxy.WelCome);

}

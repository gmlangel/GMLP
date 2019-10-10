package proxy;

import(
	"../front"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"time"
)

var(
	/*sql相关*/
	sqlType = "mysql";
	sqlFullURL = "gmlmaster:123456@tcp(39.106.135.11:32306)/RemoteDecitionDB?charset=utf8";
	sqlpro *SQLProxy;
	/*sessionManager相关*/
	sm *sessions.Sessions= sessions.New(
		sessions.Config{
			Cookie: "cookieNameForSessionIDAndUserInfo",
			Expires:time.Duration(30) * time.Minute})
)

type WebServiceProxy struct{
	app *iris.Application;
}

/*
webService管理
*/
func NewWebService(_app *iris.Application)(*WebServiceProxy){
	pro := &WebServiceProxy{app:_app};
	pro.init();
	return pro;
}



func (webs *WebServiceProxy)init(){
	//连接数据库
	sqlpro = NewSQL(sqlType,sqlFullURL);
	go sqlpro.Start();

}

func (webs *WebServiceProxy)Start(){
	//开启前端服务监听
	frontSignGroup := webs.app.Party("/front/sign");//前端需要验证签名的服务组
	f_loginProxy := front.LoginService{SqlPro:sqlpro,App:webs.app,SignGroup:&frontSignGroup,Sm:sm};
	frontSignGroup.Use(f_loginProxy.MW_CheckSinged);//添加登录校验
	f_loginProxy.Start();//启动login相关服务

	allser := &front.AllService{SQL:sqlpro}
	webs.app.Get("AddUser",allser.AddUser);//新增后台管理账号
	webs.app.Get("GetAllRoleType",allser.GetAllRoleType);//查询后台管理账号可选角色
	webs.app.Get("GetAllAuth",allser.GetAllAuth);//查询后台角色的权限说明
	webs.app.Get("GetAllUsers",allser.GetAllUsers);//查询后台角色的权限说明
	webs.app.Get("DeleteUser",allser.DeleteUser);//删除后台管理账号
	webs.app.Get("UpdateUserInfo",allser.UpdateUserInfo);//更新管理账号的信息
	webs.app.Get("GetAllConditionInfo",allser.GetAllConditionInfo);//分页获取策略条件信息
	webs.app.Get("GetAllConditionTypeInfo",allser.GetAllConditionTypeInfo);//获取所有策略条件类型的信息
	webs.app.Get("AddConditionType",allser.AddConditionType);//添加条件类型
	webs.app.Get("AddCondition",allser.AddCondition);//新增条件
	webs.app.Get("UpdateConditionInfo",allser.UpdateConditionInfo);//更新条件信息接口
	webs.app.Get("DeleteCondition",allser.DeleteCondition);//删除策略条件接口
	webs.app.Get("AddStrategyCategroy",allser.AddStrategyCategroy);//新增策略组
	// //开启后端服务监听
	// //signServiceGroup2 := webs.app.Party("/rear/sign");
	// // rearProxy := rear.LoginService{};
	// // webs.app.Any("/",frontProxy.WelCome);
}


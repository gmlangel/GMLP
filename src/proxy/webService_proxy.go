package proxy;

import(
	"../GMLWebService/front"
	_ "../GMLWebService/rear"
	"github.com/kataras/iris"
)

var(
	/*sql相关*/
	sqlType = "mysql";
	sqlFullURL = "gmlmaster:123456@tcp(39.106.135.11:32306)/GMLPlanDB?charset=utf8";
	sqlpro *SQLProxy;
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
	frontProxy := front.LoginService{SqlPro:sqlpro,App:webs.app};
	frontProxy.Start();
	//开启后端服务监听
	// rearProxy := rear.LoginService{};
	// webs.app.Any("/",frontProxy.WelCome);
}



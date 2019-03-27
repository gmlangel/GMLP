package proxy;

import(
	"fmt"
	"../GMLWebService/front"
	"../GMLWebService/rear"
	"github.com/kataras/iris"
)

var(
	/*sql相关*/
	sqlType = "mysql";
	sqlFullURL = "gmlmaster:123456@tcp(39.106.135.11:32306)/GMLPlanDB?charset=utf8";
	sqlpro SQLProxy;
)

type WebServiceProxy struct{
	app *iris.Application;
}

/*
webService管理
*/
func NewWebService(_app *iris.Application)WebServiceProxy{
	pro := WebServiceProxy{app:_app};
	pro.init();
	return pro;
}



func (webs *WebServiceProxy)init(){
	//连接数据库
	sqlpro = NewSQL(sqlType,sqlFullURL);
	go sqlpro.Start();

}

func (webs *WebServiceProxy)Start(){
	t1 := front.LoginService{};
	t2 := rear.LoginService{};
	t1.F();
	t2.F();
	webs.app.Any("/",welCome);
}


func welCome(ctx iris.Context){
	fmt.Println("欢迎使用GMLP");
	res,err := sqlpro.Query("select `uid` from `users`");
	if err == nil{
		fmt.Println("获取users表总数据为:",len(res),"条");
	}
	ctx.WriteString("<H1>欢迎使用GMLP</H1>")
}
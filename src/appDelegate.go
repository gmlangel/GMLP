package main

import(
	"fmt"
	"github.com/kataras/iris"
	"./proxy"
)
var(
	/*sql相关*/
	sqlType = "mysql";
	sqlFullURL = "gmlmaster:123456@tcp(39.106.135.11:32306)/GMLPlanDB?charset=utf8";
	sqlpro proxy.SQLProxy;
	/*静态资源管理相关*/
	staticPro proxy.StaticManagerProxy;
	/*动态资源管理相关*/
	/*webservice相关*/
)

func main(){
	fmt.Println("ok");
	//连接数据库
	sqlpro = proxy.NewSQL(sqlType,sqlFullURL);
	go sqlpro.Start();
	//初始化服务器
	app := iris.New();
	//设置服务器图标
	app.Favicon("./GMLClient/static/myico.ico");
	//启动静态服务
	staticPro = proxy.NewStaticManager(app);
	//staticPro.Start();
	//启动动态服务

	//启动接口服务
	app.Any("/",welCome);
	app.Run(iris.Addr("0.0.0.0:8080"));
}

func welCome(ctx iris.Context){
	fmt.Println("欢迎使用GMLP");
	res,err := sqlpro.Query("select `uid` from `users`");
	if err == nil{
		fmt.Println("获取users表总数据为:",len(res),"条");
	}
	ctx.WriteString("<H1>欢迎使用GMLP</H1>")
}
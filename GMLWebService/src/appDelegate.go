package main

import(
	"fmt"
	"github.com/kataras/iris"
	"./proxy"
)
var(
	sqlType = "mysql";
	sqlFullURL = "gmlmaster:123456@tcp(39.106.135.11:32306)/GMLPlanDB?charset=utf8";
	sqlpro proxy.SQLProxy;
)

func main(){
	fmt.Println("ok");
	//连接数据库
	sqlpro = proxy.New(sqlType,sqlFullURL);
	sqlpro.Start();
	//启动webservice
	app := iris.New();
	app.Any("/",welCome);
	app.Run(iris.Addr("0.0.0.0:7777"));
}

func welCome(ctx iris.Context){
	fmt.Println("欢迎使用GMLP");
	res,err := sqlpro.Query("select `uid` from `users`");
	if err == nil{
		fmt.Println("获取users表总数据为:",len(res),"条");
	}
}
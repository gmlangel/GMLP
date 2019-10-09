package main

import(
	"fmt"
	"github.com/kataras/iris"
	"./proxy"
)
var(
	/*静态资源管理相关*/
	staticPro *proxy.StaticManagerProxy;
	/*动态资源管理相关*/
	/*webservice相关*/

)

func main(){
	fmt.Println("ok");
	
	//初始化服务器
	app := iris.New();
	//启动静态服务
	staticPro = proxy.NewStaticManager(app);
	staticPro.Start();
	//启动动态服务
	//启动接口服务
	webs := proxy.NewWebService(app);
	webs.Start();
	//绑定服务器端口监听
	app.Run(iris.Addr("0.0.0.0:8080"));
}

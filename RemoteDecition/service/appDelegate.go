package main
import(
	"fmt"
	"github.com/kataras/iris"
	"./proxy"
)


func main(){
	fmt.Println("RemoteDecition远程决策服务开始启动....");
	//初始化服务器
	app := iris.New();
	//启动静态服务
	staticPro := proxy.NewStaticManager(app);
	staticPro.Start();
	//启动动态服务
	//启动接口服务
	webs := proxy.NewWebService(app);
	webs.Start();
	//绑定服务器端口监听
	app.Run(iris.Addr("0.0.0.0:8080"));
}
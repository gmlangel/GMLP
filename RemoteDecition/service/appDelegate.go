package main

import (
	"fmt"

	"./proxy"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func main() {
	debug := true //调试模式
	fmt.Println("RemoteDecition远程决策服务v 1.0开始启动....")
	//初始化服务器
	app := iris.New()
	//修改跨域访问限制
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts."*"代表允许所有域访问，这是一个数组，可以添加多个域名
		AllowCredentials: true,
		AllowedMethods:   []string{"HEAD", "GET", "POST"},
		AllowedHeaders:   []string{"accept, content-type", "Access-Control-Allow-Origin"}})
	app.Use(crs)
	//启动静态服务
	staticPro := proxy.NewStaticManager(app)
	staticPro.Start()
	//启动动态服务

	//启动接口服务
	webs := proxy.NewWebService(app)
	webs.Start()
	//绑定服务器端口监听
	if debug == true {
		app.Run(iris.Addr("0.0.0.0:8080"))
	} else {
		app.Run(iris.TLS("0.0.0.0:8080", "./1540854920368.pem", "./1540854920368.key"))
	}

}

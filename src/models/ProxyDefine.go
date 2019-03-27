package models;


import(
	"github.com/kataras/iris"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	"time"
)
/**
所有委托相关的模型，请定义到这里
*/

type SQLProxy struct{
	SQLType string //数据库类型
	DBFullURL string //"gmlmaster:123456@tcp(39.106.135.11:32306)/GMLPlanDB?charset=utf8"
	MaxIdleConns int //设置连接池的空闲数大小
	MaxOpenConns int //设置最大打开连接数
	LogLevel core.LogLevel//日志级别
	SqlHeartOffset time.Duration//心跳间隔
	sqlEngine *xorm.Engine//数据库引擎
	isConnected bool//是否已经连接成功
}

type WebServiceProxy struct{
	app *iris.Application;
}

type StaticManagerProxy struct{
	app *iris.Application;
}
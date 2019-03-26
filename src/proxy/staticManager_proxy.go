package proxy;

import(
	"github.com/kataras/iris"
)
/**
静态资源管理
*/
func NewStaticManager(_app *iris.Application) StaticManagerProxy{
	sm := StaticManagerProxy{app:_app};
	return sm;
}

type StaticManagerProxy struct{
	app *iris.Application;
}

/**
启动静态服务
*/
func (sm *StaticManagerProxy)Start(){
	sm.app.StaticWeb("/static","./GMLClient/static");
}
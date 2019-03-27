package front;

import(
	"fmt"
	m "../../models"
	"github.com/kataras/iris"
)

/**
前端 登录相关服务
*/
type LoginService struct{
	SqlPro m.SQLInterface;
}
func (ls *LoginService)F(){
	fmt.Println("我是front");
}

func (ls *LoginService)WelCome(ctx iris.Context){
	fmt.Println("欢迎使用GMLP");
	res,err := ls.SqlPro.Query("select `uid` from `users`");
	if err == nil{
		fmt.Println("获取users表总数据为:",len(res),"条");
	}
	ctx.WriteString("<H1>欢迎使用GMLP</H1>")
}
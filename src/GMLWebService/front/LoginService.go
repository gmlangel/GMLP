package front;

import(
	m "../../models"
	"github.com/kataras/iris"
	"fmt"
	"../../tools"
)

/**
前端 登录相关服务
*/
type LoginService struct{
	SqlPro m.SQLInterface;
	App *iris.Application;
}

/**
开始监听前端服务接口的调用
*/
func (ls LoginService)Start(){
	ls.App.Get("/front/singByAccount",ls.sign);
	ls.App.Get("/front/signOut",ls.signOut);
	ls.App.Get("/front/registerAccount",ls.registerAccount);
	ls.App.Get("/front/changepwd",ls.changepwd);
	ls.App.Get("/front/getVerificationCode",ls.getVerificationCode);
}

/**
通过账号密码登录
*/
func (ls LoginService)sign(ctx iris.Context){
	ctx.Write([]byte("登录成功"));
}

/**
登出
*/
func (ls LoginService)signOut(ctx iris.Context){
	ctx.Write([]byte("登出成功"));
}

/**
注册账号
*/
func (ls LoginService)registerAccount(ctx iris.Context){
	 //name := ctx.URLParam("ln");//gv(registerForm["ln"]);
	 
	 fmt.Println(tools.Pack(ctx).GetStr("ln"));
// 	pwd := gv(registerForm["pwd"]);
// 	resStruct := models.CurrentResponse{Code:resFaild,Msg:"注册失败"};
// 	//检查同名用户是否存在
// 	queryRes,queryErr := SQLQuery(fmt.Sprintf("SELECT `ln` FROM `BusinessUsers` WHERE `ln` = '%s';",name));
// 	if queryErr == nil{
// 		if len(queryRes) > 0{
// 			//数据中存在相同账号，注册失败
// 			resStruct.Code = resFaild;
// 			resStruct.Msg = "账号已存在，注册失败";
// 			fmt.Fprintln(response,structToJSONStr(resStruct));
// 			return;
// 		}
// 	}else{
// 		fmt.Fprintln(response,structToJSONStr(resStruct));
// 		return;
// 	}

// 	gmluuid,uuidErr:= uuid.NewV4();//生成唯一ID
// 	if uuidErr != nil{
// 		fmt.Println("uuid生成失败");
// 		//插入失败
// 		fmt.Fprintln(response,structToJSONStr(resStruct))
// 		return;
// 	}
// 	//向bussinessUsers插入数据
// 	res,err := SQLExec(fmt.Sprintf("insert `BusinessUsers`(`bid_str`,`ln`,`pwd`) values('%s','%s','%s');",gmluuid,name,pwd));
// 	if err == nil{
// 		if line,_ := res.RowsAffected() ;line > 0{
// 			//向BusinessInfo插入数据
// 			res2,err2 := SQLExec(fmt.Sprintf("insert `BusinessInfo`(`bid_str`,`des`,`bname`) values('%s','%s','%s');",gmluuid,"",""));
// 			if err2 == nil{
// 				if line2,_:=res2.RowsAffected();line2 >0{
// 					resStruct.Code = resOK;
// 					resStruct.Msg = "注册成功"
// 					//添加成功
// 					fmt.Fprintln(response,structToJSONStr(resStruct))
// 				}else{
// 					//插入失败
// 					fmt.Fprintln(response,structToJSONStr(resStruct))
// 				}
// 			}else{
// 				//插入失败
// 				fmt.Fprintln(response,structToJSONStr(resStruct))
// 			}
			
// 		}else{
// 			//插入失败
// 			fmt.Fprintln(response,structToJSONStr(resStruct))
// 		}
// 	}else{
// 		//插入失败
// 		fmt.Fprintln(response,structToJSONStr(resStruct))
// 	}
}

/**
修改密码
*/
func (ls LoginService)changepwd(ctx iris.Context){

}

/**
获取验证码
*/
func (ls LoginService)getVerificationCode(ctx iris.Context){

}

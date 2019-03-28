package front;

import(
	"../../models"
	"github.com/kataras/iris"
	"fmt"
	"../../tools"
	"github.com/satori/go.uuid"
)

var(
	resFaild = "0";
	resOK = "1";
	InvalideLoginName = "2";//账号不存在
	PWDFaild = "3";//密码错误
	InvalideMethod = "1001";//无效的method请求
)

/**
前端 登录相关服务
*/
type LoginService struct{
	SqlPro models.SQLInterface;
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
	sc := tools.Pack(ctx);
	name := sc.GetStr("ln");
	pwd := sc.GetStr("pwd");
 	resStruct := models.CurrentResponse{Code:resFaild,Msg:"注册失败"};
 	//检查同名用户是否存在
 	queryRes,queryErr := ls.SqlPro.Query(fmt.Sprintf("SELECT `ln` FROM `BusinessUsers` WHERE `ln` = '%s';",name));
	if queryErr == nil{
		if len(queryRes) > 0{
			//数据中存在相同账号，注册失败
			resStruct.Code = resFaild;
			resStruct.Msg = "账号已存在，注册失败";
			ctx.WriteString(tools.StructToJSONStr(resStruct))
			ctx.WriteString(tools.StructToJSONStr(resStruct));
			return;
		}
	}else{
		ctx.WriteString(tools.StructToJSONStr(resStruct));
		return;
	}

	gmluuid,uuidErr:= uuid.NewV4();//生成唯一ID
	if uuidErr != nil{
		fmt.Println("uuid生成失败");
		//插入失败
		ctx.WriteString(tools.StructToJSONStr(resStruct))
		return;
	}
 	//向bussinessUsers插入数据
 	res,err := ls.SqlPro.Exec(fmt.Sprintf("insert `BusinessUsers`(`bid_str`,`ln`,`pwd`) values('%s','%s','%s');",gmluuid,name,pwd));
	if err == nil{
		if line,_ := res.RowsAffected() ;line > 0{
			//向BusinessInfo插入数据
			res2,err2 := ls.SqlPro.Exec(fmt.Sprintf("insert `BusinessInfo`(`bid_str`,`des`,`bname`) values('%s','%s','%s');",gmluuid,"",""));
			if err2 == nil{
				if line2,_:=res2.RowsAffected();line2 >0{
					resStruct.Code = resOK;
					resStruct.Msg = "注册成功"
					//添加成功
					ctx.WriteString(tools.StructToJSONStr(resStruct))
				}else{
					//插入失败
					ctx.WriteString(tools.StructToJSONStr(resStruct))
				}
			}else{
				//插入失败
				ctx.WriteString(tools.StructToJSONStr(resStruct))
			}
			
		}else{
			//插入失败
			ctx.WriteString(tools.StructToJSONStr(resStruct))
		}
	}else{
		//插入失败
		ctx.WriteString(tools.StructToJSONStr(resStruct))
	}
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

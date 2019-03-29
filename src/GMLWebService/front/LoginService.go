package front;

import(
	"../../models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"fmt"
	"../../tools"
	"github.com/satori/go.uuid"
	"github.com/kataras/iris/core/router"
	"strconv"
)

var(
	resFaild = "0";
	resOK = "1";
	InvalideLoginName = "2";//账号不存在
	PWDFaild = "3";//密码错误
	UnSinged = "4";//登录时效已过期或未登录
	InvalideMethod = "1001";//无效的method请求
)

/**
前端 登录相关服务
*/
type LoginService struct{
	SqlPro models.SQLInterface;
	App *iris.Application;
	SignGroup *router.Party;
	Sm *sessions.Sessions;
}

/**
开始监听前端服务接口的调用
*/
func (ls *LoginService)Start(){
	ls.App.Get("/front/signByAccount",ls.sign);
	ls.App.Get("/front/signOut",ls.signOut);
	ls.App.Get("/front/registerAccount",ls.registerAccount);
	ls.App.Get("/front/findpwd",ls.findpwd);
	(*ls.SignGroup).Get("/changepwd",ls.changepwd);
	(*ls.SignGroup).Get("/getVerificationCode",ls.getVerificationCode);
	
}

/**
检查是否登录成功
*/
func (ls *LoginService)MW_CheckSinged(ctx iris.Context){
	sess := ls.Sm.Start(ctx);
	if signed,err:= sess.GetBoolean("isSinged");err == nil && signed==true{
		ctx.Next();//执行下一个处理环节
	}else{
		resStruct := models.CurrentResponse{Code:UnSinged,Msg:"登录已过期，请重新登录"};
		ctx.WriteString(tools.StructToJSONStr(resStruct));
	}
}

/**
通过账号密码登录
*/
func (ls *LoginService)sign(ctx iris.Context){
	//取参数
	sc := tools.Pack(ctx);
	loginName := sc.GetStr("ln");
	loginPWD := sc.GetStr("pwd");
	res,err := ls.SqlPro.Query(fmt.Sprintf("SELECT * FROM `BusinessUsers` WHERE `ln` = '%s';",loginName))
	//PWDFaild = "3";//密码错误
	resStruct := models.CurrentResponse{Code:InvalideLoginName};
	if err == nil{
		if len(res) > 0{
			if pwd,gcontains:=res[0]["pwd"] ;gcontains==true && string(pwd) == loginPWD{
				//存session到服务器，存cookie到用户本地
				sess := ls.Sm.Start(ctx);//获取session，如果没有则自动创建，并且自动将sessionID存入cookie
				//将用户信息存入session
				sess.Set("ln",loginName);
				sess.Set("pwd",loginPWD);
				uid,_:= strconv.ParseUint(string(res[0]["bid"]),0,64);
				sess.Set("uid",uid);
				sess.Set("isSinged",true);
				ls.Sm.ShiftExpiration(ctx);//更新session的失效时间
				//返回成功登录结果
				resLoginStruct := models.LoginStruct{};
				resLoginStruct.Code = resOK;
				resLoginStruct.BidStr = string(res[0]["bid_str"]);
				resLoginStruct.Msg = "登录成功";
				ctx.WriteString(tools.StructToJSONStr(resLoginStruct));
			}else{
				resStruct.Code = PWDFaild;
				resStruct.Msg = "密码错误";
				ctx.WriteString(tools.StructToJSONStr(resStruct));
			}
		}else{
			resStruct.Msg = "账号不存在";
			ctx.WriteString(tools.StructToJSONStr(resStruct));
		}
	}else{
		resStruct.Msg = "账号不存在";
		ctx.WriteString(tools.StructToJSONStr(resStruct));
	}
}

/**
登出
*/
func (ls *LoginService)signOut(ctx iris.Context){
	sess := ls.Sm.Start(ctx);
	sess.Destroy();
	ctx.Write([]byte("登出成功"));
}

/**
注册账号
*/
func (ls *LoginService)registerAccount(ctx iris.Context){
	//name := ctx.URLParam("ln");//gv(registerForm["ln"]);
	sc := tools.Pack(ctx);
	name := sc.GetStr("ln");
	pwd := sc.GetStr("pwd");
	resStruct := models.CurrentResponse{Code:resFaild,Msg:"注册失败"};
	if len(name) < 5{
		resStruct.Msg = "账号长度不能小于5";
		ctx.WriteString(tools.StructToJSONStr(resStruct));
		return;
	}else if len(pwd) < 6{
		resStruct.Msg = "密码长度不能小于6";
		ctx.WriteString(tools.StructToJSONStr(resStruct));
		return;
	}
 	//检查同名用户是否存在
 	queryRes,queryErr := ls.SqlPro.Query(fmt.Sprintf("SELECT `ln` FROM `BusinessUsers` WHERE `ln` = '%s';",name));
	if queryErr == nil{
		if len(queryRes) > 0{
			//数据中存在相同账号，注册失败
			resStruct.Code = resFaild;
			resStruct.Msg = "账号已存在，注册失败";
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
func (ls *LoginService)changepwd(ctx iris.Context){
	sc := tools.Pack(ctx);
	oldp := sc.GetStr("oldp");
	newp := sc.GetStr("newp");
	bid_str := sc.GetStr("bidstr");
	res,err := ls.SqlPro.Exec(fmt.Sprintf("update `BusinessUsers` set `pwd`='%s' where `bid_str`='%s' and `pwd`='%s';",newp,bid_str,oldp));
	resObj := models.CurrentResponse{Code:resFaild,Msg:"修改密码失败"};
	if err == nil{
		if line,_ :=res.RowsAffected();line > 0{
			resObj.Code = resOK;
			resObj.Msg = "修改密码成功";
			//更新session
			sess := ls.Sm.Start(ctx);
			sess.Set("pwd",newp);//更新session中存储的密码
		}
	}
	ctx.WriteString(tools.StructToJSONStr(resObj));
}

/**
找回密码
*/
func (ls *LoginService)findpwd(ctx iris.Context){

}

/**
获取验证码
*/
func (ls *LoginService)getVerificationCode(ctx iris.Context){

}

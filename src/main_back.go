package main
import(
	"fmt"
	"net/http"
	"time"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	"database/sql"
	"github.com/satori/go.uuid"
	"./models"
	"encoding/json"
	"bytes"
	"encoding/binary"
	"strconv"
)
//定义response的错误码
var(
	resFaild = "0";
	resOK = "1";
	InvalideLoginName = "2";//账号不存在
	PWDFaild = "3";//密码错误
	InvalideMethod = "1001";//无效的method请求
)
//定义请求处理函数的集合
var(
	postExec = map[string]func(http.ResponseWriter, *http.Request){};
	getExec = map[string]func(http.ResponseWriter, *http.Request){};
	putExec = map[string]func(http.ResponseWriter, *http.Request){};
) 
var (
	sqlEngine *xorm.Engine
	sqlErr error
	sqlHeartOffset  = time.Second * 30;
)


func main(){
	fmt.Printf("欢迎使用");
	startSQL();
	//添加接口处理函数
	postExec["/getCompensationInfo"] = getCompensationInfo;
	postExec["/businessLogin"] = businessLogin;
	postExec["/businessRegister"] = businessRegister;
	postExec["/businessFindPassword"] = businessFindPassword;
	postExec["/businessLogout"] = businessLogout;
	postExec["/getBusinessInfo"] = getBusinessInfo;
	postExec["/updateBusinessInfo"] = updateBusinessInfo;
	postExec["/businessChangePassword"] = businessChangePassword;
	postExec["/createProject"] = createProject;
	postExec["/updateProject"] = updateProject;
	postExec["/deleteProject"] = deleteProject;
	postExec["/getProjectList"] = getProjectList;
	postExec["/createLesson"] = createLesson;
	postExec["/updateLesson"] = updateLesson;
	postExec["/deleteLesson"] = deleteLesson;
	postExec["/getLessonListByBidstr"] = getLessonListByBidstr;
	postExec["/getLessonInfoByCid"] = getLessonInfoByCid;
	postExec["/bookLesson"] = bookLesson;
	postExec["/getOwnerUIDByCid"] = getOwnerUIDByCid;
	postExec["/getUserInfoByUID"] = getUserInfoByUID;
	//初始化http接口监听
	server := &http.Server{
		Addr: "0.0.0.0:39855",
		Handler: HttpHandlerImpl{},
		ReadTimeout:    30 * time.Second/*设置请求超时为30秒*/,
		WriteTimeout:   30 * time.Second/*设置请求超时为30秒*/,
		MaxHeaderBytes: 50 << 20/*限制包大小50mb*/};
	server.ListenAndServe();//Https监听请使用server.ListenAndServeTLS("./cert/1540854920368.pem", "./cert/1540854920368.key")；//这两个文件在main的统计目录的cert文件夹中
	fmt.Printf("代码千万别添加到我之后，因为上面那孙子开启了一个消息循环，代码无法执行到这里");
}

/*
启动sql框架
*/
func startSQL(){
	//user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
	//sqlEngine,sqlErr = xorm.NewEngine("mysql","gmlmaster:123456@tcp(39.106.135.11:32306)/GMLPlanDB?charset=utf8");
	sqlEngine,sqlErr = xorm.NewEngine("mysql","test:123456@tcp(172.16.70.11:3306)/51talk_svc?charset=utf8");
	sqlEngine.Logger().SetLevel(core.LOG_WARNING);//控制台打印sql日志
	sqlEngine.SetMaxIdleConns(50);//设置连接池的空闲数大小
	sqlEngine.SetMaxOpenConns(50);//设置最大打开连接数
	if sqlErr == nil{
		// sqlinfo,err := sqlEngine.DBMetas();
		// if err != nil{
		// 	fmt.Printf("错误信息%v",err);
		// }
		// fmt.Printf("表数据%v",sqlinfo);
		fmt.Println("\n数据库连接成功");
		//维持sql长连接
		go sqlHeart();
	}else{
		fmt.Printf("\n数据库连接错误:%v\n",sqlErr);
	}
}

/*
维持sql长连接
*/
func sqlHeart(){
	for{
		sqlEngine.Ping();
		time.Sleep(sqlHeartOffset);
	}
}

/**
执行sql语句
*/
func SQLExec(str string) (res sql.Result, err error){
	res,err = sqlEngine.Exec(str);
	if err != nil{
		fmt.Println("sql语句执行错误:",str,"错误原因:",err);
	}
	return res,err;
}

func SQLQuery(str string)(res []map[string][]byte, err error){
	res,err = sqlEngine.Query(str);
	if err != nil{
		fmt.Println("sql查询语句执行失败:",str,"错误原因:",err);
	}
	return res,err;
}

func int32ToByte(arg int32)[]byte{
	tempBytes := make([]byte,0);
	buf := bytes.NewBuffer(tempBytes);
	err := binary.Write(buf,binary.BigEndian,arg);
	if err == nil{
		return buf.Bytes();
	}
	return make([]byte,0);
}

func byteToInt32(arg []byte)int32{
	buf := bytes.NewBuffer(arg)
    var res int32;
	err := binary.Read(buf, binary.BigEndian, &res)//以网络字节序（打字节序）的方式从流中读取一个int32的整数
	if err == nil{
		return res;
	}
	return 0;
}

//声明用于实现Handler接口的结构
type HttpHandlerImpl struct{

}
//实现Handler接口的 ServeHTTP函数
func (hand HttpHandlerImpl) ServeHTTP(response http.ResponseWriter,request *http.Request){
	response.Header().Set("Access-Control-Allow-Origin", "*")//允许访问所有域
	response.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	response.Header().Set("content-type", "application/json") //返回数据为json
	var hostPort = request.RemoteAddr;//请求端的ip和端口
	var appName = searchAppNameByURI(request.RequestURI);//要请求的接口名
	fmt.Printf("%s请求了接口:%s\n",hostPort,appName);
	request.ParseForm();//如果不调用这段代码则PostForm 和PostForm中，无数据。
	var execMap map[string]func(http.ResponseWriter,*http.Request);
	var metho = request.Method;
	var resObj = models.CurrentResponse{};
	switch metho {
	case "POST":execMap = postExec;
	case "GET":execMap = getExec;
	default:
		resObj.Code = InvalideMethod;
		resObj.Msg = "不支持的Method";
		fmt.Fprintln(response,structToJSONStr(resObj))
	}

	if execF,isContains := execMap[appName];isContains{
		execF(response,request);//执行指定的处理函数
	}else{
		resObj.Code = InvalideMethod;
		resObj.Msg = fmt.Sprintf("%s不支持%s请求",appName,metho);
		fmt.Fprintln(response,structToJSONStr(resObj))
	}
	
}

/**
struct转字符串
*/
func structToJSONStr(argStruct interface{})string{
	bytes,jsonerr := json.Marshal(argStruct);
	if jsonerr == nil{
		return string(bytes);
	}else{
		return "{}";
	}
}

/**
从URI中寻找appName（处理函数名）
*/
func searchAppNameByURI(uri string) string{
	idx := strings.Index(uri,"?");
	if idx < 0{
		return uri;
	}else{
		return uri[0:idx];
	}
}
/**
获取请求参数中的值
*/
func gv(arg []string) string{
	if len(arg) > 0{
		return arg[0];
	}else{
		return "";
	}
}

/**
获取配课信息
*/
func getCompensationInfo(response http.ResponseWriter, request *http.Request){
	//取参数
	reqInfo := request.Form;
	appoint_id := gv(reqInfo["appoint_id"]);
	appoint_type := gv(reqInfo["appoint_type"]);
	stu_id := gv(reqInfo["stu_id"]);
	teacher_id := gv(reqInfo["teacher_id"]);
	start_time := gv(reqInfo["start_time"]);
	end_time := gv(reqInfo["end_time"]);
	appoint_status := gv(reqInfo["appoint_status"]);
	resStruct := models.CompensationInfoResponse{Code:"0",Message:"",Timestamp:time.Now().Unix(),Res:0,Success:false}
	if "" == appoint_id || "" == appoint_type || "" == stu_id || "" == teacher_id || "" == start_time || "" == end_time || "" == appoint_status{
		resStruct.Message = fmt.Sprintf("请求参数无效，请求内容为：%s",structToJSONStr(reqInfo));
		fmt.Fprintln(response,structToJSONStr(resStruct));
	}else{
		//开始数据库查询
		//res,err := SQLQuery(fmt.Sprintf("SELECT * FROM `BusinessUsers` WHERE `ln` = '%s';","abc"))
	}
}

/**
企业账号登录
*/
func businessLogin(response http.ResponseWriter, request *http.Request){
	//取参数
	loginInfo := request.Form;
	loginName := gv(loginInfo["ln"]);
	loginPWD := gv(loginInfo["pwd"]);
	res,err := SQLQuery(fmt.Sprintf("SELECT * FROM `BusinessUsers` WHERE `ln` = '%s';",loginName))
	//PWDFaild = "3";//密码错误
	resStruct := models.CurrentResponse{Code:InvalideLoginName};
	if err == nil{
		if len(res) > 0{
			if pwd,gcontains:=res[0]["pwd"] ;gcontains==true && string(pwd) == loginPWD{
				resLoginStruct := models.LoginStruct{};
				resLoginStruct.Code = resOK;
				resLoginStruct.BidStr = string(res[0]["bid_str"]);
				resLoginStruct.Msg = "登录成功";
				fmt.Fprintln(response,structToJSONStr(resLoginStruct));
			}else{
				resStruct.Code = PWDFaild;
				resStruct.Msg = "密码错误";
				fmt.Fprintln(response,structToJSONStr(resStruct));
			}
		}else{
			resStruct.Msg = "账号不存在";
			fmt.Fprintln(response,structToJSONStr(resStruct));
		}
	}else{
		resStruct.Msg = "账号不存在";
		fmt.Fprintln(response,structToJSONStr(resStruct));
	}
}

/*
企业账号注册
*/
func businessRegister(response http.ResponseWriter,request *http.Request){
	registerForm := request.Form;
	name := gv(registerForm["ln"]);
	pwd := gv(registerForm["pwd"]);
	resStruct := models.CurrentResponse{Code:resFaild,Msg:"注册失败"};
	//检查同名用户是否存在
	queryRes,queryErr := SQLQuery(fmt.Sprintf("SELECT `ln` FROM `BusinessUsers` WHERE `ln` = '%s';",name));
	if queryErr == nil{
		if len(queryRes) > 0{
			//数据中存在相同账号，注册失败
			resStruct.Code = resFaild;
			resStruct.Msg = "账号已存在，注册失败";
			fmt.Fprintln(response,structToJSONStr(resStruct));
			return;
		}
	}else{
		fmt.Fprintln(response,structToJSONStr(resStruct));
		return;
	}

	gmluuid,uuidErr:= uuid.NewV4();//生成唯一ID
	if uuidErr != nil{
		fmt.Println("uuid生成失败");
		//插入失败
		fmt.Fprintln(response,structToJSONStr(resStruct))
		return;
	}
	//向bussinessUsers插入数据
	res,err := SQLExec(fmt.Sprintf("insert `BusinessUsers`(`bid_str`,`ln`,`pwd`) values('%s','%s','%s');",gmluuid,name,pwd));
	if err == nil{
		if line,_ := res.RowsAffected() ;line > 0{
			//向BusinessInfo插入数据
			res2,err2 := SQLExec(fmt.Sprintf("insert `BusinessInfo`(`bid_str`,`des`,`bname`) values('%s','%s','%s');",gmluuid,"",""));
			if err2 == nil{
				if line2,_:=res2.RowsAffected();line2 >0{
					resStruct.Code = resOK;
					resStruct.Msg = "注册成功"
					//添加成功
					fmt.Fprintln(response,structToJSONStr(resStruct))
				}else{
					//插入失败
					fmt.Fprintln(response,structToJSONStr(resStruct))
				}
			}else{
				//插入失败
				fmt.Fprintln(response,structToJSONStr(resStruct))
			}
			
		}else{
			//插入失败
			fmt.Fprintln(response,structToJSONStr(resStruct))
		}
	}else{
		//插入失败
		fmt.Fprintln(response,structToJSONStr(resStruct))
	}
	
}

/**
企业账号找回密码
*/
func businessFindPassword(response http.ResponseWriter,request *http.Request){
	fmt.Fprintln(response,fmt.Sprintf("{\"code\":\"%s\",\"msg\":\"功能待开发\"}",resOK))
}

/*
企业账户登出
*/
func businessLogout(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	reqLoginName := gv(reqForm["ln"]);//企业账号
	reqBusinessCode := gv(reqForm["bc"]);//企业代码
	res,err:=SQLQuery(fmt.Sprintf("Select `ln` from `BusinessUsers` where `ln` = '%s' and `bid_str` = '%s'",reqLoginName,reqBusinessCode));
	resStruct := models.CurrentResponse{};
	if err == nil && len(res) > 0{
		//登出成功
		resStruct.Code = resOK;
		resStruct.Msg = "登出成功";
		fmt.Fprintln(response,structToJSONStr(resStruct));
	}else{
		//登出失败
		resStruct.Code = resFaild;
		resStruct.Msg = "登出失败，账号无效";
		fmt.Fprintln(response,structToJSONStr(resStruct));
	}
}

/**
获取企业信息
*/
func getBusinessInfo(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	bc := gv(reqForm["bc"]);
	res,err := SQLQuery(fmt.Sprintf("Select * from `BusinessInfo` where `bid_str` = '%s';",bc));
	if err == nil && len(res) > 0{
		resSQL := res[0];
		result := models.BusinessStruct{};
		result.Code = resOK;
		result.BidStr = string(resSQL["bid_str"]);
		result.BusinessName = string(resSQL["bname"]);
		result.BusinessDes = string(resSQL["des"]);
		fmt.Fprintln(response,structToJSONStr(result));
	}else{
		//获取信息失败
		resStruct := models.CurrentResponse{Code:resFaild,Msg:"获取信息失败"}
		fmt.Fprintln(response,structToJSONStr(resStruct));
	}
}

/**
更新企业信息
*/
func updateBusinessInfo(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	bc := gv(reqForm["bc"]);
	businessName := gv(reqForm["bmn"]);
	businessDes := gv(reqForm["bmdes"]);
	res,err := SQLExec(fmt.Sprintf("update `BusinessInfo` set `des`='%s',`bname`='%s' where `bid_str`='%s';",businessDes,businessName,bc))
	if err == nil{
		if line,_:=res.RowsAffected();line > 0{
			//更新成功
			result := models.BusinessStruct{};
			result.Code = resOK;
			result.BidStr = bc;
			result.BusinessName = businessName;
			result.BusinessDes = businessDes;
			fmt.Fprintln(response,structToJSONStr(result));
		}else{
			//更新失败
			resStruct := models.CurrentResponse{Code:resFaild,Msg:"更新企业信息失败"}
			fmt.Fprintln(response,structToJSONStr(resStruct));
		}
	}else{
		//更新失败
		resStruct := models.CurrentResponse{Code:resFaild,Msg:"更新企业信息失败"}
		fmt.Fprintln(response,structToJSONStr(resStruct));
	}
}

/**
创建项目
*/
func createProject(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	pName := gv(reqForm["pName"]);
	pDes := gv(reqForm["pDes"]);
	bid_str := gv(reqForm["bidstr"]);
	res,err := SQLExec(fmt.Sprintf("insert `projectInfo`(`pname`,`bid_str`,`pdes`) values('%s','%s','%s');",pName,bid_str,pDes));
	resObj := models.CurrentResponse{Code:resFaild,Msg:"新增项目失败"};
	if err == nil{
		if line,_ :=res.RowsAffected();line > 0{
			if pid,errtemp := res.LastInsertId();errtemp == nil{
				resultObj := models.ProjectStruct{Code:resOK};
				resultObj.Pname = pName;
				resultObj.BidStr = bid_str;
				resultObj.Pid = fmt.Sprintf("%v",pid);
				resultObj.Pdes = pDes;
				fmt.Fprintln(response,structToJSONStr(resultObj));
			}else{
				fmt.Fprintln(response,structToJSONStr(resObj));
			}
		}else{
			fmt.Fprintln(response,structToJSONStr(resObj));
		}
	}else{
		fmt.Fprintln(response,structToJSONStr(resObj));
	}
	
}

/**
更新项目
*/
func updateProject(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	pid := gv(reqForm["pid"]);
	bid_str := gv(reqForm["bidstr"]);
	pname := gv(reqForm["pName"]);
	pdes := gv(reqForm["pDes"]);
	res,err := SQLExec(fmt.Sprintf("update `projectInfo` set `pname`='%s',`pdes`='%s' where `pid`=%s and `bid_str` = '%s';",pname,pdes,pid,bid_str));
	resObj := models.CurrentResponse{Code:resFaild,Msg:"修改项目信息失败"};
	if err == nil{
		if line,_ :=res.RowsAffected();line > 0{
			resObj.Code = resOK;
			resObj.Msg = "更新项目信息成功";
		}
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}

/**
删除项目
*/
func deleteProject(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	pid := gv(reqForm["pid"]);
	bid_str := gv(reqForm["bidstr"]);
	res,err := SQLExec(fmt.Sprintf("delete from `projectInfo` where `pid`=%s and `bid_str` = '%s';",pid,bid_str));
	resObj := models.CurrentResponse{Code:resFaild,Msg:"删除项目信息失败"};
	if err == nil{
		if line,_ :=res.RowsAffected();line > 0{
			resObj.Code = resOK;
			resObj.Msg = "删除项目信息成功";
		}
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}

/*
获取项目列表
*/
func getProjectList(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	bid_str := gv(reqForm["bidstr"]);
	res,err := SQLQuery(fmt.Sprintf("select * from `projectInfo` where `bid_str` = '%s';",bid_str));
	resObj := models.ProjectListStruct{Code:resFaild,Msg:[]models.ProjectListItem{}};
	if err == nil{
		j := len(res);
		if j > 0{
			resObj.Code = resOK;
			for i := 0;i<j;i++{
				obj := res[i];
				pid, converErr := strconv.ParseInt(string(obj["pid"]), 10, 32)
				if converErr == nil{
					item := models.ProjectListItem{
						BidStr:string(obj["bid_str"]),
						Pid:fmt.Sprintf("%v",pid),
						Pname:string(obj["pname"]),
						Pdes:string(obj["pdes"])};
					resObj.Msg = append(resObj.Msg,item);
				}
			}
		}
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}

/*
企业账户更改密码
*/
func businessChangePassword(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	oldp := gv(reqForm["oldp"]);
	newp := gv(reqForm["newp"]);
	bid_str := gv(reqForm["bidstr"]);
	res,err := SQLExec(fmt.Sprintf("update `BusinessUsers` set `pwd`='%s' where `bid_str`='%s' and `pwd`='%s';",newp,bid_str,oldp));
	resObj := models.CurrentResponse{Code:resFaild,Msg:"修改密码失败"};
	if err == nil{
		if line,_ :=res.RowsAffected();line > 0{
			resObj.Code = resOK;
			resObj.Msg = "修改密码成功";
		}
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}

/**
创建课程信息
*/
func createLesson(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	b_cid := gv(reqForm["b_cid"]);
	b_key := gv(reqForm["bidstr"]);
	b_lessonInfo := gv(reqForm["b_lessonInfo"]);
	startTime := gv(reqForm["startTime"]);
	lessonTimeLengh := gv(reqForm["lessonTimeLengh"]);
	maxCap := gv(reqForm["maxCap"]);
	maxLine := gv(reqForm["maxLine"]);
	pid := gv(reqForm["pid"]);
	res,err := SQLExec(fmt.Sprintf("insert `lessones`(`b_cid`,`b_key`,`b_lessonInfo`,`startTime`,`lessonTimeLengh`,`maxCap`,`maxLine`,`pid`) values('%s','%s','%s',%s,%s,%s,%s,%s);",b_cid,b_key,b_lessonInfo,startTime,lessonTimeLengh,maxCap,maxLine,pid));
	resObj := models.CreateLessonCallBack{Code:resFaild,Msg:map[string]string{}};
	if err == nil{
		if line,_ :=res.RowsAffected();line > 0{
			if lastIdx,lastIdxerr :=res.LastInsertId();lastIdxerr == nil{
				resObj.Code = resOK;
				resObj.Msg["Cid"] = fmt.Sprintf("%v",lastIdx);
				resObj.Msg["Bcid"] = b_cid;
			}
		}
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}

/**
更新课程信息
*/
func updateLesson(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	cid := gv(reqForm["cid"]);
	b_cid := gv(reqForm["b_cid"]);
	b_key := gv(reqForm["bidstr"]);
	b_lessonInfo := gv(reqForm["b_lessonInfo"]);
	startTime := gv(reqForm["startTime"]);
	lessonTimeLengh := gv(reqForm["lessonTimeLengh"]);
	maxCap := gv(reqForm["maxCap"]);
	maxLine := gv(reqForm["maxLine"]);
	pid := gv(reqForm["pid"]);

	res,err := SQLExec(fmt.Sprintf("update `lessones` set `b_cid`='%s',`b_lessonInfo`='%s',`startTime`=%s,`lessonTimeLengh`=%s,`maxCap`=%s,`maxLine`=%s,`pid`=%s where `cid`=%s and `b_key`='%s';",b_cid,b_lessonInfo,startTime,lessonTimeLengh,maxCap,maxLine,pid,cid,b_key));
	resObj := models.CurrentResponse{Code:resFaild,Msg:"更新课程信息失败"};
	if err == nil{
		if line,_ :=res.RowsAffected();line > 0{
			resObj.Code = resOK;
			resObj.Msg = "更新课程信息成功";
		}
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}

/**
删除课程信息
*/
func deleteLesson(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	cids := gv(reqForm["delcids"]);
	arr :=strings.Split(cids,",");
	j := len(arr);
	resultCids := "";
	for i:=0;i<j;i++{
		if n,convererr := strconv.ParseInt(arr[i],10,32);convererr == nil{
			resultCids += fmt.Sprintf(",%v",n);
		}
	}
	if len(resultCids) > 0{
		resultCids = resultCids[1:];
		res,err := SQLExec(fmt.Sprintf("delete from `lessones` where `cid` in (%s);",resultCids));
		resObj := models.CurrentResponse{Code:resFaild,Msg:"删除课程信息失败"};
		if err == nil{
			if line,_ :=res.RowsAffected();line > 0{
				resObj.Code = resOK;
				resObj.Msg = "删除课程信息成功";
			}
		}
		fmt.Fprintln(response,structToJSONStr(resObj));
	}
}

/**
获取课程信息列表
*/
func getLessonListByBidstr(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	bidstr := gv(reqForm["bid_str"]);
	res,err := SQLQuery(fmt.Sprintf("select * from `lessones` where `b_key` = '%s'",bidstr));
	resobj := models.GetLessonsInfo{Code:resFaild,Msg:[]models.LessonInfoItem{}};
	if err == nil{
		if j:= len(res);j>0{
			resobj.Code = resOK;
			for i:= 0;i<j;i++{
				obj := res[i];
				item := models.LessonInfoItem{};
				item.Cid = tempFunc1(string(obj["cid"]));
				item.Bcid = string(obj["b_cid"]);
				item.BlescustomInfo = string(obj["b_lessonInfo"]);
				item.StartTimeInterval = tempFunc1(string(obj["startTime"]));
				item.LessonTimeLength = tempFunc1(string(obj["lessonTimeLengh"]));
				item.MaxCap = tempFunc1(string(obj["maxCap"]));
				item.MaxLine = tempFunc1(string(obj["maxLine"]));
				item.Pid = tempFunc1(string(obj["pid"]));
				resobj.Msg = append(resobj.Msg,item);
			}
		}
	}
	fmt.Fprintln(response,structToJSONStr(resobj));
}

/**
获取课程信息根据CID
*/
func getLessonInfoByCid(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	cid := gv(reqForm["cid"]);
	res,err := SQLQuery(fmt.Sprintf("select * from `lessones` where `cid` = %s;",cid));
	resobj := models.GetLessonsInfo{Code:resFaild,Msg:[]models.LessonInfoItem{}};
	if err == nil{
		if j:= len(res);j>0{
			resobj.Code = resOK;
			obj := res[0];
				item := models.LessonInfoItem{};
				item.Cid = tempFunc1(string(obj["cid"]));
				item.Bcid = string(obj["b_cid"]);
				item.BlescustomInfo = string(obj["b_lessonInfo"]);
				item.StartTimeInterval = tempFunc1(string(obj["startTime"]));
				item.LessonTimeLength = tempFunc1(string(obj["lessonTimeLengh"]));
				item.MaxCap = tempFunc1(string(obj["maxCap"]));
				item.MaxLine = tempFunc1(string(obj["maxLine"]));
				item.Pid = tempFunc1(string(obj["pid"]));
				resobj.Msg = append(resobj.Msg,item);
		}
	}
	fmt.Fprintln(response,structToJSONStr(resobj));
}

func tempFunc1(str string)string{
	res,err := strconv.ParseInt(str, 10, 32);
	if err == nil{
		return fmt.Sprintf("%v",res);
	}else{
		return "0";
	}
}

/**
用户约课
*/
func bookLesson(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	bid_key := gv(reqForm["bid_key"]);
	cid := gv(reqForm["cid"]);
	b_uid := gv(reqForm["b_uid"]);
	b_user_nickName := gv(reqForm["b_user_nickName"]);
	b_user_headerImage := gv(reqForm["b_user_headerImage"]);
	b_user_sex := gv(reqForm["b_user_sex"]);
	isOwnner := gv(reqForm["isOwnner"]);//是否是主讲
	//获取userID
	uid := getUidBybidkeyAndBuid(bid_key,b_uid);
	resObj := models.BookLessonRes{Code:resFaild,Msg:"约课失败"};
	if uid == "0"{
		//无用户信息，则先添加用户信息到用户表
		uid = addUserInfo(bid_key,b_uid,b_user_nickName,b_user_headerImage,b_user_sex);
	}
	if uid != "0"{
		//添加预约信息
		res,err := SQLExec(fmt.Sprintf("insert `BookLesson`(`cid`,`uid`,`isOwnner`) values(%s,%s,%s);",cid,uid,isOwnner));
		if err == nil{
			if line,_ :=res.RowsAffected();line > 0{
				resObj.Code = resOK;
				resObj.Msg = "约课成功";
				resObj.Uid = uid;
				resObj.Cid = cid;
			}
		}
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}

/**
获取课程的拥有者UID
*/
func getOwnerUIDByCid(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	cid := gv(reqForm["cid"]);
	resObj := models.CurrentResponse{Code:resFaild,Msg:"获取课程的拥有者失败"};
	res,err := SQLQuery(fmt.Sprintf("select `uid` from `BookLesson` where `cid` = %s and `isOwnner` = 1",cid))
	if err == nil && len(res) > 0{
		obj := res[0];
		resObj.Code = resOK;
		resObj.Msg = tempFunc1(string(obj["uid"]));
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}

/**
根据企业ID和企业内部用户ID，获取平台userID
*/
func getUidBybidkeyAndBuid(bid_key,b_uid string)string{
	resUid := "0";
	res,err := SQLQuery(fmt.Sprintf("select `uid` from `users` where `b_Key` = '%s' and `b_uid` = '%s';",bid_key,b_uid));
	if err == nil && len(res) > 0{
		resUid = tempFunc1(string(res[0]["uid"]));
	}
	return resUid;
}

/**
向平台用户表中添加用户信息，返回平台userID
*/
func addUserInfo(bid_key,b_uid,b_user_nickName,b_user_headerImage,b_user_sex string)string{
	resUid := "0";
	createTime := time.Now().Unix();
	res,err := SQLExec(fmt.Sprintf("insert `users`(`nickname`,`b_Key`,`b_uid`,`createTime`,`headerImage`,`sex`) values('%s','%s','%s',%v,'%s',%s);",b_user_nickName,bid_key,b_uid,createTime,b_user_headerImage,b_user_sex));
	if err == nil{
		if line,_:=res.RowsAffected() ;line > 0{
			if idx,idxerr := res.LastInsertId();idxerr == nil{
				resUid = fmt.Sprintf("%v",idx);
			}
		}
	}
	return resUid;
}

/**
根据用户ID获取用户信息
*/
func getUserInfoByUID(response http.ResponseWriter,request *http.Request){
	reqForm := request.Form;
	uid := gv(reqForm["uid"])
	res,err := SQLQuery(fmt.Sprintf("select * from `users` where `uid` = %s;",uid));
	resObj := models.UserInfo{Code:resFaild};
	if err == nil && len(res) > 0{
		obj := res[0];
		resObj.Code = resOK;
		resObj.Uid = uid;
		resObj.NickName = string(obj["nickname"]);
		resObj.BUID = string(obj["b_uid"]);
		resObj.CreateTime = tempFunc1(string(obj["createTime"]));
		resObj.HeaderImg = string(obj["headerImage"]);
		resObj.Sex = fmt.Sprintf("%v",obj["sex"][0]);
	}
	fmt.Fprintln(response,structToJSONStr(resObj));
}
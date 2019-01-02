package models

type CurrentResponse struct{
	Code string
	Msg string
}

type LoginStruct struct{
	Code string
	Msg string
	BidStr string/*企业标识码*/
}

type BusinessStruct struct{
	Code string
	BidStr string/*企业标识码*/
	BusinessName string/*企业名称*/
	BusinessDes string/*企业描述*/
}

type ProjectStruct struct{
	Code string
	BidStr string/*企业标识码*/
	Pid string/*企业id*/
	Pname string/*项目名称*/
	Pdes string/*项目描述*/
}

type ProjectListItem struct{
	BidStr string/*企业标识码*/
	Pid string/*企业id*/
	Pname string/*项目名称*/
	Pdes string/*项目描述*/
}

type ProjectListStruct struct{
	Code string
	Msg []ProjectListItem
}

//创建课程回调
type CreateLessonCallBack struct{
	Code string
	Msg map[string]string
}

type LessonInfoItem struct{
	Cid string
	Bcid string
	BlescustomInfo string
	StartTimeInterval string
	LessonTimeLength string
	MaxCap string
	MaxLine string
	Pid string
}


type GetLessonsInfo struct{
	Code string
	Msg []LessonInfoItem
}

type UserInfo struct{
	Code string
	Msg string
	Uid string
	NickName string
	BUID string
	CreateTime string
	HeaderImg string
	Sex string
}

type BookLessonRes struct{
	Code string
	Msg string
	Cid string
	Uid string
}
package tools;

import(
	"github.com/kataras/iris/context"
	cov "strconv"
)

/**
将 iris.Context 封装为SupperContext，便于快捷取参
*/
func Pack(_ctx context.Context)(sc *SupperContext){
	sc = &SupperContext{ctx:_ctx};
	return sc;
}


type SupperContext struct{
	ctx context.Context;
}

/**
获取完整的参数列表，queryParams是get请求的参数列表，formParams是post等等请求的参数列表，二者类型不同，同一时间只有一个不为空
*/
func (sc *SupperContext)GetAllParams()(queryParams map[string]string,formParams map[string][]string){
	method := sc.ctx.Method();
	if method == "GET"{
		return sc.ctx.URLParams(),nil;
	}else if method == "PUT" || method == "POST" || method == "PATCH"{
		return nil,sc.ctx.FormValues();
	}else{
		return nil,nil;
	}
}

/**
获取参数列表中的指定参数的 string value
*/
func (sc *SupperContext)GetStr(key string) string{
	result := "";
	method := sc.ctx.Method();
	if method == "GET"{
		result = sc.ctx.URLParamTrim(key);
	}else if method == "PUT" || method == "POST" || method == "PATCH"{
		result = sc.ctx.PostValueTrim(key);
	}
	return result;
}


/**
获取指定key对应的bool值
*/
func (sc *SupperContext)GetBool(key string)bool{
	method := sc.ctx.Method();
	result := false;
	if method == "GET"{
		result,_= sc.ctx.URLParamBool(key);
	}else if method == "PUT" || method == "POST" || method == "PATCH"{
		result,_= sc.ctx.PostValueBool(key);
	}
	return result;
}

/**
获取指定key对应的int值
*/
func (sc *SupperContext)GetInt(key string)int{
	method := sc.ctx.Method();
	result := -1;
	if method == "GET"{
		result,_= sc.ctx.URLParamInt(key);
	}else if method == "PUT" || method == "POST" || method == "PATCH"{
		result,_= sc.ctx.PostValueInt(key);
	}
	return result;
} 


/**
获取指定key对应的int32值
*/
func (sc *SupperContext)GetInt32(key string)int32{
	strV := sc.GetStr(key);
	var result int32= -1;
	if temp,err := cov.ParseInt(strV,0,32);err == nil{
		result = int32(temp);
	}
	return result;
} 

/**
获取指定key对应的int64值
*/
func (sc *SupperContext)GetInt64(key string)int64{
	strV := sc.GetStr(key);
	var result int64= -1;
	if temp,err := cov.ParseInt(strV,0,64);err == nil{
		result = temp;
	}
	return result;
} 


/**
获取指定key对应的uint32值
*/
func (sc *SupperContext)GetUInt32(key string)uint32{
	strV := sc.GetStr(key);
	var result uint32= 0;
	if temp,err := cov.ParseUint(strV,0,32);err == nil{
		result = uint32(temp);
	}
	return result;
} 

/**
获取指定key对应的uint64值
*/
func (sc *SupperContext)GetUInt64(key string)uint64{
	strV := sc.GetStr(key);
	var result uint64= 0;
	if temp,err := cov.ParseUint(strV,0,64);err == nil{
		result = temp;
	}
	return result;
} 

/**
获取指定key对应的float32值
*/
func (sc *SupperContext)GetFloat32(key string)float32{
	strV := sc.GetStr(key);
	var result float32= 0;
	if temp,err := cov.ParseFloat(strV,32);err == nil{
		result = float32(temp);
	}
	return result;
} 

/**
获取指定key对应的float64值
*/
func (sc *SupperContext)GetFloat64(key string)float64{
	strV := sc.GetStr(key);
	var result float64= 0;
	if temp,err := cov.ParseFloat(strV,64);err == nil{
		result = temp;
	}
	return result;
} 
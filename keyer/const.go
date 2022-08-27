package keyer

//配置中的方法枚举类型
/*
	Method_Value：直接返回一个值，附带值
	Method_Code：响应值，依据其他响应字段获取结果来判断，附带成功码和失败码
	Method_Function：执行函数，附带GetResponseStruct
*/
//type ConfigMethod int
const (
	Method_None     = "none" // ConfigMethod = iota
	Method_Value    = "value"
	Method_Code     = "code"
	Method_Function = "function"
)

//配置中通用的字符常量,$K和SV开头用来区分用户的配置
const (
	SuccssCode    = "success_code"
	FailedCode    = "failed_code"
	PrimaryKey    = "primary"
	PostParamsKey = "post_key"
	LogicConfig   = "logic"
	Method        = "method"
	Value         = "value"
	Default       = "Default"
	None          = "None"
	PostData      = "PostData"
	PostParams    = "PostParams"
	CodeKey       = "code"
)

const (
	DefaultSuccessCode = 0
	DefaultFailedCode  = 1
)

//错误码和信息
var errorMsg = map[int]string{
	100: "there is no post data, router %s",
	101: "%s->%s->%s->%s->%s can not be %s, only can be %s, %s, %s or %s, %s",
	102: "method is not register, config is %s->%s->%s->%s->%s, call %s to register",
	//200之后开头的错误拼接格式："Not Found in generalConfig: %s or %s, check args of http_server.SetGeneralConfig"
	200: "%s: %s or %s, %s",
	//300之后的开头错误拼接格式："Get args in generalConfig error: %s or %s, err is: %s, check args of http_server.SetGeneralConfig"
	300: "%s: %s->%s or %s->%s, err is: %s, %s",
	301: "%s: %s->%s->%s, err is:%s, %s",
	//400之后的开头错误拼接格式："Get args in post data error: url is %s, key is %s, err is %s, check args of post data"
	400: "%s: url is %s, key is %s, err is %s, %s",
}

const (
	getConfigErrorTip   = "Get args in generalConfig error"
	getPostDataErrorTip = "Get args in post data error"
	notFoundConfigTip   = "Not Found in generalConfig"
	checkConfigTip      = "check args of http_server.SetGeneralConfig"
	checkPostDataTip    = "check args of post data"
)

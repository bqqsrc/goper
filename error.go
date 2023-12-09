package goper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bqqsrc/loger"
)

//自定义的通用错误类型
type gpError struct {
	caller    string        //调用的函数
	msgFormat string        //错误的format字符串
	args      []interface{} //错误的字符串
}

func (e *gpError) Error() string {
	fmtStr := fmt.Sprintf("call by: %s,\n error is : %s", e.caller, e.msgFormat)
	return fmt.Sprintf(fmtStr, e.args...)
}

func GoperError(funcName, format string, args ...interface{}) error {
	return &gpError{funcName, format, args}
}

//路由配置检查错误的类型定义
type configError struct {
	listen    int
	host      string
	formatStr string
	args      []interface{}
}

func (e *configError) Error() string {
	fmtStr := fmt.Sprintf("server config error: listen is: %d, host is: %s, err is: %s", e.listen, e.host, e.formatStr)
	return fmt.Sprintf(fmtStr, e.args...)
}

func RouterConfigError(port int, hostName, format string, args ...interface{}) error {
	return &configError{port, hostName, format, args}
}

//路由错误
type routerErr struct {
	port      int
	host      string
	url       string
	formatStr string
	args      []interface{}
}

func (e *routerErr) Error() string {
	fmtStr := fmt.Sprintf("router error, listen is: %d, host is: %s, url is: %s, err is: %s", e.port, e.host, e.url, e.formatStr)
	return fmt.Sprintf(fmtStr, e.args...)
}

func RouterError(listen int, hostName, urlName, format string, args ...interface{}) error {
	return &routerErr{listen, hostName, urlName, format, args}
}

type generalError struct {
	code      int
	caller    string
	formatStr string
	args      []interface{}
}

func (e *generalError) Error() string {
	fmtStr := fmt.Sprintf("call by: %s,\n error is : %s", e.caller, e.formatStr)
	return fmt.Sprintf(fmtStr, e.args...)
}

func GeneralError(code int, caller string, errMsg map[int]string, args ...interface{}) error {
	format, ok := errMsg[code]
	if !ok {
		log.Printf("%d is not found in errMsg", code)
	}
	return &generalError{code, caller, format, args}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Title   string `json:"title"`
	Context string `json:"content"`
}

func Response(e error) []byte {
	loger.Debugf("err is Response1: %s", e)
	loger.Errorf("error: %s", e)
	if genErr, ok := e.(*generalError); ok {
		return genErr.Response()
	} else {
		return []byte(e.Error())
	}
}
func (e *generalError) Response() []byte {
	msgStr := fmt.Sprintf(e.formatStr, e.args...)
	response := ErrorResponse{Code: e.code, Title: e.caller, Context: msgStr}
	loger.Debugf("Response response is %v", response)
	ret, err := json.Marshal(response)
	if err != nil {
		log.Printf("Respose error: json.Marshal(%v) err, err is: %s", response, err)
	}
	loger.Debugf("Response response is %s, err is %s", ret, err)
	return ret
}

type CodeResponse struct {
	Code int `json:"code"`
}

// func ResponseCode(code int) []byte {
// 	c := CodeResponse{code}
// 	ret, err := json.Marshal(c)
// 	if err != nil {
// 		log.Printf("CodeResponse.Respose error: json.Marshal(%v) err, err is: %s", c, err)
// 	}
// 	return ret
// }

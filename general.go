package goper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	GoUrl "net/url"
	"strings"

	"github.com/bqqsrc/loger"
)

type GeneralResponse func(url string, urlParams, bodyData map[string]interface{}, cookies []*http.Cookie) ([]byte, []*http.Cookie, error)

var generalResponses map[string]GeneralResponse

func RegisterGeneralResponse(funcName string, method GeneralResponse) {
	if generalResponses == nil {
		generalResponses = make(map[string]GeneralResponse)
	}
	if _, ok := generalResponses[funcName]; ok {
		loger.Warnf("RegisterGeneralResponse %s twice", funcName)
	}
	generalResponses[funcName] = method
}

func UnRegisterGeneralResponse(funcName string) {
	delete(generalResponses, funcName)
}

func ResetGeneralResponse() {
	generalResponses = nil
}

type CookieChecker func(url string, cookie []*http.Cookie) ([]*http.Cookie, error)

var generalCookieCheckers map[string]CookieChecker

func RegisterCookieChecker(funcName string, checker CookieChecker) {
	if generalCookieCheckers == nil {
		generalCookieCheckers = make(map[string]CookieChecker)
	}
	if _, ok := generalCookieCheckers[funcName]; ok {
		loger.Warnf("RegisterCookieChecker %s twice", funcName)
	}
	generalCookieCheckers[funcName] = checker
}

func UnRegisterCookieChecker(funcName string) {
	delete(generalCookieCheckers, funcName)
}

func ResetCookieChecker() {
	generalCookieCheckers = nil
}

type CookieWriter func(url string, w http.ResponseWriter) http.ResponseWriter

var generalCookieWriters map[string]CookieWriter

func RegisterCookieWriter(funcName string, writer CookieWriter) {
	if generalCookieWriters == nil {
		generalCookieWriters = make(map[string]CookieWriter)
	}
	if _, ok := generalCookieWriters[funcName]; ok {
		loger.Warnf("RegisterCookieWriter %s twice", funcName)
	}
	generalCookieWriters[funcName] = writer
}

func UnRegisterCookieWriter(funcName string) {
	delete(generalCookieWriters, funcName)
}

func ResetCookieWriter() {
	generalCookieWriters = nil
}

func (f *FuncRouter) GeneralServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !CheckRequest(r) {
		return
	}
	funcName := "GeneralServeHTTP"
	r.ParseForm()     // 解析参数，默认是不会解析的
	url := r.URL.Path /// r.URL.Path //路由路径
	// cookieMethod := f.Cookie
	// loger.Debugf("%s, url: %s, method: %s, cookie: %s", funcName, url, f.Method, cookieMethod)
	// cookieCheck, ok := generalCookieCheckers[cookieMethod]
	// if ok {
	// 	if cookies, err := cookieCheck(url, r.Cookies()); err != nil {
	// 		GoResponse(w, Response(err))
	// 		return
	// 	} else if cookies != nil {
	// 		for _, value := range cookies {
	// 			http.SetCookie(w, value)
	// 		}
	// 	}
	// }
	method := f.Method
	reponseFunc, ok := generalResponses[method]
	if !ok {
		loger.Errorf("error: %s, method not foud, url is: %s, method is %s", funcName, url, method)
		http.NotFound(w, r)
		return
	}
	urlParams, err := parseUrlParams(r.URL.RawQuery)
	if err != nil {
		GoResponse(w, Response(err))
		//	w.Write(Response(err))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = GeneralError(1000, funcName, errorMsg, err)
		GoResponse(w, Response(err))
		//	w.Write(Response(err))
		return
	}
	var bodyData map[string]interface{}
	if len(body) > 0 {
		loger.Debugf("body is %v, \n%s", body, string(body))
		if err = json.Unmarshal(body, &bodyData); err != nil {
			err = GeneralError(1001, funcName, errorMsg, err)
			GoResponse(w, Response(err))
			//	w.Write(Response(err))
			return
		}
	}
	//var responseData []byte
	if responseData, cookies, err := reponseFunc(url, urlParams, bodyData, r.Cookies()); err != nil {
		loger.Debugf("err is responseData: %s", err)
		GoResponse(w, Response(err))
		//	w.Write(Response(err))
	} else {
		if cookies != nil && len(cookies) > 0 {
			for _, value := range cookies {
				//loger.Debugf("value is %v", value)
				http.SetCookie(w, value)
			}
		}
		GoResponse(w, responseData)
		//	w.Write(responseData)
	}
}

func (f *FuncRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.GeneralServeHTTP(w, r)
}

func parseUrlParams(rawQuery string) (map[string]interface{}, error) {
	funcName := "parseUrlParams"
	result := make(map[string]interface{})
	if len(rawQuery) > 0 {
		sepStr := "&"
		paramStrArr := strings.Split(rawQuery, sepStr)
		for _, value := range paramStrArr {
			sepStr = "="
			i := strings.Index(value, sepStr)
			k, err := GoUrl.QueryUnescape(value[:i])
			if err != nil {
				return nil, GeneralError(1002, funcName, errorMsg, err)
			}
			v, err := GoUrl.QueryUnescape(value[i+1:])
			if err != nil {
				return nil, GeneralError(1002, funcName, errorMsg, err)
			}
			result[k] = v
		}
	}
	//loger.Println(result)
	return result, nil
}

package keyer

import (
	"net/http"

	"github.com/bqqsrc/loger"
)

type GetResponseFunc func(url, ac, key string, postData, postParams map[string]interface{}, cookies []*http.Cookie) (bool, interface{}, []*http.Cookie, error)

var responseFunc map[string]GetResponseFunc

func RegisterResponseFunc(name string, fun GetResponseFunc) {
	if responseFunc == nil {
		responseFunc = make(map[string]GetResponseFunc)
	}
	if _, ok := responseFunc[name]; ok {
		loger.Warnf("RegisterResponseFunc %s twice", name)
	}
	responseFunc[name] = fun
}

func UnRegisterResponseFunc(name string) {
	delete(responseFunc, name)
}

func ResetResponseFunc() {
	responseFunc = nil
}

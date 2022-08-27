package goper

import (
	"net/http"

	"github.com/bqqsrc/loger"
)

type RouterServeHTTP func(w http.ResponseWriter, r *http.Request)

var routerServeHttp map[string]RouterServeHTTP

func RegisterRouterServeHTTP(funcName string, method RouterServeHTTP) {
	if routerServeHttp == nil {
		routerServeHttp = make(map[string]RouterServeHTTP)
	}
	if _, ok := routerServeHttp[funcName]; ok {
		loger.Warnf("RegisterRouterServeHTTP %s twice", funcName)
	}
	routerServeHttp[funcName] = method
}

func UnRegisterRouterServeHTTP(funcName string) {
	delete(routerServeHttp, funcName)
}

func ResetRouterServeHTTP() {
	routerServeHttp = nil
}

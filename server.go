package goper

import (
	"fmt"

	//"log"
	"net/http"
	//	"github.com/bqqsrc/loger"
)

func runServerListen(domains Domains) error {
	funcName := "runServerListen"
	//loger.Debugf("domains is %s", domains)
	for key, value := range domains.Map {
		mux := http.NewServeMux()
		if value.length() > 1 {
			//TODO 如果有多个域名，所有路由收归一个函数，再由这个函数分发
			//如果多个域名，默认为重定向端口，不允许同时有其他的路由类型
			mux.Handle("/", &MutilDomainHandler{value})
		} else {
			//遍历所有路由，将静态路由存起来，将代理路由存起来，
			for host, routers := range value.Map {
				for url, router := range routers.Map {
					switch router.Type {
					case "static":
						fileHandle := http.StripPrefix(router.Static.Prefix, http.FileServer(http.Dir(router.Static.HttpDir)))
						mux.Handle(url, fileHandle)
						break
					case "proxy":
						mux.Handle(url, &ProxyHandler{router.Proxy})
						break
					case "func":
						methodName := router.Func.Method
						method, ok := routerServeHttp[methodName]
						if !ok {
							return RouterError(key, host, url, "router is not found, type is: %s, method is: %s, call %s to register a method", "func", methodName, "RegisterRouterServeHTTP")
						}
						mux.HandleFunc(url, method)
						break
					case "general":
						//loger.Debugf("testtest url %s, router %s", url, router)
						methodName := router.Func.Method
						if _, ok := generalResponses[methodName]; !ok {
							return RouterError(key, host, url, "router is not found, type is: %s, method is: %s, call %s to register a method", "general", methodName, "RegisterGeneralResponse")
						}
						mux.Handle(url, &FuncRouter{Method: router.Func.Method}) //, Cookie: router.Func.Cookie})
						break
					case "redirect":
						mux.Handle(url, &Gp_RedirectHandler{url, router.Redirect.Host})
						break
					default:
						return RouterError(key, host, url, "unsupport router type, type port is: %s, only can be static, proxy, func, redirect or general")
					}
				}
			}
		}
		host := fmt.Sprintf(":%d", key)
		proto := value.Proto
		switch proto.Name {
		case "http":
			go http.ListenAndServe(host, mux)
			break
		case "https":
			crt := proto.Crt
			key := proto.Key
			go http.ListenAndServeTLS(host, crt, key, mux)
			break
		default:
			return GoperError(funcName, "unsupport proto %s, listen is: %d, proto only can be http or https", proto, key)
		}
	}
	select {}
}

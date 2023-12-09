package goper

import (
	"fmt"

	//"log"
	"strings"
)

//一个路由组包括了一组路由名称和路由信息映射表
type Routers struct {
	Map map[string]Router //将路由和相关的信息以字典形式记录
}

func (r *Routers) equal(routers Routers) bool {
	for key, value := range r.Map {
		if tmpValue, ok := routers.Map[key]; !ok {
			return false
		} else {
			if !value.equal(tmpValue) {
				return false
			}
		}
	}
	return true
}
func (r *Routers) addRouter(router Router) error {
	url := router.Url
	if _, ok := r.Map[url]; ok {
		return fmt.Errorf("redeclare router, url is %s", url)
	}
	if router.Proto.equal(EmptyProtocol) {
		router.Proto = DefaultProtocol
	}
	r.Map[url] = router
	return nil
}
func (r *Routers) getRouter(url string) (Router, error) {
	router, ok := r.Map[url]
	ordinalUrl := url
	if !ok && !strings.HasSuffix(url, "/") {
		var build strings.Builder
		build.WriteString(url)
		build.WriteString("/")
		url = build.String()
		router, ok = r.Map[url]
	}
	if !ok {
		for {
			url = url[:len(url)-1]
			index := strings.LastIndex(url, "/")
			if index == -1 {
				break
			}
			url = url[:index]
			if router, ok = r.Map[url]; ok {
				break
			}
		}
	}
	if !ok {
		return EmptyRouter, fmt.Errorf("url not found in routers: url is %s", ordinalUrl)
	} else {
		return router, nil
	}
}

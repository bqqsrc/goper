package goper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bqqsrc/loger"
)

type MutilDomainHandler struct {
	domain Domain
}

func (m *MutilDomainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	loger.Debug(m.domain)
	fmt.Fprintf(w, "MutilDomainHandler %v", m.domain)
}

type ProxyHandler struct {
	Proxy ProxyRouter
}

//TODO 代理实现的类型：1.直接重定向，网页上的地址也会修改2.返回代理页面，3.不重定向，做一个中转，网页上的地址不修改
//TODO 实现路由重定向
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	loger.Debug(p.Proxy)
	//获取到host，url，获取到重定向的路径

}

type Gp_RedirectHandler struct {
	Url      string
	Redirect string
}

func (rd *Gp_RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !CheckRequest(r) {
		return
	}
	loger.Debugf("Test Host %s, url %s, proto %s, raw %s", r.Host, r.RequestURI, r.Proto, r.URL.RawQuery)
	url := r.RequestURI
	url = strings.Replace(url, rd.Url, "", 1)
	loger.Debugf("url is %s, rd.Url is %s", url, rd.Url)
	redirectTarget := rd.Redirect
	if !strings.HasSuffix(redirectTarget, "/") {
		redirectTarget = fmt.Sprintf("%s/", redirectTarget)
	}
	if strings.HasSuffix(url, "/") {
		url = url[1:]
	}
	redirectTarget = fmt.Sprintf("%s%s", redirectTarget, url)
	loger.Debugf("redirectTarget is %s", redirectTarget)
	w = WriteGoResponseHeader(w)
	http.Redirect(w, r, redirectTarget, http.StatusTemporaryRedirect)
}

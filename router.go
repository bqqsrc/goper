package goper

import (
	"errors"
	"strings"
)

//静态路由设置
type StaticRouter struct {
	HttpDir string `json:"http_dir"`
	Prefix  string `json:"prefix"`
}

var EmptyStaticRouter = StaticRouter{}
var DefaultStaticRouter = StaticRouter{"", ""}

//静态路由的配置都可以为空
func (s *StaticRouter) check() error {
	return nil
}
func (s *StaticRouter) equal(staticRouter StaticRouter) bool {
	return s.HttpDir == staticRouter.HttpDir && s.Prefix == staticRouter.Prefix
}

//函数路由
type FuncRouter struct {
	Method string `json:"method"`
	//	Cookie string `json:"cookie"`
}

func (f *FuncRouter) check() error {
	if f.Method == "" {
		return errors.New("method must not be empty")
	}
	return nil
}
func (f *FuncRouter) equal(funcRouter FuncRouter) bool {
	return f.Method == funcRouter.Method //&& f.Cookie == funcRouter.Cookie
}

//代理路由
type ProxyRouter struct {
	Host string `json:"host"`
}

func (p *ProxyRouter) check() error {
	if p.Host == "" {
		return errors.New("host must not be empty")
	}
	return nil
}
func (p *ProxyRouter) equal(proxyRouter ProxyRouter) bool {
	return p.Host == proxyRouter.Host
}

//重定向路由
type RedirectRouter struct {
	Host string `json:"host"`
}

func (r *RedirectRouter) check() error {
	if r.Host == "" {
		return errors.New("host must not be empty")
	}
	return nil
}
func (r *RedirectRouter) equal(redirectRouter RedirectRouter) bool {
	return r.Host == redirectRouter.Host
}

//一个路由信息
type Router struct {
	Url      string         `json:"url"`
	Type     string         `json:"type"`
	Proto    Protocol       `json:"protocol"`
	Static   StaticRouter   `json:"static"`
	Proxy    ProxyRouter    `json:"proxy"`
	Func     FuncRouter     `json:"func"`
	Redirect RedirectRouter `json:"redirect"`
}

var EmptyRouter = Router{}

//对路由的配置是否合法进行检查
func (r *Router) check() error {
	var build strings.Builder
	if r.Url == "" {
		build.WriteString("url must not be empty;\n")
	}
	if r.Type == "" {
		build.WriteString("type must not be empty;\n")
	}
	proto := r.Proto
	if err := proto.check(); err != nil {
		build.WriteString("protocol has error:\n")
		build.WriteString(err.Error())
		build.WriteString("\n")
	}
	switch r.Type {
	case "static":
		static := r.Static
		if err := static.check(); err != nil {
			build.WriteString("static has error:\n")
			build.WriteString(err.Error())
			build.WriteString("\n")
		}
		break
	case "proxy":
		proxy := r.Proxy
		if err := proxy.check(); err != nil {
			build.WriteString("proxy has error:\n")
			build.WriteString(err.Error())
			build.WriteString("\n")
		}
		break
	case "func", "general":
		fun := r.Func
		if err := fun.check(); err != nil {
			build.WriteString("func has error:\n")
			build.WriteString(err.Error())
			build.WriteString("\n")
		}
		break
	case "redirect":
		redirect := r.Redirect
		if err := redirect.check(); err != nil {
			build.WriteString("redirect has error:\n")
			build.WriteString(err.Error())
			build.WriteString("\n")
		}
		break
	default:
		build.WriteString("unsupport type: %s, type must be static, proxy, func, redirect or general")
		break
	}
	if build.Len() > 0 {
		return errors.New(build.String())
	}
	return nil
}
func (r *Router) equal(router Router) bool {
	return r.Url == router.Url && r.Type == router.Type && r.Proto.equal(router.Proto) &&
		r.Static.equal(router.Static) && r.Proxy.equal(router.Proxy) && r.Func.equal(router.Func) &&
		r.Redirect.equal(router.Redirect)
}

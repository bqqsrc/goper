package goper

import (
	"encoding/json"
	"errors"
	"strings"
	//"github.com/bqqsrc/loger"
)

type Domains struct {
	Map map[int]Domain
}

func (d *Domains) addRouter(port int, host string, proto Protocol, router Router) error {
	//loger.Debugf("before domainn addRouter, port is %d, host is %s, proto is %v, router is %v, domains is %v", port, host, proto, router, d)
	if d.Map == nil {
		d.Map = make(map[int]Domain)
	}
	domainMap := d.Map
	domains, ok := domainMap[port]
	//loger.Debug("test555 domain ", portServes, domains, ok)
	if !ok {
		domains = Domain{Protocol{}, make(map[string]Routers)}
	}
	if err := domains.addRouter(host, proto, router); err != nil {
		return RouterConfigError(port, host, "addRouter error, err: %s", err.Error())
	}
	domainMap[port] = domains
	d.Map = domainMap
	//loger.Debugf("after domainn addRouter, domains is %v", d)
	return nil
}

func (d *Domains) addRouters(port int, host string, proto Protocol, routers []Router) error {
	//loger.Debug("test333 AddRouters", port, host, routers)
	//loger.Debugf("domains addRouters begin, domains is %v", d)
	for _, value := range routers {
		//loger.Debug("test444 AddRouters", value, portServes)
		//loger.Debugf("domain addRouter begin, value is %v, domains is %v", value, d)
		if err := d.addRouter(port, host, proto, value); err != nil {
			return err
		}
		//loger.Debugf("domain addRouter after, domains is %v", d)
		//loger.Debug("test555 ", portServes)
	}
	return nil
}

func (d *Domains) getRouter(port int, host, url string) (Router, error) {
	if d.Map == nil || len(d.Map) == 0 {
		return EmptyRouter, RouterError(port, host, url, "there is no any router")
	}
	domain, ok := d.Map[port]
	if !ok {
		return EmptyRouter, RouterError(port, host, url, "there is no router, port is: %s", port)
	}
	router, err := domain.getRouter(host, url)
	if err != nil {
		err = RouterError(port, host, url, "get a router error, port is: %s, err is: %s", port, err.Error())
	}
	return router, err
}

func (d *Domains) Bytes2Config(filePath string, b []byte) error {
	funcName := "Domains.Bytes2Config"
	var serverList []ServeConfig
	err := json.Unmarshal(b, &serverList)
	if err != nil {
		return GoperError(funcName, "json.Unmarshal error, file is: %s, err is: %s", filePath, err)
	}
	//loger.Debugf("%s, after unmarshal, serverList is %v", funcName, serverList)
	if err := checkServerConfigList(serverList); err != nil {
		return GoperError(funcName, "configfile error, file is: %s, err is: %s", filePath, err.Error())
	}
	for _, value := range serverList {
		//loger.Debugf("%s, value is %v, before addRouters domains is %v", funcName, value, d)
		if err = d.addRouters(value.Listen, value.Host, value.Proto, value.Routers); err != nil {
			return err
		}
		//loger.Debugf("%s, after addRouters domains is %v", funcName, d)
	}
	return nil
}
func checkServerConfigList(serverList []ServeConfig) error {
	var build strings.Builder
	for _, value := range serverList {
		if err := value.check(); err != nil {
			build.WriteString(err.Error())
			build.WriteString("\n")
		}
	}
	if build.Len() > 0 {
		return errors.New(build.String())
	}
	return nil
}

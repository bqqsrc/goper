package goper

import (
	"fmt"
	"strings"
)

//一个域名组包括了所有域名和路由组映射表，以及域名协议
type Domain struct {
	Proto Protocol
	Map   map[string]Routers
}

func (d *Domain) equal(domain Domain) bool {
	if !d.Proto.equal(domain.Proto) {
		return false
	}
	for key, value := range d.Map {
		if tmpValue, ok := domain.Map[key]; !ok {
			return false
		} else {
			if !value.equal(tmpValue) {
				return false
			}
		}
	}
	return true
}

func (d *Domain) addRouter(host string, proto Protocol, router Router) error {
	routers, ok := d.Map[host]
	if !ok {
		routers = Routers{make(map[string]Router)}
	}
	if err := routers.addRouter(router); err != nil {
		return err
	}
	d.Map[host] = routers
	if !proto.equal(EmptyProtocol) {
		d.Proto = proto
	} else if d.Proto.equal(EmptyProtocol) {
		d.Proto = DefaultProtocol
	}
	return nil
}

func (d *Domain) getRouter(host, url string) (Router, error) {
	routers, ok := d.Map[host]
	ordinalHost := host
	if !ok {
		for {
			index := strings.Index(host, ".")
			if index == -1 {
				break
			}
			host = host[index:]
			host = fmt.Sprintf("*%s", host)
			if routers, ok = d.Map[host]; ok {
				break
			}
		}
	}
	if !ok {
		return EmptyRouter, fmt.Errorf("host is not found in routers, host is: %s", ordinalHost)
	}
	return routers.getRouter(url)
}

func (d *Domain) length() int {
	return len(d.Map)
}

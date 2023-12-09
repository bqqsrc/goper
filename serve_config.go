package goper

import "strings"

//一份服务的配置信息
type ServeConfig struct {
	Listen  int      `json:"listen"`
	Host    string   `json:"host"`
	Proto   Protocol `json:"protocol"`
	Routers []Router `json:"routers"`
}

//这只是对配置是否合法进行检查
func (s *ServeConfig) check() error {
	var build strings.Builder
	if s.Listen <= 0 {
		build.WriteString("listen must be large than 0;\n")
	}
	if s.Host == "" {
		build.WriteString("host must not be empty;\n")
	}
	if len(s.Routers) == 0 {
		build.WriteString("routers not contain any router;\n")
	} else {
		var tmpBuild strings.Builder
		for _, router := range s.Routers {
			if err := router.check(); err != nil {
				tmpBuild.WriteString(err.Error())
				tmpBuild.WriteString("\n")
			}
		}
		if tmpBuild.Len() > 0 {
			build.WriteString("some routers config error:\n")
			build.WriteString(tmpBuild.String())
		}
	}
	if build.Len() > 0 {
		return RouterConfigError(s.Listen, s.Host, build.String())
	} else {
		return nil
	}
}

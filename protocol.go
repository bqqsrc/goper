package goper

import (
	"errors"
	"strings"
)

type Protocol struct {
	Name string `json:"name"`
	Crt  string `json:"crt"`
	Key  string `json:"key"`
}

var EmptyProtocol = Protocol{}
var DefaultProtocol = Protocol{Name: "http"}

func (p *Protocol) check() error {
	var build strings.Builder
	if p.Name != "" {
		switch p.Name {
		case "http":
		case "https":
			if p.Crt == "" || p.Key == "" {
				build.WriteString("a https protocol must add crt and key")
			}
			break
		default:
			build.WriteString("unsupport protocol type: ")
			build.WriteString(p.Name)
			break
		}
	}
	if build.Len() > 0 {
		return errors.New(build.String())
	}
	return nil
}
func (p *Protocol) equal(proto Protocol) bool {
	return p.Name == proto.Name && p.Crt == proto.Crt && p.Key == proto.Key
}

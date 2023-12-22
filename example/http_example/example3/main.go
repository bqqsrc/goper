package main

import (
	"fmt"

	"github.com/bqqsrc/goper"
	"github.com/bqqsrc/goper/config"
	"github.com/bqqsrc/goper/core"
	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/http/hcore"
	"github.com/bqqsrc/goper/object"
)

type ExampleHandler struct {
	Ac int
}

func (h *ExampleHandler) Handler(c *http.Context) http.HttpPhase {
	switch h.Ac {
	case 1:
		c.Response.SetData("key", "HelloWorld, ac is 1")
	case 2:
		c.Response.SetData("key", "HelloWorld, ac is 2")
	default:
		c.Response.SetData("key", "HelloWorld, ac is default")
	}
	return http.HttpNext
}

type Example struct {
	http.HttpComponent
}

func (e *Example) CreateSrvConfig(key string) http.HttpCommands {
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"action",
				&ExampleHandler{},
			},
		},
	}
}

func (e *Example) MergeConfig(mainConf, srvConf, locConf, mthConf http.HttpConfigs) (object.ConfigValue, error) {
	if srvConf != nil && len(srvConf) > 0 {
		return srvConf[0].Value, nil
	}
	return nil, nil
}

func (e *Example) CreateHandler(dataConfig object.ConfigValue) (http.HttpHandler, http.HttpPhase, error) {
	if h, ok := dataConfig.(*ExampleHandler); ok {
		return h, http.HttpLogic, nil
	}
	return nil, http.HttpLogic, nil
}

var compts = []object.Componenter{
	&core.Core{},
	&config.Config{},
	&http.Http{},
	&hcore.HCore{},
	&Example{},
}

func main() {
	if err := goper.Launch(compts); err != nil {
		fmt.Printf("goper.Launch err: %v", err)
	}
}

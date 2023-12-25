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

type DemoHandler struct {
	Action int    `gson:"action"`
	Name   string `gson:"name"`
}

func (h *DemoHandler) Handler(c *http.Context) http.HttpPhase {
	fmt.Printf("DemoHandler action is %d, Name is %s", h.Action, h.Name)
	c.Response.SetData("action", h.Action)
	c.Response.SetData("name", h.Name)
	return http.HttpNext
}

type Demo struct {
	http.HttpComponent
}

func (d *Demo) CreateSrvConfig(key string) http.HttpCommands {
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"demo",
				&DemoHandler{},
			},
		},
	}
}
func (d *Demo) CreateMthConfig(key string) http.HttpCommands {
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"demo",
				&DemoHandler{},
			},
		},
	}
}

func (d *Demo) MergeConfig(mainConf, srvConf, locConf, mthConf http.HttpConfigs) (object.ConfigValue, error) {
	if mthConf != nil && len(mthConf) > 0 {
		return mthConf[0].Value, nil
	}
	if mainConf != nil && len(mainConf) > 0 {
		return mainConf[0].Value, nil
	}
	return nil, nil
}

func (d *Demo) CreateHandler(dataConfig object.ConfigValue) (http.HttpHandler, http.HttpPhase, error) {
	if h, ok := dataConfig.(*DemoHandler); ok {
		return h, http.HttpLogic, nil
	}
	return nil, http.HttpLogic, nil
}

var compts = []object.Componenter{
	&core.Core{},
	&config.Config{},
	&http.Http{},
	&hcore.HCore{},
	&Demo{},
}

func main() {
	if err := goper.Launch(compts); err != nil {
		fmt.Printf("goper.Launch err: %v", err)
	}
}

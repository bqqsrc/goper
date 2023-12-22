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

type Example struct {
	http.HttpComponent
	ac int
}

func (e *Example) CreateSrvConfig(key string) http.HttpCommands {
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"action",
				&e.ac,
			},
		},
	}
}

func (e *Example) Handler(c *http.Context) http.HttpPhase {
	switch e.ac {
	case 1:
		c.Response.SetData("key", "HelloWorld, ac is 1")
	case 2:
		c.Response.SetData("key", "HelloWorld, ac is 2")
	default:
		c.Response.SetData("key", "HelloWorld, ac is default")
	}
	return http.HttpNext
}

func (e *Example) CreateHandler(dataConfig object.ConfigValue) (http.HttpHandler, http.HttpPhase, error) {
	return e, http.HttpLogic, nil
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

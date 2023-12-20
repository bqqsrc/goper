//  Copyright (C) 晓白齐齐,版权所有.

package main

import (
	"fmt"

	"github.com/bqqsrc/goper"
	"github.com/bqqsrc/goper/config"
	"github.com/bqqsrc/goper/core"
	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/http/hcore"
	"github.com/bqqsrc/goper/http/hdatabase"
	"github.com/bqqsrc/goper/log"
	"github.com/bqqsrc/goper/object"
)

var compts = []object.Componenter{
	// compts的第一个一定是core.Core，不可改动
	&core.Core{},
	// config.Config模块用于从配置文件读取各组件的配置参数值，应当紧跟在core.Core之后
	&config.Config{},
	&log.Log{},
	&http.Http{},
	&hcore.HCore{},
	&hdatabase.HDatabase{},
}

func main() {
	if err := goper.Launch(compts); err != nil {
		fmt.Printf("err is %v, %T", err, err)
	} else {
		fmt.Println("err is nil")
	}
}

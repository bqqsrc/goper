//  Copyright (C) 晓白齐齐,版权所有.

package goper

import (
	"flag"
	"os"

	"github.com/bqqsrc/bqqg/errors"
	"github.com/bqqsrc/bqqg/file"
	"github.com/bqqsrc/goper/core"
	"github.com/bqqsrc/goper/log"
	"github.com/bqqsrc/goper/mode"
	"github.com/bqqsrc/goper/object"
)

// 启动
func Launch() error {
	var lErr errors.ErrorGroup
	var pCycle *object.Cycle
	defer func() {
		if lErr != nil {
			log.Errorln(lErr)
		}
	}()
	pCycle, lErr = analyArgs()

	if cnt := len(compts); cnt > 0 {
		if compt, ok := compts[0].(*core.Core); ok {
			if err := compt.Start(pCycle); err != nil {
				lErr = lErr.AddErrors(err)
			}
		} else {
			lErr = lErr.AddErrorf("type of compts[0] must be (*core.Core), but got (%T)", compt)
		}
	} else {
		lErr = lErr.AddErrorf("(Array: compts) must has some elements, but got %v", compts)
	}
	if lErr != nil {
		return lErr
	} else {
		log.Infof("goper Launch Success")
		log.LogConsoleln("goper Launch Success")
	}
	select {}
	return nil
}

// 分析传入参数
func analyArgs() (*object.Cycle, errors.ErrorGroup) {
	var lErr errors.ErrorGroup
	prefix, err := os.Getwd()
	if err != nil {
		lErr = lErr.AddErrors(err)
	}
	confFile, binFile := _GPConfFile, _GPBinFile

	flag.StringVar(&prefix, "prefix", prefix, "goper的安装路径，默认指定为可执行文件所在路径")
	flag.StringVar(&confFile, "conf-file", confFile, "goper的配置文件相对安装路径的文件（包含文件名）")
	flag.StringVar(&binFile, "bin-file", binFile, "goper的可执行文件相对安装路径的文件（包含文件名）")
	flag.Parse()

	//检查参数是否和法
	ok := false
	if confFile, ok = file.JoinPathIsFile(prefix, confFile); !ok {
		lErr = lErr.AddErrorf("conf-file is %s and must be a file, but it is not exist or not a file", confFile)
	}
	if binFile, ok = file.JoinPathIsFile(prefix, binFile); !ok && mode.EnvMode != mode.DEBUG {
		lErr = lErr.AddErrorf("bin-file is %s and must be a file, but it is not exist or not a file", binFile)
	}

	pc := &object.Cycle{
		Prefix: prefix,
		Compts: compts,
	}
	if lErr == nil {
		pc.ConfFile = confFile
		pc.BinFile = binFile
	}
	return pc, lErr
}

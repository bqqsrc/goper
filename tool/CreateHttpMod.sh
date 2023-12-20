#!/bin/bash

#  Copyright (C) 晓白齐齐,版权所有.

usage() {
	cat << END

使用参数说明：
参数使用规范：参数选项=参数值，或直接参数选项
例如：./CreateBaseMod.sh --help
    ./CreateBaseMod.sh --[arg]=[value] 

    参数选项 参数意义

    --help  显示帮助

    --skipAsk  如果脚本无参数，是否跳过询问获取参数的环节，如果不设置该参数默认不跳过

		--name  http组件名称

		--dir  组件模块相对goper根路径的相对路径
END
}

help=no 
skipAsk=no

Name=
Dir=

for option 
do 
	case "${option}" in
		-*=*) value=`echo "$option" | sed -e 's/[-_a-zA-Z0-9]*=//'` ;;
		*) value="" ;;
	esac

	case "${option}" in 
    --help | -help | --h | -h) help=yes ;;
		--skipAsk) skipAsk=yes ;;

		--name=*)
			Name=${value:-${Name}} 
			;;
		--dir=*)
			Dir=${value:-${Dir}} 
			;;

		*) 
			echo "Warn: 未知的参数 ${option}，查看是否遗漏了参数值" 
			;;
	esac
done

if [ ${help} = yes ]; then
    usage
    exit 0
fi

if [ $skipAsk = "no" ]; then
	if [ -z ${Name} ]; then 
		echo "输入新基本组件名称（大写开头名称）："
		read Name
	fi 
	if [ -z ${Dir} ]; then 
		echo "输入新基本组件路径："
		read Dir
	fi 
fi 

if [ -z ${Name} ]; then 
	echo "没有输入组件名称，未创建任何新组件" 
	exit 0
fi 

srcName=${Name,,}

if [ -z ${Dir} ]; then 
	Dir=${srcName}
fi 

cd ..
mkdir -p ${Dir}
cd ${Dir}
go mod init github.com/bqqsrc/goper/${srcName}
echo '//  Copyright (C) 晓白齐齐,版权所有.

package '${srcName}'

import (
	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/log"
	"github.com/bqqsrc/goper/object"
	"time"
)

type '${Name}' struct {
	http.HttpComponent
}

type '${srcName}'MainConf struct {
	Args1 string `gson:"Args1"`
	Args2 string `gson:"Args2"`
	Args3 string `gson:"Args3"`
}

type '${srcName}'SrvConf struct {
	Args1 string `gson:"Args1"`
	Args2 string `gson:"Args2"`
	Args3 string `gson:"Args3"`
}

type '${srcName}'LocConf struct {
	Args1 string `gson:"Args1"`
	Args2 string `gson:"Args2"`
	Args3 string `gson:"Args3"`
}

type '${srcName}'MthConf struct {
	Args1 string `gson:"Args1"`
	Args2 string `gson:"Args2"`
	Args3 string `gson:"Args3"`
}

type '${srcName}' struct {
	Args1 string
	Args2 string
	Args3 string
}

func (d *'${srcName}') Handler(c *http.Context) http.HttpPhase {
	//TODO 执行Handler逻辑
	return http.HttpNext
}

func (h *'${Name}') CreateMainConfig(key string) http.HttpCommands {
	//TODO 返回main级别的配置项
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"'${srcName}'main",
				&'${srcName}'MainConf{},
			},
		},
	}
}

func (h *'${Name}') CreateSrvConfig(key string) http.HttpCommands {
	//TODO 返回srv级别的配置项
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"'${srcName}'srv",
				&'${srcName}'SrvConf{},
			},
		},
	}
}

func (h *'${Name}') CreateLocConfig(key string) http.HttpCommands {
	//TODO 返回loc级别的配置项
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"'${srcName}'loc",
				&'${srcName}'LocConf{},
			},
		},
	}
}

func (h *'${Name}') CreateMthConfig(key string) http.HttpCommands {
	//TODO 返回mth级别的配置项
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"'${srcName}'mth",
				&'${srcName}'MthConf{},
			},
		},
	}
}

func (h *'${Name}') MergeConfig(mainConf, srvConf, locConf, mthConf http.HttpConfigs) (object.ConfigValue, error) {
	//TODO 合并4个级别的配置项，返回总配置项
	return &'${srcName}'{}, nil
}

func (h *'${Name}') CreateHandler(dataConfig object.ConfigValue) (http.HttpHandler, http.HttpPhase, error) {
	//TODO 返回Handler
	if dataConfig != nil {
		return dataConfig.(*'${srcName}'), http.HttpLogic, nil
	}
	return nil, http.HttpLogic, nil
}

func (h *'${Name}') Awake() error {
	log.Debugln("'${Name}' Awake")
	return nil
}
func (h *'${Name}') Start(pc *object.Cycle) error {
	log.Debugln("'${Name}' Start")
	return nil
}
func (h *'${Name}') Update(pc *object.Cycle) (time.Duration, error) {
	log.Debugln("'${Name}' Update")
	return 1000, nil
}
func (h *'${Name}') OnExit(pc *object.Cycle) error {
	log.Debugln("'${Name}' OnExit")
	return nil
}' > "${srcName}.go"

go get
go test

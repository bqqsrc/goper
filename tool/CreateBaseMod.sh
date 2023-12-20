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

		--name  基本组件名称

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
	"fmt"
	"time"

	"github.com/bqqsrc/goper/log"
	"github.com/bqqsrc/goper/object"
)

type '${Name}' struct {
	object.BaseComponent

	//TODO 添加要支持业务逻辑的字段
	val1 string
	val2 int
	val3 struct {
		Val1 string  `gson:"Val1"`
		Val2 float64 `gson:"Val2"`
	}
	val4 []float64
	val5 bool
}

func (t *'${Name}') CreateConfig(key string) []object.Command {
	//TODO 返回要从配置文件中解析的参数配置项和关注的配置
	// key为当前从配置文件获取读取到的键，有可能为空字符串，该字段一般在配置块时有用
	// 返回值：
	//     []object.Command，需要解析的参数配置项
	//     []string，关注的配置名称
	//     Componenter，关注这个配置键的组件，如果传nil则表示本组件关注，解析完会回调Componenter的Init
	// []object.Command记录了组件关注的键和解析的目标值的对应关系，配置文件读取到对应的键后会将其解析为所传的目标值
	// []string记录了组件关注的键，如果在[]object.Command找到对应的键，按照上面的方法处理
	// 如果找不到则会记录下这个键，后面读取到配置文件的对应的键时，不会解析为某个目标值，但是后面这个配置块的所有键都会调用到该组件的CreateConfig获取要解析的object.Command
	// 如果这两个返回值都是nil或者无值，表示该组件不需要任何的参数

	return []object.Command{
		{
			object.ConfigPair{"key1", // 要关注的key1
				&t.val1, // 要解析的值
			},
			false,
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				fmt.Println("found key", keyValue.Key, keyValue.Value, t.val1)
				return nil
			},
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				fmt.Println("done key", keyValue.Key, keyValue.Value, t.val1)
				return nil
			},
		},
		{object.ConfigPair{"key2",
			&t.val2,
		},
			false,
			nil,
			nil,
		},
		{object.ConfigPair{"key3",
			&t.val3,
		},
			false,
			nil,
			nil,
		},
		{object.ConfigPair{"key4",
			&t.val4,
		},
			false,
			nil,
			nil,
		},
		{object.ConfigPair{"key5",
			&t.val5,
		},
			false,
			nil,
			nil,
		},
	} //要关注的配置项字段
}

func (t *'${Name}') InitConfig(pc *object.Cycle) error {
	log.Debugln("'${Name}' InitConfig")
	return nil
}

func (t *'${Name}') Awake() error {
	log.Debugln("'${Name}' Awake")
	return nil
}
func (t *'${Name}') Start(pc *object.Cycle) error {
	log.Debugln("'${Name}' Start")
	return nil
}
func (t *'${Name}') Update(pc *object.Cycle) (time.Duration, error) {
	log.Debugln("'${Name}' Update")
	return 1000, nil
}
func (t *'${Name}') OnExit(pc *object.Cycle) error {
	log.Debugln("'${Name}' OnExit")
	return nil
}
func (t *'${Name}') GetKind() object.Kind {
	return object.ComptBase
}' > "${srcName}.go"

go get
go test
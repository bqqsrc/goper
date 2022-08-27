package goper

import (
	"strings"

	"gopkg.in/ini.v1"
)

var configPath string
var serverName = "[Goper/1.0 (Power by go 1.18)]"

//设置配置文件
func SetConfig(confPath string) {
	configPath = confPath
}

//获取服务配置文件列表
func GetConfigArr(section, key, sep string) ([]string, error) {
	funcName := "getServerCnf"
	if configPath == "" {
		return nil, GoperError(funcName, "configPath is empty, call SetConfig to resolve it")
	}
	dfIni, err := ini.Load(configPath)
	if err != nil {
		return nil, GoperError(funcName, "ini.Load(%s) error, err is %s", configPath, err)
	}
	config := dfIni.Section(section).Key(key).String()
	if config == "" {
		return nil, GoperError(funcName, "not found %s->%s in %s", section, key, configPath)
	}
	return strings.Split(config, sep), nil
}

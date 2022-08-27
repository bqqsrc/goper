package goper

import (
	"log"

	"github.com/bqqsrc/jsoner"
	//"github.com/bqqsrc/loger"
)

//开启监听服务
func ListenAndServe() error {
	serverConfig, err := GetConfigArr("server", "hostDir", ",")
	if err != nil {
		return err
	}
	//log.Debugf("serverConfig %s", serverConfig)
	domains := Domains{}
	//log.Debugf("domains is begin: %s", domains)
	if err = jsoner.ReadAllConfig(serverConfig, &domains); err != nil {
		return err
	}
	//log.Debugf("domains is 11 after: %v", domains)
	//执行监听
	if err = runServerListen(domains); err != nil {
		log.Fatalf("start serve error, error message is:\n %s", err)
		return err
	}
	return nil
}

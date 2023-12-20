//  Copyright (C) 晓白齐齐,版权所有.

package hdatabase

import (
	"github.com/bqqsrc/goper/http"

	"fmt"

	"github.com/bqqsrc/bqqg/databasehelper"
	"github.com/bqqsrc/goper/object"
	_ "github.com/go-sql-driver/mysql"
)

type HDatabase struct {
	http.HttpComponent
}

type database struct {
	Register string `gson:"register"`
	Driver   string `gson:"driver"`
	User     string `gson:"user"`
	Passwd   string `gson:"passwd"`
	Host     string `gson:"host"`
	Port     int    `gson:"port"`
	Name     string `gson:"name"`
	init     bool
}

func (d *database) Handler(c *http.Context) http.HttpPhase {
	if d.init {
		return http.HttpNext
	}
	conStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Passwd, d.Host, d.Port, d.Name)
	err := databasehelper.RegistController(d.Register, d.Driver, conStr)
	databasehelper.SetDefaultRegister(d.Register)
	if err != nil {
		c.Errors = c.Errors.AddErrors(err)
		return http.HttpError
	}
	return http.HttpNext
}

func (h *HDatabase) CreateSrvConfig(key string) http.HttpCommands {
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"database",
				&database{},
			},
		},
	}
}

func (h *HDatabase) MergeConfig(mainConf, srvConf, locConf, mthConf http.HttpConfigs) (object.ConfigValue, error) {
	if srvConf != nil && len(srvConf) > 0 {
		return srvConf[0].Value, nil
	}
	return nil, nil
}

func (h *HDatabase) CreateHandler(dataConfig object.ConfigValue) (http.HttpHandler, http.HttpPhase, error) {
	if dataConfig != nil {
		return dataConfig.(*database), http.HttpLogic, nil
	}
	return nil, http.HttpLogic, nil
}

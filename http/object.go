//  Copyright (C) 晓白齐齐,版权所有.

package http

import (
	"github.com/bqqsrc/goper/object"
)

type HttpPhase = uint8

type HttpHandler interface {
	Handler(*Context) HttpPhase
}

type HttpHandlerFunc func(*Context) HttpPhase

func (f HttpHandlerFunc) Handler(context *Context) HttpPhase {
	return f(context)
}

func HandleFunc(f func(*Context) HttpPhase) HttpHandlerFunc {
	return f
}

type HttpCommands = []object.Command
type HttpConfigs = []object.ConfigPair

type AllHttpConfigs struct {
	MainConfig []HttpConfigs
	SrvConfig  [][]HttpConfigs
	LocConfig  [][][]HttpConfigs
	MthConfig  [][][][]HttpConfigs
}

func (c *AllHttpConfigs) createConfig(mainIndex, srvIndex, locIndex, mthIndex, count int) {
	if c == nil {
		c = &AllHttpConfigs{}
	}
	if mainIndex >= 0 && c.MainConfig == nil {
		c.MainConfig = make([]HttpConfigs, count, count)
	}
	if srvIndex >= 0 {
		if c.SrvConfig == nil {
			c.SrvConfig = make([][]HttpConfigs, 0)
		}
		if c.LocConfig == nil {
			c.LocConfig = make([][][]HttpConfigs, 0)
		}
		if c.MthConfig == nil {
			c.MthConfig = make([][][][]HttpConfigs, 0)
		}
		index := len(c.SrvConfig)
		for index <= srvIndex {
			c.SrvConfig = append(c.SrvConfig, make([]HttpConfigs, count, count))
			c.LocConfig = append(c.LocConfig, nil)
			c.MthConfig = append(c.MthConfig, nil)
			index++
		}
	}
	if locIndex >= 0 {
		if c.LocConfig[srvIndex] == nil {
			c.LocConfig[srvIndex] = make([][]HttpConfigs, 0)
		}
		index := len(c.LocConfig[srvIndex])
		for index <= locIndex {
			c.LocConfig[srvIndex] = append(c.LocConfig[srvIndex], make([]HttpConfigs, count, count))
			c.MthConfig[srvIndex] = append(c.MthConfig[srvIndex], nil)
			index++
		}
	}
	if mthIndex >= 0 {
		if c.MthConfig[srvIndex][locIndex] == nil {
			c.MthConfig[srvIndex][locIndex] = make([][]HttpConfigs, 0)
		}
		index := len(c.MthConfig[srvIndex][locIndex])
		for index <= mthIndex {
			c.MthConfig[srvIndex][locIndex] = append(c.MthConfig[srvIndex][locIndex], make([]HttpConfigs, count, count))
			index++
		}
	}
}

type HttpComponenter interface {
	object.Componenter
	CreateMainConfig(string) HttpCommands
	CreateSrvConfig(string) HttpCommands
	CreateLocConfig(string) HttpCommands
	CreateMthConfig(string) HttpCommands
	MergeConfig(HttpConfigs, HttpConfigs, HttpConfigs, HttpConfigs) (object.ConfigValue, error)
	CreateHandler(object.ConfigValue) (HttpHandler, HttpPhase, error)
	setConfigs(*AllHttpConfigs)
}

type HttpComponent struct {
	object.SuperComponent
	Configs *AllHttpConfigs
}

func (hc *HttpComponent) CreateMainConfig(key string) HttpCommands { return nil }
func (hc *HttpComponent) CreateSrvConfig(key string) HttpCommands  { return nil }
func (hc *HttpComponent) CreateLocConfig(key string) HttpCommands  { return nil }
func (hc *HttpComponent) CreateMthConfig(key string) HttpCommands  { return nil }
func (hc *HttpComponent) MergeConfig(mainConf, srvConf, locConf, mthConf HttpConfigs) (object.ConfigValue, error) {
	return nil, nil
}
func (hc *HttpComponent) CreateHandler(config object.ConfigValue) (HttpHandler, HttpPhase, error) {
	return nil, HttpLogic, nil
}
func (hc *HttpComponent) GetKind() object.Kind { return object.ComptHttp }
func (hc *HttpComponent) setConfigs(configs *AllHttpConfigs) {
	if hc != nil {
		hc.Configs = configs
	}
}

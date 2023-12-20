//  Copyright (C) 晓白齐齐,版权所有.

package hcore

import (
	"fmt"
	httpp "net/http"
	"strings"

	"github.com/bqqsrc/bqqg/errors"
	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/log"
	"github.com/bqqsrc/goper/object"
)

var (
	anyMethods = []string{
		httpp.MethodGet, httpp.MethodPost, httpp.MethodPut, httpp.MethodPatch,
		httpp.MethodHead, httpp.MethodOptions, httpp.MethodDelete, httpp.MethodConnect,
		httpp.MethodTrace,
	}
)

type protocol struct {
	Name    string `gson:"name"`
	CrtFile string `gson:"crt"`
	KeyFile string `gson:"key"`
}

type listen struct {
	protocol
	port int
}

type Url struct {
	listen
	methods []string
	domain  string
	router  string
}

type HCore struct {
	http.HttpComponent
	defaultListen  listen
	defaultDomain  string
	defaultRouter  string
	defaultMethods []string
}

func (h *HCore) Awake() error {
	h.defaultListen.port = 80
	h.defaultListen.protocol.Name = "http"
	h.defaultRouter = "/"
	h.defaultMethods = []string{"any"}
	return nil
}

func getMainSrvConfig(l *listen, domain, router *string, methods *[]string) http.HttpCommands {
	ret := http.HttpCommands{
		{
			Config: object.ConfigPair{
				"listen",
				&l.port,
			},
		},
		{
			Config: object.ConfigPair{
				"protocol",
				&l.protocol,
			},
		},
		{
			Config: object.ConfigPair{
				"domain",
				domain,
			},
		},
	}
	ret = append(ret, getMainSrvLocConfig(router, methods)...)
	return ret
}

func getMainSrvLocConfig(router *string, methods *[]string) http.HttpCommands {
	ret := http.HttpCommands{
		{
			Config: object.ConfigPair{
				"url",
				router,
			},
		},
	}
	ret = append(ret, getMainSrvLocMthConfig(methods)...)
	return ret
}

func getMainSrvLocMthConfig(methods *[]string) http.HttpCommands {
	return http.HttpCommands{
		{
			Config: object.ConfigPair{
				"methods",
				methods,
			},
		},
	}
}

func (h *HCore) CreateMainConfig(key string) http.HttpCommands {
	return getMainSrvConfig(&h.defaultListen, &h.defaultDomain, &h.defaultRouter, &h.defaultMethods)
}

func (h *HCore) CreateSrvConfig(key string) http.HttpCommands {
	domain, router := "", "/"
	methods := []string{"any"}
	return getMainSrvConfig(&listen{port: 80, protocol: protocol{Name: "http"}}, &domain, &router, &methods)
}

func (h *HCore) CreateLocConfig(key string) http.HttpCommands {
	router := "/"
	methods := []string{"any"}
	return getMainSrvLocConfig(&router, &methods)
}

func (h *HCore) CreateMthConfig(key string) http.HttpCommands {
	methods := []string{"any"}
	return getMainSrvLocMthConfig(&methods)
}

func (h *HCore) MergeConfig(mainConf, srvConf, locConf, mthConf http.HttpConfigs) (object.ConfigValue, error) {
	url := Url{listen: h.defaultListen, methods: h.defaultMethods, domain: h.defaultDomain, router: h.defaultRouter}
	setConfig := func(httpConfs http.HttpConfigs, conf Url) Url {
		for _, value := range httpConfs {
			switch value.Key {
			case "listen":
				conf.listen.port = *value.Value.(*int)
			case "protocol":
				conf.listen.protocol = *value.Value.(*protocol)
			case "domain":
				conf.domain = *value.Value.(*string)
			case "url":
				conf.router = *value.Value.(*string)
			case "methods":
				conf.methods = *value.Value.(*[]string)
			}
		}
		return conf
	}
	url = setConfig(srvConf, url)
	url = setConfig(locConf, url)
	url = setConfig(mthConf, url)
	return url, nil
}

func (h *HCore) Start(pc *object.Cycle) error {
	servers, errGroup := mergeServers(pc)
	engines, errs := createTrees(servers)
	if errs != nil {
		errGroup = errGroup.AddErrors(errs)
	}
	errs = runHttpServe(engines)
	if errs != nil {
		errGroup = errGroup.AddErrors(errs)
	}
	if errGroup != nil {
		return errGroup
	}
	return nil
}

func mergeConfAndHandler(compt http.HttpComponenter, servers []server,
	mainConf, srvConf, locConf, mthConf http.HttpConfigs,
	serverIndex, routerIndex, restIndex, comptIndex int) ([]server, errors.ErrorGroup) {
	var errGroup errors.ErrorGroup
	if serverIndex < 0 || routerIndex < 0 || restIndex < 0 || comptIndex < 0 {
		return servers, errGroup.AddErrorf("serverIndex, routerIndex, restIndex, comptIndex should not less than 0, but got serverIndex: %d, routerIndex: %d, restIndex: %d, comptIndex: %d",
			serverIndex, routerIndex, restIndex, comptIndex)
	} else if comptConf, err := compt.MergeConfig(mainConf, srvConf, locConf, mthConf); err != nil {
		return servers, errGroup.AddErrors(err)
	} else {
		if comptIndex == 0 {
			conf, _ := comptConf.(Url)

			if servers == nil {
				servers = make([]server, 0, 4)
			}
			l := len(servers)
			if routerIndex == 0 {
				if serverIndex != l {
					return servers, errGroup.AddErrorf("routerIndex is %d, serverIndex is %d, so servers len should be %d, but got %d", routerIndex, serverIndex, serverIndex, l)
				}
				servers = append(servers, server{listen: conf.listen})
			} else {
				if serverIndex != l-1 {
					return servers, errGroup.AddErrorf("routerIndex is %d, serverIndex is %d, so servers len should be %d, but got %d", routerIndex, serverIndex, serverIndex+1, l)
				}
			}

			routers := servers[serverIndex].routers
			if routers == nil {
				routers = make([]router, 0, 4)
			}
			l = len(routers)
			if restIndex == 0 {
				if routerIndex != l {
					return servers, errGroup.AddErrorf("restIndex is %d, routerIndex is %d, so servers[%d].routers len should be %d, but got %d",
						restIndex, routerIndex, serverIndex, routerIndex, l)
				}
				routers = append(routers, router{domain: conf.domain, path: conf.router})
			} else {
				if routerIndex != l-1 {
					return servers, errGroup.AddErrorf("restIndex is %d, routerIndex is %d, so servers[%d].routers len should be %d, but got %d",
						restIndex, routerIndex, serverIndex, routerIndex+1, l)
				}
			}

			rests := routers[routerIndex].rests
			if rests == nil {
				rests = make([]rest, 0)
			}
			l = len(rests)
			if restIndex != l {
				return servers, errGroup.AddErrorf("comptIndex is %d, restIndex is %d, so servers[%d].routers[%d].rests len should be %d, but got %d",
					comptIndex, restIndex, serverIndex, routerIndex, restIndex, l)
			}
			rests = append(rests, rest{methods: conf.methods})
			routers[routerIndex].rests = rests
			servers[serverIndex].routers = routers
		}

		if handler, phase, err := compt.CreateHandler(comptConf); err != nil {
			return servers, errGroup.AddErrors(err)
		} else {
			if phase == http.HttpNext {
				return servers, errGroup.AddErrorf("CreateHandler should not return HttpNext, but %T return HttpNext", compt)
			}
			if handler != nil {
				servers[serverIndex].addHander(routerIndex, restIndex, phase, handler)
			}
		}
	}
	return servers, errGroup
}

func mergeServers(pc *object.Cycle) ([]server, errors.ErrorGroup) {
	httpIndexs := pc.KindIndexs[object.ComptHttp]
	var errGroup errors.ErrorGroup
	servers := make([]server, 0, 5)
	var allConfigs *http.AllHttpConfigs
	for index, httpIndex := range httpIndexs {
		if compt := pc.Compts[httpIndex]; compt != nil {
			if index == 0 {
				if httpCompt, ok := compt.(*HCore); !ok {
					errGroup.AddErrorf("first HttpComponent shoud be *hcore.HCore, but got %T", compt)
					return nil, errGroup
				} else {
					allConfigs = httpCompt.Configs
					if allConfigs == nil || allConfigs.MainConfig == nil || len(allConfigs.MainConfig) == 0 {
						return nil, nil
					}
				}
			}
			if httpCompt, ok := compt.(http.HttpComponenter); ok {
				var err errors.ErrorGroup
				if allConfigs.SrvConfig == nil || len(allConfigs.SrvConfig) == 0 {
					if servers, err = mergeConfAndHandler(httpCompt, servers, allConfigs.MainConfig[index], nil, nil, nil, 0, 0, 0, index); err != nil {
						errGroup = errGroup.AddErrors(err)
					}
				} else {
					for srvIndex := 0; srvIndex < len(allConfigs.SrvConfig); srvIndex++ {
						if allConfigs.LocConfig == nil || len(allConfigs.LocConfig) == 0 ||
							allConfigs.LocConfig[srvIndex] == nil || len(allConfigs.LocConfig[srvIndex]) == 0 {
							if servers, err = mergeConfAndHandler(httpCompt, servers, allConfigs.MainConfig[index], allConfigs.SrvConfig[srvIndex][index], nil, nil, srvIndex, 0, 0, index); err != nil {
								errGroup = errGroup.AddErrors(err)
							}
						} else {
							for locIndex := 0; locIndex < len(allConfigs.LocConfig[srvIndex]); locIndex++ {
								if allConfigs.MthConfig == nil || len(allConfigs.MthConfig) == 0 ||
									allConfigs.MthConfig[srvIndex] == nil || len(allConfigs.MthConfig[srvIndex]) == 0 ||
									allConfigs.MthConfig[srvIndex][locIndex] == nil || len(allConfigs.MthConfig[srvIndex][locIndex]) == 0 {
									if servers, err = mergeConfAndHandler(httpCompt, servers, allConfigs.MainConfig[index],
										allConfigs.SrvConfig[srvIndex][index], allConfigs.LocConfig[srvIndex][locIndex][index], nil, srvIndex, locIndex, 0, index); err != nil {
										errGroup = errGroup.AddErrors(err)
									}
								} else {
									for mthIndex := 0; mthIndex < len(allConfigs.MthConfig[srvIndex][locIndex]); mthIndex++ {
										if servers, err = mergeConfAndHandler(httpCompt, servers, allConfigs.MainConfig[index],
											allConfigs.SrvConfig[srvIndex][index], allConfigs.LocConfig[srvIndex][locIndex][index],
											allConfigs.MthConfig[srvIndex][locIndex][mthIndex][index], srvIndex, locIndex, mthIndex, index); err != nil {
											errGroup = errGroup.AddErrors(err)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return servers, errGroup
}

func createTrees(servers []server) (map[int]*engine, errors.ErrorGroup) {
	var errGroup errors.ErrorGroup
	ret := make(map[int]*engine)
	for serverIndex, serve := range servers {
		lis := serve.listen
		port := lis.port
		var eng *engine
		var ok bool
		if eng, ok = ret[port]; ok {
			if err := eng.checkListen(lis); err != nil {
				errGroup = errGroup.AddErrors(err)
				continue
			}
		} else {
			eng = new()
			if err := eng.setListen(lis); err != nil {
				errGroup = errGroup.AddErrors(err)
				continue
			}
			ret[port] = eng
		}
		for routerIndex, route := range serve.routers {
			methMap := make(map[string]uint8, 8)
			for _, rst := range route.rests {
				methods := rst.methods
			walkMth:
				for mthIndex, meth := range methods {
					meth = strings.ToUpper(meth)
					if meth == "ANY" {
						methods = append(anyMethods, methods[mthIndex+1:]...)
						goto walkMth
					}
					if _, ok := methMap[meth]; ok {
						errGroup = errGroup.AddErrorf("redeclare method in a router, serverIndex: %d, routerIndex: %d, method: %s", serverIndex+1, routerIndex+1, meth)
						continue
					} else {
						methMap[meth] = 1
					}
					if err := eng.addRouter(route.domain, route.path, meth, rst.handlers); err != nil {
						errGroup = errGroup.AddErrors(err)
					}
				}
			}
		}
	}
	return ret, errGroup
}

func runHttpServe(engines map[int]*engine) errors.ErrorGroup {
	var errGroup errors.ErrorGroup
	for port := range engines {
		go runAServe(port, engines[port])
	}
	return errGroup
}

func runAServe(port int, eng *engine) {
	switch eng.listenInfo.protocol.Name {
	case "http":
		if err := httpp.ListenAndServe(fmt.Sprintf(":%d", port), eng); err != nil {
			errMsg := fmt.Sprintf("ListenAndServe(:%d) error: %s", port, err)
			log.Errorln(errMsg)
			log.LogConsoleln("Error: ", errMsg)
		}
	case "https":
		if err := httpp.ListenAndServeTLS(fmt.Sprintf(":%d", port),
			eng.listenInfo.protocol.CrtFile, eng.listenInfo.protocol.KeyFile, eng); err != nil {
			errMsg := fmt.Sprintf("ListenAndServeTLS(:%d, %s, %s) error: %s",
				port, eng.listenInfo.protocol.CrtFile, eng.listenInfo.protocol.KeyFile, err)
			log.Errorln(errMsg)
			log.LogConsoleln("Error: ", errMsg)
		}
	default:
		errMsg := fmt.Sprintf("unsopport protocol type: %s, port: %d", eng.listenInfo.protocol.Name, port)
		log.Errorln(errMsg)
		log.LogConsoleln("Error: ", errMsg)
	}
}

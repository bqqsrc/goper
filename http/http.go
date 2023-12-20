//  Copyright (C) 晓白齐齐,版权所有.

package http

import (
	"github.com/bqqsrc/bqqg/errors"
	"github.com/bqqsrc/goper/object"
)

// 解析状态
type parseStatus = uint8

const (
	_ParseNone     parseStatus = 0
	_ParseHttpMain             = 1 << iota
	_ParseHttpSrv
	_ParseHttpLoc
	_ParseHttpMth
	_ParseHttp = _ParseHttpMain | _ParseHttpSrv | _ParseHttpLoc | _ParseHttpMth
)

type Http struct {
	object.BaseComponent
	createConfigFunc func(string) HttpCommands
}

func (h *Http) Awake() error {
	h.createConfigFunc = getCreateConfigFunc()
	return nil
}

func (h *Http) CreateConfig(key string) HttpCommands {
	return h.createConfigFunc(key)
}

type comptComdInfo struct {
	compt      HttpComponenter
	comd       *object.Command
	comptIndex int
}

func getCreateConfigFunc() func(string) HttpCommands {
	var errGroup errors.ErrorGroup
	var parseState parseStatus = _ParseNone
	var pCycle *object.Cycle = nil
	var httpComptIndexs []int
	var httpComptCount, srvCount, locCount, mthCount, comptMainIndex, comptSrvIndex, comptLocIndex, comptMthIndex int
	var mainComptComds, srvComptComds, locComptComds, mthComptComds, keyWordComptComds map[string]comptComdInfo
	allHttpConfigs := &AllHttpConfigs{}

	httpFoundCall := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		pCycle = pc
		errGroup = nil
		parseState |= _ParseHttpMain
		httpComptIndexs = pc.KindIndexs[object.ComptHttp]
		httpComptCount = len(httpComptIndexs)

		allHttpConfigs.createConfig(0, -1, -1, -1, httpComptCount)

		for _, index := range httpComptIndexs {
			if compt := pc.Compts[index]; compt != nil {
				if hcompt, ok := compt.(HttpComponenter); ok {
					hcompt.setConfigs(allHttpConfigs)
				} else {
					errGroup = errGroup.AddErrorf("%T is defined as a ComptHttp, but it does not implement all method of HttpComponenter", compt)
				}
			}
		}

		return nil
	}

	httpDoneCall := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		pCycle = nil
		parseState ^= _ParseHttpMain
		if errGroup != nil {
			return errGroup
		}
		return nil
	}

	serverFoundCall := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		parseState |= _ParseHttpSrv
		allHttpConfigs.createConfig(0, srvCount, -1, -1, httpComptCount)
		return nil
	}

	serverDoneCall := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		parseState ^= _ParseHttpSrv
		srvCount++
		locCount = 0
		mthCount = 0
		return nil
	}

	locationFoundCall := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		parseState |= _ParseHttpLoc
		allHttpConfigs.createConfig(0, srvCount, locCount, -1, httpComptCount)
		return nil
	}

	locationDoneCall := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		parseState ^= _ParseHttpLoc
		locCount++
		mthCount = 0
		return nil
	}

	methodFoundCall := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		parseState |= _ParseHttpMth
		allHttpConfigs.createConfig(0, srvCount, locCount, mthCount, httpComptCount)
		return nil
	}

	methodDoneCall := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		parseState ^= _ParseHttpMth
		mthCount++
		return nil
	}

	return func(key string) HttpCommands {
		if parseState&_ParseHttp == 0 {
			return []object.Command{
				{
					object.ConfigPair{
						"http",
						nil,
					}, false, httpFoundCall, httpDoneCall,
				},
			}

		} else {
			switch {
			case key == "server" && parseState&_ParseHttp == _ParseHttpMain:
				return []object.Command{
					{
						object.ConfigPair{
							"server",
							nil,
						}, false, serverFoundCall, serverDoneCall,
					},
				}

			case key == "location" && parseState&_ParseHttp == _ParseHttpMain|_ParseHttpSrv:
				return []object.Command{
					{
						object.ConfigPair{
							"location",
							nil,
						}, false, locationFoundCall, locationDoneCall,
					},
				}

			case key == "method" && parseState&_ParseHttp == _ParseHttp^_ParseHttpMth:
				return []object.Command{
					{
						object.ConfigPair{
							"method",
							nil,
						}, false, methodFoundCall, methodDoneCall,
					},
				}

			default:
				comptComd, ok := keyWordComptComds[key]
				if ok {
					return HttpCommands{*comptComd.comd}
				}

				var ret HttpCommands
				var errs errors.ErrorGroup
				var keywords []comptComdInfo

				switch {
				case parseState&_ParseHttpMth != 0:
					ret, keywords, allHttpConfigs.MthConfig[srvCount][locCount][mthCount], comptMthIndex, mthComptComds, errs =
						getConfig(key, parseState, mthComptComds, comptMthIndex, httpComptIndexs, pCycle.Compts, allHttpConfigs.MthConfig[srvCount][locCount][mthCount])
				case parseState&_ParseHttpLoc != 0:
					ret, keywords, allHttpConfigs.LocConfig[srvCount][locCount], comptLocIndex, locComptComds, errs =
						getConfig(key, parseState, locComptComds, comptLocIndex, httpComptIndexs, pCycle.Compts, allHttpConfigs.LocConfig[srvCount][locCount])
				case parseState&_ParseHttpSrv != 0:
					ret, keywords, allHttpConfigs.SrvConfig[srvCount], comptSrvIndex, srvComptComds, errs =
						getConfig(key, parseState, srvComptComds, comptSrvIndex, httpComptIndexs, pCycle.Compts, allHttpConfigs.SrvConfig[srvCount])
				case parseState&_ParseHttpMain != 0:
					ret, keywords, allHttpConfigs.MainConfig, comptMainIndex, mainComptComds, errs =
						getConfig(key, parseState, mainComptComds, comptMainIndex, httpComptIndexs, pCycle.Compts, allHttpConfigs.MainConfig)
				}

				if keywords != nil {
					allComptComds := [...]map[string]comptComdInfo{keyWordComptComds, mainComptComds, srvComptComds, locComptComds, mthComptComds}
					if keyWordComptComds, errs = isKeywordRedefined(keywords, keyWordComptComds, allComptComds[:]); errs != nil {
						errs = errs.AddErrors(errs)
					}
				}

				if errs != nil {
					errGroup = errGroup.AddErrors(errs)
					return nil
				}
				return ret
			}
		}
	}
}

func getConfig(key string, state parseStatus, comptComds map[string]comptComdInfo, comptIndex int, httpIndexs []int,
	compts []object.Componenter, config []HttpConfigs) (ret HttpCommands, keywords []comptComdInfo, newConfig []HttpConfigs,
	newComptIndex int, newComptComds map[string]comptComdInfo, errs errors.ErrorGroup) {
	newComptIndex = comptIndex
	newComptComds = comptComds
	newConfig = config
	comptIndex = -1
	if newComptComds == nil {
		newComptComds = make(map[string]comptComdInfo)
	}
	if comptComd, ok := newComptComds[key]; ok {
		if comptComd.comd != nil {
			ret = HttpCommands{*comptComd.comd}
			comptComd.comd = nil
			newComptComds[key] = comptComd
			comptIndex = comptComd.comptIndex
		} else {
			comds := createCommands(key, state, comptComd.compt)
			var found bool
			if found, _, ret, _, newComptComds = findCommand(key, comds, newComptComds, comptComd.compt, comptComd.comptIndex); found {
				comptIndex = comptComd.comptIndex
			}
		}
	} else {
		httpCnt := len(httpIndexs)
		for newComptIndex < httpCnt {
			com := compts[httpIndexs[newComptIndex]]
			if compt, ok := com.(HttpComponenter); ok {
				comds := createCommands(key, state, compt)
				newComptIndex++
				if err := isKeyRedefined(compt, comds, comptComds, state); err != nil {
					errs = errs.AddErrors(err)
					continue
				}
				var keyword []comptComdInfo
				var found, isKeyWord bool
				found, isKeyWord, ret, keyword, newComptComds = findCommand(key, comds, newComptComds, compt, newComptIndex-1)
				if keyword != nil && len(keyword) > 0 {
					if keywords == nil {
						keywords = keyword
					} else {
						keywords = append(keywords, keyword...)
					}
				}
				if found {
					if !isKeyWord {
						comptIndex = newComptIndex - 1
					}
					break
				}
			}
		}
	}
	if comptIndex >= 0 {
		newConfig[comptIndex] = setConfig(ret, newConfig[comptIndex])
	}
	return
}

func createCommands(key string, state parseStatus, compt HttpComponenter) HttpCommands {
	if createFunc := getCreateMethod(compt, state); createFunc != nil {
		return createFunc(key)
	}
	return nil
}

func findCommand(key string, comds HttpCommands, comptComds map[string]comptComdInfo, compt HttpComponenter, comptIndex int) (bool, bool, HttpCommands, []comptComdInfo, map[string]comptComdInfo) {
	var found, isKeywords bool
	var ret HttpCommands
	var keywords []comptComdInfo
	for index, comd := range comds {
		if comd.KeyWord {
			if keywords == nil {
				keywords = make([]comptComdInfo, 0)
			}
			keywords = append(keywords, comptComdInfo{compt, &comds[index], comptIndex})
		}
		if comd.Config.Key == key {
			found = true
			ret = HttpCommands{comds[index]}
			if comd.KeyWord {
				isKeywords = true
			} else {
				comptComds[key] = comptComdInfo{compt, nil, comptIndex}
			}
		} else if !comd.KeyWord {
			comptComds[comd.Config.Key] = comptComdInfo{compt, &comds[index], comptIndex}
		}
	}
	return found, isKeywords, ret, keywords, comptComds
}

func isKeyRedefined(compt HttpComponenter, comds HttpCommands, comptComds map[string]comptComdInfo, state parseStatus) errors.ErrorGroup {
	var errs errors.ErrorGroup
	keyMap := make(map[string]uint8)
	configName := ""
	switch {
	case state&_ParseHttpMth != 0:
		configName = "CreateMthConfig"
	case state&_ParseHttpLoc != 0:
		configName = "CreateLocConfig"
	case state&_ParseHttpSrv != 0:
		configName = "CreateSrvConfig"
	case state&_ParseHttpMain != 0:
		configName = "CreateMainConfig"
	}
	for _, comd := range comds {
		if err := comd.IsValid(); err != nil {
			errs = errs.AddErrors(err)
			continue
		}
		key := comd.Config.Key
		if _, ok := keyMap[key]; ok {
			errs = errs.AddErrorf("redeclare config key %s in %T.%s(string) HttpCommands", key, compt, configName)
			continue
		} else {
			keyMap[key] = 1
		}
		if comptcomd, ok := comptComds[key]; ok && comptcomd.compt != compt {
			errs = errs.AddErrorf("config of [%s(string) HttpCommands] in http component should be different, but %T and %T has same config key %s in [%s(string) HttpCommands]",
				configName, comptcomd.compt, compt, key, configName)
			continue
		}
	}
	return errs
}

func isKeywordRedefined(keywords []comptComdInfo, keyWordComptComds map[string]comptComdInfo, allComptComdInfo []map[string]comptComdInfo) (map[string]comptComdInfo, errors.ErrorGroup) {
	if keywords == nil || len(keywords) == 0 {
		return keyWordComptComds, nil
	}
	var errs errors.ErrorGroup
	for index, checkInfo := range keywords {
		key := checkInfo.comd.Config.Key
		err := false
		for _, targetInfos := range allComptComdInfo {
			if comptcomd, ok := targetInfos[key]; ok {
				errs = errs.AddErrorf("%s is declare as a keyword in %T, and redeclare in %T", key, comptcomd.compt, checkInfo.compt)
				err = true
				break
			}
		}
		if !err {
			if keyWordComptComds == nil {
				keyWordComptComds = make(map[string]comptComdInfo)
			}
			keyWordComptComds[key] = keywords[index]
		}
	}
	return keyWordComptComds, errs
}

func getCreateMethod(compt HttpComponenter, state parseStatus) func(string) HttpCommands {
	switch {
	case state&_ParseHttpMth != 0:
		return compt.CreateMthConfig
	case state&_ParseHttpLoc != 0:
		return compt.CreateLocConfig
	case state&_ParseHttpSrv != 0:
		return compt.CreateSrvConfig
	case state&_ParseHttpMain != 0:
		return compt.CreateMainConfig
	default:
		return nil
	}
}

func setConfig(comds HttpCommands, config HttpConfigs) HttpConfigs {
	if comds != nil && len(comds) > 0 {
		for _, comd := range comds {
			if comd.Config.Value != nil {
				if config == nil {
					config = make(HttpConfigs, 0)
				}
				config = append(config, comd.Config)
			}
		}
	}
	return config
}

//  Copyright (C) 晓白齐齐,版权所有.

package config

import (
	"fmt"

	"github.com/bqqsrc/bqqg/errors"

	// "time"
	"io/ioutil"

	"github.com/bqqsrc/bqqg/file"
	"github.com/bqqsrc/goper/log"
	"github.com/bqqsrc/goper/object"
	"github.com/bqqsrc/gson"
)

type Config struct {
	object.BaseComponent
	includeFiles []string
	initConfFunc func(*object.Cycle, string) error
}

func (c *Config) Awake() error {
	c.initConfFunc = getParseConfig()
	return nil
}

func (c *Config) CreateConfig(key string) []object.Command {
	log.Debugf("Config CreateConfig(%s), config key: include", key)
	return []object.Command{{object.ConfigPair{"include", &c.includeFiles}, true, nil, c.include}}
}

func (c *Config) Start(pc *object.Cycle) error {
	log.Debugln("Config Start(*object.Cycle)")
	log.Infof("confFile is : %s", pc.ConfFile)
	if pc.ConfFile == "" {
		return fmt.Errorf("Config Start Error, ConfFile of pc is empty")
	}
	var errGroup errors.ErrorGroup
	if err := c.initConfFunc(pc, pc.ConfFile); err != nil {
		errGroup = errGroup.AddErrors(err)
	}
	// baseIndexs := pc.KindIndexs[object.ComptBase]
	for _, index := range pc.KindIndexs[object.ComptBase] {
		if compt := pc.Compts[index]; compt != nil {
			if baseCompt, ok := compt.(object.BaseComponenter); ok {
				if err := baseCompt.InitConfig(pc); err != nil {
					errGroup = errGroup.AddErrors(err)
				}
			} else {
				errGroup = errGroup.AddErrorf("%T is defined as a ComptBase, but it does not implement all method of BaseComponenter", compt)
			}
		}
	}
	if errGroup != nil {
		return errGroup
	}
	return nil
}

type comptComdInfo struct {
	compt object.BaseComponenter
	comd  object.Command
}

// 解析一个配置文件
func getParseConfig() func(*object.Cycle, string) error {
	//所有配置文件公用的
	var baseIndexs []int
	baseLen := -1
	var baseComptComds, keyWordComptComds map[string]comptComdInfo

	//一个配置文件公用的
	var blockCompt object.BaseComponenter
	keyStack := make([]string, 0, 4)
	comdStack := make([]object.Command, 0, 4)

	return func(pc *object.Cycle, filePath string) error {
		log.Debugf("parseConfig: %s", filePath)
		if baseLen < 0 {
			baseIndexs = pc.KindIndexs[object.ComptBase]
			baseLen = len(baseIndexs)
			baseComptComds = make(map[string]comptComdInfo, baseLen>>1)
			keyWordComptComds = make(map[string]comptComdInfo, baseLen>>1)
		}
		index := 0

		//TODO 这个地方改为可以根据密码进行解密读取
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read conf-file err: %s, conf-file: %s", err, filePath)
		}
		var errGroup errors.ErrorGroup

		keyEventCallback := func() func(*gson.Decoder, *gson.Lexer, bool) bool {
			var targetCompt object.BaseComponenter
			var targetComd object.Command
			return func(d *gson.Decoder, l *gson.Lexer, isFound bool) bool {
				key := l.String()
				if isFound {
					parseComptComd := func(compt object.BaseComponenter, comd object.Command) {
						if comd.ConfigFoundCallback != nil {
							if err = comd.ConfigFoundCallback(object.ConfigPair{key, comd.Config.Value}, pc); err != nil {
								errGroup = errGroup.AddErrors(err)
							}
						}
						if comd.Config.Value != nil {
							if err = d.SetAnyTarget(comd.Config.Value, true); err != nil {
								errGroup = errGroup.AddErrorf("SetAnyTarget Error: %s, Component: %v, key: %s", err, compt, key)
							}
						} else if blockCompt == nil {
							blockCompt = compt
						}
					}
					keyStack = append(keyStack, key)
					if comptComd, ok := keyWordComptComds[key]; ok { // 优先给到关键字
						targetCompt = comptComd.compt
						targetComd = comptComd.comd
						parseComptComd(targetCompt, targetComd)
					} else if blockCompt != nil { // 再给到代码块
						targetCompt = blockCompt
						if comds := targetCompt.CreateConfig(key); len(comds) > 0 {
							targetComd = comds[0]
						} else {
							targetComd = object.Command{}
						}
						parseComptComd(targetCompt, targetComd)
					} else if comptComd, ok = baseComptComds[key]; ok {
						targetCompt = comptComd.compt
						targetComd = comptComd.comd
						parseComptComd(targetCompt, targetComd)
					} else {
						var toBreak bool
						for index < baseLen {
							baseIndex := baseIndexs[index]
							index++
							var compt object.BaseComponenter
							if comp := pc.Compts[baseIndex]; comp != nil {
								if compt, ok = comp.(object.BaseComponenter); !ok {
									errGroup = errGroup.AddErrorf("%T is defined as a ComptBase, but it does not implement all method of BaseComponenter", comp)
									continue
								}
							}
							if commands := compt.CreateConfig(key); commands != nil {
								for _, command := range commands {
									if err := command.IsValid(); err != nil {
										errGroup.AddErrors(err)
										continue
									}
									name := command.Config.Key
									if comptComd, ok := keyWordComptComds[name]; ok {
										if comptComd.compt == compt {
											errGroup.AddErrorf("redeclare config key %s in %T", key, compt)
										} else {
											errGroup.AddErrorf("config of base component should be different, %T and %T has same config key %s", comptComd.compt, compt, name)
										}
										continue
									}
									if comptComd, ok := baseComptComds[name]; ok {
										if comptComd.compt == compt {
											errGroup.AddErrorf("redeclare config key %s in %T", key, compt)
										} else {
											errGroup.AddErrorf("config of base component should be different, %T and %T has same config key %s", comptComd.compt, compt, name)
										}
										continue
									}
									if command.KeyWord {
										keyWordComptComds[name] = comptComdInfo{compt, command}
									} else {
										baseComptComds[name] = comptComdInfo{compt, command}
									}
									if name == key {
										targetCompt = compt
										targetComd = command
										toBreak = true
									}
								}
							}
							if toBreak {
								parseComptComd(targetCompt, targetComd)
								break
							}
						}
						if !toBreak {
							targetComd = object.Command{}
						}
					}
					comdStack = append(comdStack, targetComd)
				} else {
					l := len(comdStack)
					if l > 0 {
						targetComd = comdStack[l-1]
						comdStack = comdStack[:l-1]
					} else {
						targetComd = object.Command{}
					}
					l = len(keyStack)
					if l > 0 {
						keyStack = keyStack[:l-1]
					}
					if targetComd.ConfigDoneCallback != nil {
						if err := targetComd.ConfigDoneCallback(object.ConfigPair{key, targetComd.Config.Value}, pc); err != nil {
							errGroup = errGroup.AddErrors(err)
						}
					}
					targetCompt = nil
					targetComd = object.Command{}
					if l-1 == 0 {
						blockCompt = nil
					}
				}
				return true
			}
		}()
		if err = gson.DecodeData(data, keyEventCallback, nil, nil); err != nil {
			errGroup.AddErrors(err)
		}
		if errGroup != nil {
			return errGroup
		}
		return nil
	}
}

func (c *Config) include(keyValue object.ConfigPair, pc *object.Cycle) error {
	log.Debugf("include file: %v", *keyValue.Value.(*[]string))
	if c.includeFiles == nil || len(c.includeFiles) == 0 {
		return nil
	}
	var errGroup errors.ErrorGroup
	for _, includeFile := range c.includeFiles {
		if ok := file.IsFile(includeFile); !ok {
			errGroup = errGroup.AddErrorf("include file %s and must be a file, but it is not exist or not a file", includeFile)
		} else {
			if err := c.initConfFunc(pc, includeFile); err != nil {
				errGroup = errGroup.AddErrors(err)
			}
		}
	}
	if errGroup != nil {
		return errGroup
	}
	return nil
}

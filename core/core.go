//  Copyright (C) 晓白齐齐,版权所有.

package core

import (
	"time"

	"github.com/bqqsrc/bqqg/errors"
	"github.com/bqqsrc/goper/log"
	"github.com/bqqsrc/goper/object"
)

type Core struct {
	object.BaseComponent
	frequency int // 执行Update的帧率，单位毫秒
}

func (c *Core) CreateConfig(key string) []object.Command {
	log.Debugf("Core CreateConfig(%s), config key: frequency", key)
	return []object.Command{{object.ConfigPair{"frequency", &c.frequency}, false, nil, nil}}
}

func (c *Core) Start(pc *object.Cycle) error {
	log.Debugln("Core Start(*object.Cycle)")
	errGroup := awake(pc)
	if err := start(pc); err != nil {
		errGroup = errGroup.AddErrors(err)
	}
	nextDuration, err := update(pc, true, c.frequency)
	if err != nil {
		errGroup = errGroup.AddErrors(err)
	}
	if errGroup != nil {
		return errGroup
	}
	if nextDuration >= 0 {
		log.Debugf("nextDuration: %d, need to update some components", nextDuration)
		updateAll := func() {
			for {
				time.Sleep(time.Duration(nextDuration))
				nextDuration, _ = update(pc, false, c.frequency)
			}
		}
		go updateAll()
	}
	return nil
}

func awake(pc *object.Cycle) errors.ErrorGroup {
	log.Debugln("Awake all components")
	var errGroup errors.ErrorGroup
	pc.KindIndexs = make(map[object.Kind][]int, object.KindCount)
	comptsCnt := len(pc.Compts)
	//TODO这个for循环可以用并发处理
	//Awake与顺序无关，因此可以用并发执行
	for index, compt := range pc.Compts {
		if compt != nil {
			if err := compt.Awake(); err != nil {
				errGroup = errGroup.AddErrors(err)
			}
		}
		kind := compt.GetKind()
		kindArr, ok := pc.KindIndexs[kind]
		if !ok {
			kindArr = make([]int, 0, comptsCnt>>3)
		}
		pc.KindIndexs[kind] = append(kindArr, index)
	}
	return errGroup
}

func start(pc *object.Cycle) errors.ErrorGroup {
	log.Debugln("Start all components except *core.Core")
	var errGroup errors.ErrorGroup
	comptsCnt := len(pc.Compts)
	index := 1
	//Start和Update与顺序有关，因此不能用并发处理
	for index < comptsCnt {
		if err := pc.Compts[index].Start(pc); err != nil {
			errGroup = errGroup.AddErrors(err)
		}
		index++
	}
	return errGroup
}

func update(pc *object.Cycle, fitst bool, frequency int) (int, errors.ErrorGroup) {
	log.Debugf("Update all components, frequency: %d", frequency)
	var errGroup errors.ErrorGroup
	start := time.Now()
	var updateIndexs []int
	if fitst {
		for index, compt := range pc.Compts {
			if compt != nil {
				if duration, err := compt.Update(pc); err != nil {
					errGroup = errGroup.AddErrors(err)
				} else {
					if duration > 0 {
						if updateIndexs == nil {
							updateIndexs = make([]int, 0, len(pc.Compts)>>3)
						}
						updateIndexs = append(updateIndexs, index)
					}
				}
			}
		}
		if len(updateIndexs) == 0 {
			return -1, errGroup
		}
	} else {
		for index := range updateIndexs {
			if _, err := pc.Compts[updateIndexs[index]].Update(pc); err != nil {
				errGroup = errGroup.AddErrors(err)
			}
		}
	}
	elapsed := time.Since(start)
	if frequency <= 0 {
		frequency = 1000
	}
	nextDuration := frequency*1000000 - int(elapsed)
	if nextDuration < 0 {
		nextDuration = 0
	}
	return nextDuration, errGroup
}

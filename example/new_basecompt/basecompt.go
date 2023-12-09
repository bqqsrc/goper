//  Copyright (C) 晓白齐齐,版权所有.

package basecompt

import (
	"fmt"

	"github.com/bqqsrc/goper/object"
)

type HelloWord struct {
	object.BaseComponent
	action int
	name   string
}

func (h *HelloWord) CreateConfig(key string) []object.Command {
	return []object.Command{
		{
			object.ConfigPair{"action", &h.action},
			false, nil, nil,
		},
		{
			object.ConfigPair{"name", &h.name},
			false,
			func(conf object.ConfigPair, pc *object.Cycle) error {
				fmt.Printf("found a key: %s", conf.Key)
				return nil
			},
			func(conf object.ConfigPair, pc *object.Cycle) error {
				fmt.Printf("parse finish, key: %s, name: %s ", conf.Key, h.name)
				return nil
			},
		},
	}
}

func (h *HelloWord) Start(pc *object.Cycle) error {
	if h.action > 0 {
		fmt.Printf("HelloWord, my name is %s, I got a num larger than 0", h.name)
	} else if h.action < 0 {
		fmt.Printf("HelloWord, my name is %s, I got a num less than 0", h.name)
	} else {
		fmt.Printf("HelloWord, my name is %s, I got a num equal 0", h.name)
	}
	return nil
}

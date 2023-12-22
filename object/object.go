//  Copyright (C) 晓白齐齐,版权所有.

package object

import (
	"fmt"
	"time"
)

// 组件的参数值类型的定义
type ConfigValue = any

// 组件类型的定义
type Kind = uint8

type ConfigPair struct {
	Key   string      // 参数名称
	Value ConfigValue // 参数值，这个值传入时只能是指针，可以为空，如果为nil表示关注这个key，但是不是解析为一个值，而是进行配置块解析
}

// 组件的配置指令定义
type Command struct {
	Config              ConfigPair
	KeyWord             bool                           // 是否注册为关键字
	ConfigFoundCallback func(ConfigPair, *Cycle) error // 找到参数的回调方法
	ConfigDoneCallback  func(ConfigPair, *Cycle) error // 参数解析完的回到方法
}

func (c Command) IsValid() error {
	if c.Config.Key == "" {
		return fmt.Errorf("Key of Command should be not empty")
	}
	if c.KeyWord && c.Config.Value == nil {
		return fmt.Errorf("Value of KeyWord-Command should be not nil, Key: %s", c.Config.Key)
	}
	return nil
}

// 组件接口
// 该接口定义了8个相关的生命周期方法、2个配置相关的方法、1个返回组件类型的方法，共11方法
// 实现了该组接口中的方法的对象类型均可作为组件使用
//
// CreateConfig：返回组件关注的配置参数项、配置参数键、关注配置的子类组件
//
//	该方法在解析配置文件时调用，发生在Awake之后，Start之前。解析配置文件的逻辑在config组件中
//	 该方法有一个参数，表示当前从配置文件中读取到了某个键值，该值在基础类型中的组件没有实际作用，主要用于在衍生类型中
//	 该方法返回3个参数
//	    第一个参数[]object.Command：该组件关注的键要转化的目标值对应的数组
//	    第二个参数[]string，该组件关注的键
//	    第三个参数Componenter，关注该键的衍生类型
//
// InitConfig：初始配置
//
//	发生在CreateConfig之后，Start之前
type Componenter interface {
	Awake() error
	Start(*Cycle) error
	Update(*Cycle) (time.Duration, error)
	OnExit(*Cycle) error
	//TODO 添加以下3个生命周期接口
	// OnReInit：在重新加载配置文件时触发
	// OnReStart：在重新加载可执行文件时触发
	// OnPanic：在发生panic时触发
	// OnReInit(*Cycle) error
	// OnReStart(*Cycle) error
	// OnPanic(*Cycle) error
	GetKind() Kind
}

type BaseComponenter interface {
	Componenter
	CreateConfig(string) []Command
	InitConfig(*Cycle) error
}

// 核心结构体，这个结构存储了整个系统运行过程中需要的东西
// 将伴随整个程序执行的整个周期，贯穿整个系列体
type Cycle struct {
	Prefix     string        // goper的安装路径，一般为可执行文件所在路径
	ConfFile   string        // goper的配置文件全路径（包含文件名）
	BinFile    string        // goper的可执行文件全路径（包含文件名）
	Compts     []Componenter // 所有组件
	KindIndexs map[Kind][]int
}

type SuperComponent struct{}

func (sc *SuperComponent) Awake() error                            { return nil }
func (sc *SuperComponent) Start(pc *Cycle) error                   { return nil }
func (sc *SuperComponent) Update(pc *Cycle) (time.Duration, error) { return -1, nil }
func (sc *SuperComponent) OnExit(pc *Cycle) error                  { return nil }

// func (sc *SuperComponent) OnReInit(pc *Cycle) error                { return nil }
// func (sc *SuperComponent) OnReStart(pc *Cycle) error               { return nil }
// func (sc *SuperComponent) OnPanic(pc *Cycle) error                 { return nil }
func (sc *SuperComponent) GetKind() Kind { return ComptBase }

type BaseComponent struct {
	SuperComponent
}

func (bc *BaseComponent) CreateConfig(key string) []Command { return nil }
func (bc *BaseComponent) InitConfig(pc *Cycle) error        { return nil }

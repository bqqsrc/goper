goper
---

- 作者：晓白齐齐
- 更新：2023.12.9

---
## 组件
### 组件的定义
在goper中，组件是一个实现了约定的方法的一个结构体类型（或者其他类型）。

goper中一切皆组件，整个goper是由一个一个组件组合而成的，goper的核心代码本身也是几个组件。

### Componenter接口
Componenter声明了goper中组件约定的方法：

```
type Componenter interface {
	Awake() error
	Start(*Cycle) error
	Update(*Cycle) (time.Duration, error)
	OnExit(*Cycle) error
	GetKind() Kind
}
```

它声明了5个方法，分别用于实现不同时机组件的逻辑。下面具体介绍这些方法：

- Awake：该方法在组件刚被唤醒时调用，只调用一次，在Start之前调用。该方法不接受任何参数，返回值为error表示是否发生异常，如果没有异常，返回nil。

   该方法的调用与组件顺序无关，各组件是随意调用，因此与组件顺序有关的逻辑不能放在这个方法实现。

   通常该方法用于实现一些组件还没有获取到配置参数时的临时初始化操作。例如日志组件还未获取到配置时，为保证日志得以记录，会赋值一个临时的日志实例，又例如为组件设置默认的配置参数。

- Start：该方法在启动组件时调用，只调用一次，该方法在Awake之后，Update之前调用。 该方法接受一个Cycle指针参数，Cycle是goper中定义的核心结构体之一，这个结构存储了整个goper系统运行过程中需要的东西，后面会具体介绍。该方法返回值为error表示是否发生异常，如果没有异常，返回nil。
   
   通常在该方法做一些组件初始化和启动组件的操作。调用该方法时，组件已经获得了其关注的配置参数值，组件获取配置参数是通过config组件实现的，config是goper的核心组件之一。

	 该方法的调用与modules.go的数组参数有关，会按照该数组索引顺序进行调用。

- Update：该方法在固定的每一轮循环中调用，在Start之后调用，根据返回值决定是否循环调用。该方法接受一个Cycle指针参数。

   该方法返回两个参数，第一个参数为time.Duration，如果其值小于等于0，表示组件的Update不需要循环调用，那么Update只会调用一次，如果其值大于0则会循环调用对应组件的Update。第二个参数为error类型，表示是否发生异常。

	 注意整个系统的Update循环并不是一定存在的，如果所有组件的Update都返回一个小于等于0的Duration，表示所有的组件都不需要Update，循环就会停止。

   通常在该方法实现一些需要循环发生的操作。所有组件启动完成后，每隔一定频率进行一次循环调用该方法。具体间隔多久循环，由core组件获取到的配置决定，默认时间为1秒。

- OnExit：该方法在退出goper系统时调用，做一些必要的保存逻辑。该方法接受一个Cycle指针参数，返回值为error。

   该方法在当前版本还未使用到，属于待完善功能，后续版本将会完善。

- GetKind：该方法返回组件的种类，返回值为组件种类类型。

   组件种类类型定义是uint8的重命名：
	 ```type Kind = uint8```。

	 目前定义了2种种类的组件，ComptBase表示基础组件，ComptHttp表示http组件。

	 ```
	 const (
		ComptBase Kind = iota
		ComptHttp
	)
	 ```
	 
	 用于实现goper基础功能的组件都可以是基础组件，例如核心组件（core）、配置组件（config）、日志组件（log）等都是定义为基础组件。

	 用于实现具体的http逻辑的组件都可以定义为http组件。

所有实现了Componenter声明的接口的类型，都可以作为一个goper组件参与goper系统的运行。事实上，goper中的核心代码也是几个核心组件。

---
## 基础组件
基础组件实现了goper的一些基础的功能，是goper的核心实现。

一般符合以下情况的组件，都可以实现为基础组件：
   - 实现goper的核心功能的组件，例如core、config。
   - 实现通用功能的组件，例如log组件等。
   - 衔接衍生类型组件的组件，例如http组件，作为核心组件和http组件的连接，也是定义为基础组件。

当前版本的goper有4个基础组件：
   - core：该组件是启动goper服务的组件，所有其他组件的Awake、Start、Update等接口都是通过该组件调用。
   - config：该组件用于从配置文件读取配置并分发配置值给关注的组件。
   - log：日志组件。
   - http：该组件从配置组件获取http配置项，并分发给具体的http逻辑组件。
	
### BaseComponenter接口
BaseComponenter接口声明了所有基础组件需要实现的方法。它包含了Componenter接口，因此，一个BaseComponenter接口实现同时也是一个Componenter接口实现。同时它声明了两个新的接口:

```
type BaseComponenter interface {
	Componenter
	CreateConfig(string) []Command
	InitConfig(*Cycle) error
}
```

- CreateConfig：该方法接受一个string参数，表示从配置文件读取到的键，返回值是一个Command数组，表示要进行解析的配置项，Command是组件的参数配置项类型，下面会具体介绍。

- InitConfig：该方法接受一个Cycle指针参数，返回值为error表示是否发生异常。该方法会在读取整个配置文件结束后调用。

所有实现BaseComponenter接口的方法的类型，都可以作为一个基础组件参与goper的运行。

事实上，BaseComponenter会比Componenter更常用。可以看到，BaseComponenter添加的两个接口都是和参数配置项有关，至于具体的配置参数解析行为由config实现。一个可复用的组件，往往需要配置参数来决定具体的逻辑。

---
## 参数配置项
在配置文件中，配置参数都是以键值对的形式存在，goper的config组件从配置文件中读取配置，读取到某个键后，会调用组件列表中组件的相关接口来确定具体是哪个组件关注这个键，进而将这个键对应的值解析到这个组件的具体配置中。

一个Command确定了一个组件的一个参数配置项，它的定义如下：
```
type Command struct {
	Config              ConfigPair
	KeyWord             bool                          
	ConfigFoundCallback func(ConfigPair, *Cycle) error
	ConfigDoneCallback  func(ConfigPair, *Cycle) error
}
```  

它的具体成员如下：
   - 首先，它包含了Config成员，该成员是一个ConfigPair类型，ConfigPair类型定义了一个键值对，如下：

      ```
      type ConfigPair struct {
         Key   string     
         Value ConfigValue 
      }
      ```

      ConfigValue是一个参数值类型的定义，是any的重命名。

   - KeyWord表示该键是否注册为一个关键字，后面会具体介绍到关键字。
   - ConfigFoundCallback是找到该配置项的键时的回调。它接收两个参数，ConfigPair为该配置项中的Config成员，以及一个Cycle指针，返回值为error。
   - ConfigDoneCallback是解析完该配置项时的回调。它的参数同ConfigFoundCallback。

可以看到，一个Command确定了一个配置的键、值，以及找到该配置的回调、解析完该配置的回调，我们把这一组合叫做参数配置项。

基础组件的CreateConfig返回了一组参数配置项的数组，表示该组件关注的键。

特别注意，一个组件可能关注配置文件中的多个键，因此CreateConfig返回的是一个数组，而非单独一个Command。

---
## Cycle核心结构体
Cycle包含了goper整个运行周期都需要的参数，它在不同组件的各个方法中传递，实现不同组件之间的参数传递。具体定义如下：
```
type Cycle struct {
	Prefix     string        
	ConfFile   string        
	BinFile    string        
	Compts     []Componenter
	KindIndexs map[Kind][]int
}
```

它包含以下成员：
   - Prefix：是goper的安装路径，这是在发布版本才会用到的参数，开发模式下即为可执行文件。
   - ConfFile：配置文件的全路径（包含文件名）。
   - BinFile：可执行文件的全路径（包含文件名）。
   - Compts：goper所有组件列表。
   - KindIndexs：不同种类组件的在Compts的索引，它是一个map，键为Kind类型，值为[]int表示对应种类的组件在Compts中的下标。

---
## 实现一个基础组件
一个类型只要实现Componenter声明的所有方法，就实现了一个goper组件。实践中，我们实现一个具体种类的组件可能更有意义。下面介绍如何实现一个HelloWord基础组件：

### 第一步：导入object
```
import "github.com/bqqsrc/goper/object"
```

goper/object定义了goper的基础类型，先导入对应的包。

### 第二步：定义具体的结构体

```
type HelloWord struct {
	object.BaseComponent
	action int 
	name string
}
```

goper/object中定义了BaseComponent（注意后面没有er），这是一个实现了BaseComponenter接口的空组件，它没有任何实际的逻辑操作。一般用于让具体组件包含它以实现继承。

实际上一个组件并不需要用到BaseComponenter的所有方法，有可能只需要用到其中的一两个方法。因此我们在定义自己的基础组件时，可以通过在结构体中包含BaseComponent来继承BaseComponent方法，这样就不用额外去实现其他默认方法。

当然你也可以通过实现BaseComponenter中所有的方法的方式来实现。但是我建议还是使用包含BaseComponent的方式，一方面会使得代码更加简洁，另一方面在后面的迭代更新过程中，有可能会新增方法或删除已有的方法，通过后者的方式能够降低迭代更新带来的修改代码的概率，提高兼容性。

在本例中HelloWord定义了一个int参数action，一个string参数name。

### 第三步：实现参数配置项
```
func (h *HelloWord) CreateConfig(key string) []object.Command {
	return []object.Command {
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
```

在该例子中，HelloWord组件关注两个参数，其中一个键为action，需要将其解析为int类型，另一个键为name，解析为string类型。action键不需要任何回调，而name键则添加了发现name时的回调和解析完成时的回调。

配置文件解析过程中，通过组件的CreateConfig返回的参数配置项列表，可以得知关注action和name的组件为HelloWord。如果发现了action键，会将配置文件action键对应的值解析到h.action中，发现了name键，会将配置文件action键对应的值解析到h.name中，同时发现name键时会发生回调，解析完name后会发生回调。

特别注意的是，Command中的第一个成员object.ConfigPair的Value值必须是一个指针，这样才能将结果解析到目标变量中。

因为我们知道name键会解析为一个string，并且解析到h.name，因此name键对应的最后一个回调，可以改为如下：
```
func(conf object.ConfigPair, pc *object.Cycle) error {
	fmt.Printf("parse finish, key: %s, name: %s ", conf.Key, *conf.Value.(*string))
	return nil
}
```

因为conf.Value是一个指针，指向h.name，因此实际上和前者的回调是等价的。

### 第四步：实现具体方法逻辑
并不是所有BaseComponenter声明的方法都会用到，这里我们只用到Start方法：
```
func (h *HelloWord) Start(pc *object.Cycle) error {
	if h.action > 0 {
		fmt.Printf("HelloWord, my name is %s, I got a num larger than 0", h.name)
	} else if h.action < 0{
		fmt.Printf("HelloWord, my name is %s, I got a num less than 0", h.name)
	} else {
		fmt.Printf("HelloWord, my name is %s, I got a num equal 0", h.name)
	}
	return nil
}
```

获取到了配置参数后，我们可以在Start接口中进行组件的初始化、启动组件等操作。

因为HelloWord已经包含了BaseComponent，所以其他不需要用到的方法我们不用再实现一次。如果没有包含BaseComponent，那还需要再实现其他的几个方法。

由于BaseComponent已经实现了GetKind接口，并且返回的值为object.ComptBase，因此，我们也不需要再实现一次GetKind。

到此，我们就实现了一个新的基础类型，该类型关注action和name两个参数，获取到参数后输出日志。

### 完整的代码如下
```
import (
	"fmt"
	"github.com/bqqsrc/goper/object"
)

type HelloWord struct {
	object.BaseComponent
	action int 
	name string
}

func (h *HelloWord) CreateConfig(key string) []object.Command {
	return []object.Command {
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
	} else if h.action < 0{
		fmt.Printf("HelloWord, my name is %s, I got a num less than 0", h.name)
	} else {
		fmt.Printf("HelloWord, my name is %s, I got a num equal 0", h.name)
	}
	return nil
}
```

---
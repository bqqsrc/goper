# goper

Goper 是一个用 go 语言开发的模块化、配置化的 Web 框架，有更好的扩展性、复用性、可配置化、高度可拔插特性。

## 开始使用
### 导入Goper
Goper 是一个模块化的 web 框架，使用 Goper ，需要先导入以下几个基本模块包：
```
import (
	"github.com/bqqsrc/goper"
	"github.com/bqqsrc/goper/config"
	"github.com/bqqsrc/goper/core"
	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/http/hcore"
	"github.com/bqqsrc/goper/object"
)
```

### 简单的Demo
我们运行一个简单的例子，新建一个 example.go ，添加如下代码
```
package main

import (
	"fmt"

	"github.com/bqqsrc/goper"
	"github.com/bqqsrc/goper/config"
	"github.com/bqqsrc/goper/core"
	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/http/hcore"
	"github.com/bqqsrc/goper/object"
)

type Example struct {
	http.HttpComponent
}

func (e *Example) Handler(c *http.Context) http.HttpPhase {
	c.Response.SetData("key", "HelloWorld")
	return http.HttpNext
}

func (e *Example) CreateHandler(dataConfig object.ConfigValue) (http.HttpHandler, http.HttpPhase, error) {
	return e, http.HttpLogic, nil
}

var compts = []object.Componenter{
	&core.Core{},
	&config.Config{},
	&http.Http{},
	&hcore.HCore{},
	&Example{},
}

func main() {
	if err := goper.Launch(compts); err != nil {
		fmt.Printf("goper.Launch err: %v", err)
	}
}
```

### 配置Http
在 example.go 的同级目录创建文件 goper.gs ，写入：
```
http: {}
```

### 运行
然后运行 example.go ，并在浏览器中查看： localhost:80
```
go run example.go
```

### 更多例子
如果需要了解更多例子，请前往 Goper Quick Start 进行查看。

---
## 贡献者和问题咨询
你可以通过以下方式联系作者，添加作者时请注明来意（谢绝广告、推销）。由于作者工作繁忙，一般晚上才查看一次微信和QQ，如没有马上通过，还请见谅。

作者的联系方式：

   - 微信：_xbqq_

   - QQ：1662488230

如果你出于以下目的想联系作者，作者非常欢迎，也非常期待：
1. 有意愿成为项目的贡献者
2. 想为项目提供一些思路建议、优化意见
3. 找到一些bug
4. 尝试使用这个项目，在使用过程中有任何不清楚的问题想咨询
5. 其他技术交流

Goper 为作者业余所作，当前只有作者一个人在开发维护，当然会应用到一些开源项目。

由于作者的技术能力和项目经验有限， Goper 当前版本还存在很多问题亟需改善，这里有些问题作者已经了解，有些问题作者还未了解。后期将会不断进行迭代更新，欢迎你持续关注。

---
## 结语
作者近期正在寻找新的工作机会，如果你觉得这个项目对你有用，方便点个 star 的话，还请不吝赐星，这对作者找工作会有一定的帮助。当然如果这不符合你的习惯也没关系的，因为你对项目的关注和支持也是对作者很大的支持了。

Goper 通过模块化（在Goper中，也称之为组件），将各个功能开发为彼此独立的组件，各个组件搭建成一个 Goper 整体，各组件各司其职，互相协作，将 Goper 服务搭建起来。正是这种模块化的设计，使得 Goper 真正实现了可拔插，在 Goper 中添加、删除组件是一件很容易的事，甚至替换组件也是轻而易举的事。

例如，在 Goper 中，启动其他所有组件的实现以及读取配置的实现也是两个组件，分别是 core 组件和 config 组件。如果你觉得当前版本的 config 组件效率不够高，或者想要替换成其他的配置文件类型，例如json文件，只要你按照当前读取配置的规则开发出一个新的读取配置的组件，把原版的 config 组件一替换，就可以马上无缝用上你自己的组件了。

Goper 的模块化设计，也使得 Goper 具有更高的扩展性，只要按照约定的接口实现方法的类型，就可以作为一个新的组件添加到 Goper 的运行中，这使得开发自己的组件非常容易，也可以将自己的组件作为一个单独的项目发布给其他人使用。

Goper 的配置化特性也是 Goper 复用性的保障，通过不同的配置实现不同的逻辑，使得一份代码多处使用。

Goper 当前版本还存在很多问题，不过由于 Goper 以上的特性，后期维护过程中，开发更高效率的组件、新功能组件将会更容易，后期将会有更多的功能推出。

限于作者的技术能力和项目经验，当前框架还有很多问题，当前的设计思路还有很多不足之处，如果你有任何建议意见，欢迎你联系作者。谢谢。

最后，感谢你的阅读、关注和支持。

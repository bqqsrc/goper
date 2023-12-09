goper
---

- 作者：晓白齐齐
- 更新：2023.12.9


---
## 新增一种组件种类
当前版本的goper定义了两种组件种类：基本组件和http组件。在基本组件的基础上衍生一些其他功能的组件种类，我们把它称为衍生种类。一般来说新增的组件种类都是衍生组件，例如我们可能需要一种email种类的组件来实现一套email框架等。

我们需要定义一个新的基本组件，与核心组件（core和config等）交互，从核心组件获取配置项等，然后分发给具体的衍生组件。下面我们以新增一种为ComptDemo的种类的组件为例，给出新增组件类型的步骤：

### 新增组件种类常数定义
如果你是直接修改goper的源码，那么可以直接在goper/object/constant.go中新增你的组件种类，如下
```
const (
	ComptBase Kind = iota
	ComptHttp
	ComptDemo  // 你新增的ComptDemo
)
const KindCount = 3
```

如果你没有直接在goper源码修改，那么可以通过如下方法来定义种类值：
```
import (
	"github.com/bqqsrc/goper/object"
)
```

当前版本的goper在定义新的衍生种类时，只会关注于将配置文件中对应的配置值分发给对应的衍生种类，至于其他的流程


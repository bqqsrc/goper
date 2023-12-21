# goper

Goper是一个用go语言开发的模块化Web框架，有更好的扩展性、复用性、高度可拔插特性。

## 开始使用
### 导入Goper
Goper是一个模块化的web框架，使用Goper，需要导入以下几个基本模块包：
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

#### 运行Goper
我们运行一个简单的例子，新建一个example.go，添加如下代码
```
package main 
type Example struct {
	http.HttpComponent 
}

func (e *Example) Handler(c *http.Context) http.HttpPhase {
	c.Response = ResponseData(1)
	c.Response.SetData("key1", "value1")
	c.Response.SetData("key2", c.ParamsData)
	return http.HttpNext
}


```


## 配置

## 组件

## 定制自己的组件

## 模块化和扩展性

## 结语

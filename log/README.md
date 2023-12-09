goper
---

- 作者：晓白齐齐
- 更新：2023.12.9

---
## 使用日志模块
### 第一步：导入
```
import "github.com/bqqsrc/goper/log"
```

### 第二步：调用接口
log模块支持以下接口：

Debugf(format string, v ...any)

Debug(v ...any)

Debugln(v ...any)

Infof(format string, v ...any)

Info(v ...any)

Infoln(v ...any)

Warnf(format string, v ...any)

Warn(v ...any) 

Warnln(v ...any) 

Errorf(format string, v ...any)

Error(v ...any) 

Errorln(v ...any) 

Criticalf(format string, v ...any)

Critical(v ...any)

Criticalln(v ...any)

Fatalf(format string, v ...any) 

Fatal(v ...any)

Fatalln(v ...any)

LogConsolef(format string, v ...any)

LogConsolefln(format string, v ...any) 

LogConsole(v ...any)

LogConsoleln(v ...any)

其中，Debug、Info、Warn、Error、Critical对应5个等级的日志输出接口，Fatal输出错误后会调用exit退出进程，LogConsole不管设置日志等级是多少，都输出日志到终端。

### 第三步：配置文件添加配置项
日志模块关注log配置项，配置项中的各个字段对应的意义如下：

- level:
   - 意义：默认输出的日志等级，配置为一个数字
   - 默认值: 3
   - 取值和说明如下：
      - 0：不设置等级，所有日志都输出 
      - 1：LevelDebug 输出Debug及以上等级的日志（Debug、Info、Warn、Error、Critical）
      - 2：LevelInfo 输出Info及以上等级的日志（Info、Warn、Error、Critical）
      - 3：LevelWarn 输出Warn及以上等级的日志（Warn、Error、Critical）
      - 4：LevelError 输出Error及以上等级的日志（Error、Critical）
      - 5：LevelCritical 输出Critical及以上等级的日志（Critical）
      - 大于5：所有日志都不输出

- flag：
   - 意义：日志中要输出的内容对应的相关的flag标志，配置为一个字符串，如果要输出多个内容，就用|将多个标签进行或运算
   - 默认为：LStdFlags|LLevel
   - 取值和说明如下：
      - LDate：输出当地的时区的日期 
      - LTime：输出当地的时区的时间
      - LMicroseconds：输出当地的时区的时间时，精确到微秒，这个设置了默认LTime也设置
      - LNanosceonds：输出当前时间戳的纳秒数，注意这里是输出总纳秒数，可用于在一段代码前后输出查看对应代码运行的时间耗时
      - LLongFile：输出调用日志接口所在长文件名称（即包含全路径和文件名）
      - LShortFile：输出调用日志接口所在短文件名称（即仅仅是文件名），如果LLongFile和LShortFile同时设置了，将只输出长文件名称
      - LUTC：输出时间时，将时间转换为UTC时间，这个要配合LDate或LTime或LMicroseconds使用才有效
      - LTag：输出tag字段设置的字符串
      - LPreTag：将tag字段设置的字符串输出到最前面
      - LLevel：输出日志的等级
      - LStdFlags：等于LDate|LTime
      - 注意：|前后不要留空格，由于gson的解析规则，gson对空格敏感，如果出现空格，后面的值会被gson舍弃，或者用`将字符串包起来

- tag：
   - 意义：日志输出附带的默认标签，配置为一个字符串
   - 默认值: 空字符串
   - 取值和说明：只有在flag设置了LTag或者LPreTag，tag才会在日志中输出

- logfile：
   - 意义：输出的日志文件，配置为一个字符串
	 - 默认值：空字符串
   - 取值和说明：
      - 如果填写的是相对路径，则会以当前运行路径为路径前缀，追加配置的路径，如果日志文件不存在，则新建日志文件，如果所填路径存在但是不是一个文件，则会报错
       如果不设置，goper将默认不输出到日志文件
      - 注意：路径配置用一对“\`”括起来，因为“\:”是gson的保留字符，如果没有用“\`”包起来，则会出现遗漏“\:”的情况

该配置为一个字典，用{}包起来，配置样例：
```
log: {
	level: 2  # 输出Warn以上级别的日志
	tag: [goper]  # 每个日志附带一个[goper]标签
	flag: LStdFlags|LLevel
	logfile: `./tmp/log/goper.log`
}
```

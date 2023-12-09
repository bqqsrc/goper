//  Copyright (C) 晓白齐齐,版权所有.

package log

import (
	// "github.com/bqqsrc/bqqg/errors"
	// "time"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bqqsrc/bqqg/file"
	logg "github.com/bqqsrc/bqqg/log"
	"github.com/bqqsrc/goper/object"
)

type Log struct {
	object.BaseComponent
	gpLoger *logg.Log
	logConf logConfig
}

type logConfig struct {
	Level   logg.LogLevelType `gson:"level"`
	Tag     string            `gson:"tag"`
	Flag    string            `gson:"flag"`
	Logfile string            `gson:"logfile"`
}

func (l *Log) Awake() error {
	Debugln("Log Awake()")
	l.logConf = logConfig{3, "", "", ""}
	l.gpLoger = logg.Default()
	myLog = l
	return nil
}

func (l *Log) CreateConfig(key string) []object.Command {
	Debugf("ConfigValue CreateConfig(%s), config key: log", key)
	return []object.Command{
		{
			object.ConfigPair{
				"log",
				&l.logConf,
			},
			false,
			nil,
			l.initLoger,
		},
	}
}

func (l *Log) initLoger(keyValue object.ConfigPair, pc *object.Cycle) error {
	Debugf("initLoger, config: %v, %v", l.logConf, *keyValue.Value.(*logConfig))
	flag := logg.LStdFlags | logg.LLevel
	if l.logConf.Flag != "" {
		flag = 0
		flagStr := strings.TrimSpace(l.logConf.Flag)
		flagStrs := strings.Split(flagStr, "|")
		for _, v := range flagStrs {
			switch v {
			case "LDate":
				flag = flag | logg.LDate
			case "LTime":
				flag = flag | logg.LTime
			case "LMicroseconds":
				flag = flag | logg.LMicroseconds
			case "LNanosceonds":
				flag = flag | logg.LNanosceonds
			case "LLongFile":
				flag = flag | logg.LLongFile
			case "LShortFile":
				flag = flag | logg.LShortFile
			case "LUTC":
				flag = flag | logg.LUTC
			case "LTag":
				flag = flag | logg.LTag
			case "LPreTag":
				flag = flag | logg.LPreTag
			case "LLevel":
				flag = flag | logg.LLevel
			case "LStdFlags":
				flag = flag | logg.LStdFlags
			}
		}
	}

	// Debugln("initLoger 2")
	if l.logConf.Logfile != "" {
		if file.IsDir(l.logConf.Logfile) {
			return fmt.Errorf("%s is exist but not a file", l.logConf.Logfile)
		}
		dr := path.Dir(l.logConf.Logfile)
		if !file.Exist(dr) {
			os.MkdirAll(dr, os.ModeDir)
		}

		f, err := os.OpenFile(l.logConf.Logfile, os.O_APPEND|os.O_CREATE, os.ModePerm)
		//TODO 如果目录不存在，创建目录
		if err != nil {
			return fmt.Errorf("open log file %s error, err: %s", l.logConf.Logfile, err)
		} else {

			// Debugln("initLoger 4", l.logConf.Logfile, l.logConf.Level)
			l.gpLoger = logg.New(l.logConf.Tag, l.logConf.Level, flag, f)
		}
	} else {

		// Debugln("initLoger 5")
		l.gpLoger = logg.New(l.logConf.Tag, l.logConf.Level, flag, os.Stderr)
	}
	// Warnln("initLoger 6")
	Infof("initLoger success, logfile: %s, logLevel: %d, logTag: %s, logFlag: %s", l.logConf.Logfile, l.logConf.Level, l.logConf.Tag, l.logConf.Flag)
	l.gpLoger.SetCallDepth(1)
	// cycle.Loger = gpLoger
	return nil
}

var myLog *Log

func Debugf(format string, v ...any) {
	if myLog == nil {
		LogConsolefln(format, v...)
	} else {
		myLog.gpLoger.Debugf(format, v...)
	}
}

func Debug(v ...any) {
	if myLog == nil {
		LogConsole(v...)
	} else {
		myLog.gpLoger.Debug(v...)
	}
}

func Debugln(v ...any) {
	if myLog == nil {
		LogConsoleln(v...)
	} else {
		myLog.gpLoger.Debugln(v...)
	}
}

func Infof(format string, v ...any) {
	if myLog == nil {
		LogConsolefln(format, v...)
	} else {
		myLog.gpLoger.Infof(format, v...)
	}
}

func Info(v ...any) {
	if myLog == nil {
		LogConsole(v...)
	} else {
		myLog.gpLoger.Info(v...)
	}
}

func Infoln(v ...any) {
	if myLog == nil {
		LogConsoleln(v...)
	} else {
		myLog.gpLoger.Infoln(v...)
	}
}

func Warnf(format string, v ...any) {
	if myLog == nil {
		LogConsolefln(format, v...)
	} else {
		myLog.gpLoger.Warnf(format, v...)
	}
}
func Warn(v ...any) {
	if myLog == nil {
		LogConsole(v...)
	} else {
		myLog.gpLoger.Warn(v...)
	}
}
func Warnln(v ...any) {
	if myLog == nil {
		LogConsoleln(v...)
	} else {
		myLog.gpLoger.Warnln(v...)
	}
}

func Errorf(format string, v ...any) {
	if myLog == nil {
		LogConsolefln(format, v...)
	} else {
		myLog.gpLoger.Errorf(format, v...)
	}
}
func Error(v ...any) {
	if myLog == nil {
		LogConsole(v...)
	} else {
		myLog.gpLoger.Error(v...)
	}
}
func Errorln(v ...any) {
	if myLog == nil {
		LogConsoleln(v...)
	} else {
		myLog.gpLoger.Errorln(v...)
	}
}

func Criticalf(format string, v ...any) {
	if myLog == nil {
		LogConsolefln(format, v...)
	} else {
		myLog.gpLoger.Criticalf(format, v...)
	}
}
func Critical(v ...any) {
	if myLog == nil {
		LogConsole(v...)
	} else {
		myLog.gpLoger.Critical(v...)
	}
}
func Criticalln(v ...any) {
	if myLog == nil {
		LogConsoleln(v...)
	} else {
		myLog.gpLoger.Criticalln(v...)
	}
}

func Fatalf(format string, v ...any) {
	if myLog == nil {
		LogConsolefln(format, v...)
		os.Exit(1)
	} else {
		myLog.gpLoger.Fatalf(format, v...)
	}
}
func Fatal(v ...any) {
	if myLog == nil {
		LogConsole(v...)
		os.Exit(1)
	} else {
		myLog.gpLoger.Fatal(v...)
	}
}
func Fatalln(v ...any) {
	if myLog == nil {
		LogConsoleln(v...)
		os.Exit(1)
	} else {
		myLog.gpLoger.Fatalln(v...)
	}
}
func LogConsolef(format string, v ...any) {
	fmt.Printf(format, v...)
}
func LogConsolefln(format string, v ...any) {
	fmt.Printf(format+"\n", v...)
}
func LogConsole(v ...any) {
	fmt.Print(v...)
}
func LogConsoleln(v ...any) {
	fmt.Println(v...)
}

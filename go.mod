module github.com/bqqsrc/goper

go 1.18

require (
	github.com/bqqsrc/jsoner v0.0.1
	github.com/bqqsrc/loger v0.0.0
	gopkg.in/ini.v1 v1.66.6
)

replace (
	//github.com/bqqsrc/jsoner v0.0.1 => ../jsoner
	github.com/bqqsrc/loger v0.0.0 => ../loger
)

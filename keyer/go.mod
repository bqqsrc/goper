module github.com/bqqsrc/goper/keyer

go 1.18

require (
	github.com/bqqsrc/goper v0.0.0
	github.com/bqqsrc/imaper v0.0.0
	github.com/bqqsrc/jsoner v0.0.1
	github.com/bqqsrc/loger v0.0.0
)

require gopkg.in/ini.v1 v1.66.6 // indirect

replace github.com/bqqsrc/goper v0.0.0 => ../

replace (
	github.com/bqqsrc/imaper v0.0.0 => ../../imaper
	//github.com/bqqsrc/jsoner v0.0.1 => ../../jsoner
	github.com/bqqsrc/loger v0.0.0 => ../../loger
)

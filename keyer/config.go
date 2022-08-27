package keyer

import (
	"github.com/bqqsrc/goper"
	"github.com/bqqsrc/jsoner"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	//	"github.com/bqqsrc/loger"
)

type Key struct {
	Method      string      `json:"method"`
	Value       interface{} `json:"value"`
	SuccessCode int         `json:"success_code"`
	FailedCode  int         `json:failed_code"`
}

func (k *Key) check() error {
	switch k.Method {
	case Method_None, Method_Code:
		break
	case Method_Value, Method_Function:
		if k.Value == nil {
			return fmt.Errorf("value can not be nil while method is %s", k.Method)
		}
		break
	default:
		return fmt.Errorf("unsupport method type: %s, only can't be %s, %s, %s or %s", k.Method, Method_None, Method_Value, Method_Code, Method_Function)
	}
	return nil
}

func (k *Key) equal(key Key) bool {
	return k.Method == key.Method && k.Value == key.Value && k.SuccessCode == key.SuccessCode && k.FailedCode == key.FailedCode
}

type RouterKey struct {
	PrimaryKey  string                    `json:"primary"`
	PostDataKey string                    `json:"post_key"`
	ActionMap   map[string]map[string]Key `json:"logic"`
}

var EmptyRouterKey = RouterKey{}

func (r *RouterKey) equal(routerKey RouterKey) bool {
	if r.PrimaryKey != routerKey.PrimaryKey {
		return false
	}
	if r.PostDataKey != routerKey.PostDataKey {
		return false
	}
	if len(r.ActionMap) != len(routerKey.ActionMap) {
		return false
	}
	for acKey, acValue := range r.ActionMap {
		if tmpAcValue, ok := routerKey.ActionMap[acKey]; !ok {
			return false
		} else {
			if len(acValue) != len(tmpAcValue) {
				return false
			} else {
				for key, value := range acValue {
					if tmpValue, ok := tmpAcValue[key]; !ok {
						return false
					} else {
						if !value.equal(tmpValue) {
							return false
						}
					}
				}
			}
		}
	}
	return true
}

func (r *RouterKey) check() error {
	var build strings.Builder
	if r.ActionMap != nil {
		for ac, keyMap := range r.ActionMap {
			for key, valueMap := range keyMap {
				if err := valueMap.check(); err != nil {
					build.WriteString(fmt.Sprintf("%s->%s error: ", ac, key))
					build.WriteString(err.Error())
					build.WriteString("\n")
				}
			}
		}
	}
	if build.Len() > 0 {
		var tmpBuild strings.Builder
		tmpBuild.WriteString("some router key settings error: \n")
		tmpBuild.WriteString(build.String())
		return errors.New(tmpBuild.String())
	}
	return nil
}

type RouterKeys struct {
	Map map[string]RouterKey
}

func (r *RouterKeys) Bytes2Config(filepath string, b []byte) error {
	if r.Map == nil {
		r.Map = make(map[string]RouterKey)
	}
	funcName := "RouterKeys.Bytes2Config"
	var routerKeyMap map[string]RouterKey
	err := json.Unmarshal(b, &routerKeyMap)
	//loger.Debugf("routerKeyMap is %v", routerKeyMap)
	if err != nil {
		return goper.GoperError(funcName, "json.Unmarshal error, file is: %s, err is: %s", filepath, err)
	}
	var build strings.Builder
	for key, value := range routerKeyMap {
		//loger.Debugf("key is %s, value is %v", key, value)
		if err = value.check(); err != nil {
			build.WriteString(fmt.Sprintf("url is %s, error is: \n%s\n", key, err))
		} else {
			if _, ok := r.Map[key]; ok {
				build.WriteString(fmt.Sprintf("redeclare router, url is %s\n", key))
			} else {
				r.Map[key] = value
			}
		}
	}
	if build.Len() > 0 {
		var tmpBuild strings.Builder
		tmpBuild.WriteString("some routerKeys error:\n")
		tmpBuild.WriteString(build.String())
		tmpBuild.WriteString("\n")
		return goper.GoperError(funcName, tmpBuild.String())
	}
	return nil
}

var globalRouterKeys RouterKeys

func InitKeyConfig() error {
	config, err := goper.GetConfigArr("routerKey", "configDir", ",")
	if err != nil {
		return err
	}
	if err = jsoner.ReadAllConfig(config, &globalRouterKeys); err != nil {
		return err
	}
	//loger.Debugf("InitKeyConfig %v", globalRouterKeys)
	return nil
}

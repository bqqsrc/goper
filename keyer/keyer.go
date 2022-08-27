package keyer

import (
	"encoding/json"
	"net/http"

	Goper "github.com/bqqsrc/goper"

	"github.com/bqqsrc/imaper"

	"github.com/bqqsrc/loger"
)

func KeyerRequest(url string, urlParams, bodyData map[string]interface{}, cookies []*http.Cookie) ([]byte, []*http.Cookie, error) {
	funcName := "KeyerRequest"
	if bodyData == nil || len(bodyData) == 0 {
		return nil, nil, Goper.GeneralError(100, funcName, errorMsg, url)
	}
	defaultRouterKeys, ok := globalRouterKeys.Map[Default]
	if !ok {
		loger.Infof("%s, Warning: %s not found in globalRouterKeys", funcName, Default)
	}
	urlRouterKeys, ok := globalRouterKeys.Map[url]
	if !ok {
		loger.Warnf("%s, Warning: %s not found in globalRouterKeys", funcName, url)
	}
	if defaultRouterKeys.equal(EmptyRouterKey) && urlRouterKeys.equal(EmptyRouterKey) {
		return nil, nil, Goper.GeneralError(200, funcName, errorMsg, notFoundConfigTip, url, Default, checkConfigTip)
	}
	if responseData, cookies, err := getResponseByConfig(url, urlRouterKeys, defaultRouterKeys, bodyData, cookies); err != nil {
		return nil, nil, err
	} else {
		if responseResult, err := json.Marshal(responseData); err != nil {
			return nil, nil, Goper.GoperError(funcName, "%s, json.Marshal, Error: %s", funcName, err)
		} else {
			return responseResult, cookies, nil
		}
	}
}

func getResponseByConfig(url string, urlRouterKeys, defaultRouterKeys RouterKey, postData map[string]interface{}, cookies []*http.Cookie) (interface{}, []*http.Cookie, error) {
	funcName := "getResponseByConfig"
	loger.Tracef("%s, begin, url: %s", funcName, url)
	//获取关键字key
	primaryKey := urlRouterKeys.PrimaryKey
	if primaryKey == "" {
		primaryKey = defaultRouterKeys.PrimaryKey
	}
	if primaryKey == "" {
		return nil, nil, Goper.GeneralError(300, funcName, errorMsg, getConfigErrorTip, url, PrimaryKey, Default, PrimaryKey, "key not found", checkConfigTip)
	}
	//获取关键字值
	primaryValue, err := getStringFromPostData(url, primaryKey, funcName, postData)
	if err != nil {
		return nil, nil, err
	}
	//获取所有逻辑配置
	logicConfig := urlRouterKeys.ActionMap
	if logicConfig == nil {
		logicConfig = defaultRouterKeys.ActionMap
	}
	if logicConfig == nil {
		return nil, nil, Goper.GeneralError(300, funcName, errorMsg, getConfigErrorTip, url, LogicConfig, Default, LogicConfig, "key not found", checkConfigTip)
	}
	//获取对应关键字值的路由的逻辑配置
	routerConfig, ok := logicConfig[primaryValue]
	if !ok {
		return nil, nil, Goper.GeneralError(301, funcName, errorMsg, getConfigErrorTip, url, LogicConfig, primaryValue, "key not found", checkConfigTip)
	}
	//获取postParmas的key
	postParamsKey := urlRouterKeys.PostDataKey
	if postParamsKey == "" {
		postParamsKey = defaultRouterKeys.PostDataKey
	}
	if postParamsKey == "" {
		return nil, nil, Goper.GeneralError(300, funcName, errorMsg, getConfigErrorTip, url, PostParamsKey, Default, PostParamsKey, "key not found", checkConfigTip)
	}
	//获取postPrarms
	var postParams map[string]interface{}
	if postParams, err = getMapFromPostData(url, postParamsKey, funcName, postData); err != nil {
		loger.Warnf("%s, getpost data error, err is %s, key is %s", funcName, postParamsKey, err)
		//return nil, nil, err
	}
	return getResponseData(url, primaryValue, routerConfig, postData, postParams, cookies)
}

func getStringFromPostData(url, key, funcName string, postData map[string]interface{}) (string, error) {
	if result, err := imaper.GetStringFromMaps(key, postData); err != nil {
		return "", Goper.GeneralError(400, funcName, errorMsg, getPostDataErrorTip, url, key, err, checkPostDataTip)
	} else {
		return result, nil
	}
}

func getMapFromPostData(url, key, funcName string, postData map[string]interface{}) (map[string]interface{}, error) {
	if result, err := imaper.GetMapFromMaps(key, postData); err != nil {
		return nil, Goper.GeneralError(400, funcName, errorMsg, getPostDataErrorTip, url, key, err, checkPostDataTip)
	} else {
		return result, nil
	}
}

func getResponseData(url, primary string, keyConfig map[string]Key, postData, postParams map[string]interface{}, cookies []*http.Cookie) (interface{}, []*http.Cookie, error) {
	funcName := "getResponseData"
	loger.Tracef("%s, begin, url: %s, primary: %s", funcName, url, primary)
	responseData := make(map[string]interface{})
	result := true
	codeKey := CodeKey
	responseSuccessCode := DefaultSuccessCode
	responseFailedCode := DefaultFailedCode
	newCookies := make([]*http.Cookie, 0)
	for key, value := range keyConfig {
		response, keyResponseData, err, isCodeKey, successCode, failedCode, cookie := getOneKeyResponseData(url, primary, key, value, postData, postParams, cookies)
		if response {
			return keyResponseData, cookie, err
		} else {
			if err != nil {
				loger.Errorf("%s, getOneKeyResponseData error, url is %s, key is %s, err is %s", funcName, url, key, err)
				result = false
			}
		}
		if cookie != nil {
			newCookies = append(newCookies, cookie...)
		}
		if isCodeKey {
			codeKey = key
			responseSuccessCode = successCode
			responseFailedCode = failedCode
		}
		if key != None && !isCodeKey {
			responseData[key] = keyResponseData
		}
	}
	if result {
		responseData[codeKey] = responseSuccessCode
	} else {
		responseData[codeKey] = responseFailedCode
	}
	return responseData, newCookies, nil // Goper.GoperError(funcName, "getResponseData, check log")
}

//返回值1：是否直接将返回数据作为响应发送回客户端，如果返回true，将不会继续遍历其余键值，而把第2个返回值作为响应发送到客户端
//返回值2：该键的值，如果第1个返回值为true，则该返回值作为本次请求的响应
//返回值3：本次遍历是否成功获取了数值
//返回值4：是否为返回码，如果该值返回true，将再所有键遍历完成，如果存在没有获取成功的键，这个key将赋值为返回值6的失败码，如果所有键都获取成功，这个键将赋值为返回值5的成功码
//返回值5：成功码，只有在返回值4为true时，这个返回值才有效
//返回值6：失败码，只有在返回值4为true时，这个返回值才有效
//func getOneKeyResponseData(url, primary, key string, keyConfig, postData, postParams map[string]interface{}) (bool, interface{}, bool, bool, int, int) {
func getOneKeyResponseData(url, primary, key string, keyConfig Key, postData, postParams map[string]interface{}, cookies []*http.Cookie) (bool, interface{}, error, bool, int, int, []*http.Cookie) { // bool, interface{}, bool, bool, int, int) {
	funcName := "getOneKeyResponseData"
	loger.Tracef("%s, begin, url: %s, primary: %s, key: %s", funcName, url, primary, key)
	switch keyConfig.Method {
	case Method_None: //空类型时直接返回nil，且结果为true
		return false, nil, nil, false, 0, 0, nil
	case Method_Code:
		successCode := DefaultSuccessCode
		if keyConfig.SuccessCode != 0 {
			successCode = keyConfig.SuccessCode
		}
		failedCode := DefaultFailedCode
		if keyConfig.FailedCode != 0 {
			failedCode = keyConfig.FailedCode
		}
		return false, nil, nil, true, successCode, failedCode, nil
	case Method_Value:
		value := keyConfig.Value // keyConfig[Value]
		if value != nil {
			if valueStr, ok := value.(string); ok {
				if valueStr == PostData {
					return false, postData, nil, false, 0, 0, nil
				} else if valueStr == PostParams {
					return false, postParams, nil, false, 0, 0, nil
				}
			}
		}
		return false, value, nil, false, 0, 0, nil
	case Method_Function:
		method := keyConfig.Value.(string)
		rsFunc, ok := responseFunc[method]
		if !ok {
			return true, nil, Goper.GeneralError(102, funcName, errorMsg, url, LogicConfig, primary, key, Value, "RegisterResponseFunc"), false, 0, 0, nil
		}
		response, responseData, cookies, err := rsFunc(url, primary, key, postData, postParams, cookies)
		return response, responseData, err, false, 0, 0, cookies
	default:
		return true, nil, Goper.GeneralError(101, funcName, errorMsg, url, LogicConfig, primary, key, Method,
			keyConfig.Method, Method_None, Method_Value, Method_Code, Method_Function, checkConfigTip), false, 0, 0, nil
	}
}

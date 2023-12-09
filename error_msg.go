package goper

var errorMsg = map[int]string{
	1000: "read body error, ioutil.ReadAll error, err is %s",
	1001: "parse body error, json.Unmarshal error, err is %s",
	1002: "parse rawQuery error, url.QueryUnescape error, err is %s",
}

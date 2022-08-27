package goper

import (
	"net/http"
)

func GoResponse(w http.ResponseWriter, data []byte) {
	w = WriteGoResponseHeader(w)
	//TODO统计流量改为连同包头、其他更底层的数据也统计进来
	//len(data)可以用于流量统计
	//log.Println("lenght is ", len(data))
	w.Write(data)
}

func WriteGoResponseHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Add("Server", serverName)
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	return w
}

func CheckRequest(r *http.Request) bool {
	//TODO 这个函数用于过滤一些不必要的恶意的请求
	//防刷，防重复提交（表单一节的处理）
	//从包头获取User-Agent字段的值，判断客户端，确定是否接收请求
	//r.ContentLength记录了请求body的长度，可以用于流量统计
	//log.Println("request length ", r.ContentLength, len(test))
	return true
}

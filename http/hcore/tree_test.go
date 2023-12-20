// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

package hcore

import (
	"fmt"
	// "strings"
	"testing"

	"github.com/bqqsrc/bqqg/test"
	"github.com/bqqsrc/goper/http"
)

type routeVal struct {
	val    int
	route  string
	domain string
	method string
}

type resultInfo struct {
	index int
	pre   bool
	tsr   bool
}

func TestStaticRouter(t *testing.T) {
	n := &node{}
	var index int
	var routeStr string
	routeValMap := []routeVal{
		// 普通路由测试样例
		{1, "/router1/level1/test1", "", ""}, //0 完全匹配：/router1/level1/test1, 前缀：/router1/level1/test123, /router1/level1/test1/val1, 去掉/重定向：/router1/level1/test1/
		{2, "/test/level/val1/", "", ""},     //1 添加/重定向：/test/level/val1, 完全匹配：/test/level/val1/,  前缀：/test/level/val1/tee, /test/level/val1/tee2

		// 前缀和重定向取舍，取重定向
		{3, "/pre1/pre2/pre3/", "", ""}, //2 /pre1/pre2/pre3 可以匹配到前缀/pre1/pre2/，也可以匹配到重定向/pre1/pre2/pre3/，取重定向
		{4, "/pre1/pre2/", "", ""},      //3

		// 长短不一的前缀取路由更长的
		{5, "/for1/for2/for3/", "", ""}, //4 /for1/for2/for3/for4可以匹配到两个前缀，取更长的
		{6, "/for1/for2/", "", ""},      //5

		// 测试各种类型的静态路由分叉
		{7, "/baidu/test/path1", "", ""}, //6 /baidu/test/path1, /baidu/test/path123,
		{8, "/baidu/test/path", "", ""},  //7 /baidu/test/path, /baidu/test/path3
		{9, "/baidu/test/path2", "", ""}, //8 /baidu/test/path2, /baidu/test/path23
		{10, "/baidu/tee/path1", "", ""}, //9 /baidu/tee/path1, /baidu/tee/path1/path2
		{11, "/baidu/set/path1", "", ""}, //10 /baidu/set/path1, /baidu/set/path133

		// 测试域名的静态路由
		{12, ".com.bqq.huxing/baidu/set/path1", "", ""},  // 11
		{13, ".com.bqq.hotload/baidu/set/path1", "", ""}, // 12
		{14, ".cn.bqq.huxing/baidu/set/path1", "", ""},   // 13

		// 静态路由结尾添加和不添加/的测试
		{15, "/mytest/my", "", ""},  //14 完全匹配：/mytest/my，没有重定向：因为配置了/mytest/my/，前缀：/mytest/myname
		{16, "/mytest/my/", "", ""}, //15 完全匹配：/mytest/my/，没有重定向：因为配置了/mytest/my，前缀：/mytest/my/name
	}
	for i := range routeValMap {
		err := n.addRouter(routeValMap[i].route, http.HandleFunc(func() func(*http.Context) http.HttpPhase {
			v := routeValMap[i]
			return func(c *http.Context) http.HttpPhase {
				index = v.val
				routeStr = v.route
				return http.HttpFinish
			}
		}()))
		test.MustBeValue(t, nil, err, fmt.Sprintf("err of n.addRouter(%s)", routeValMap[i].route))
	}
	testIndexVal := map[string]resultInfo{
		"/router1/level1/test1":      {0, false, false},
		"/router1/level1/test123":    {0, true, false},
		"/router1/level1/test1/val1": {0, true, false},
		"/router1/level1/test1/":     {0, false, true},
		"/test/level/val1":           {1, false, true},
		"/test/level/val1/":          {1, false, false},
		"/test/level/val1/tee":       {1, true, false},
		"/test/level/val1/tee2":      {1, true, false},

		"/pre1/pre2/pre3":      {2, false, true},
		"/pre1/pre2":           {3, false, true},
		"/for1/for2/for3/for4": {4, true, false},

		"/baidu/test/path1":      {6, false, false},
		"/baidu/test/path123":    {6, true, false},
		"/baidu/test/path":       {7, false, false},
		"/baidu/test/path3":      {7, true, false},
		"/baidu/test/path2":      {8, false, false},
		"/baidu/test/path23":     {8, true, false},
		"/baidu/tee/path1":       {9, false, false},
		"/baidu/tee/path1/path2": {9, true, false},
		"/baidu/set/path1":       {10, false, false},
		"/baidu/set/path133":     {10, true, false},

		".com.bqq.huxing/baidu/set/path1":      {11, false, false},
		".com.bqq.huxing/baidu/set/path1/":     {11, false, true},
		".com.bqq.huxing/baidu/set/path123":    {11, true, false},
		".com.bqq.huxing/baidu/set/path1/test": {11, true, false},

		".com.bqq.hotload/baidu/set/path1":      {12, false, false},
		".com.bqq.hotload/baidu/set/path1/":     {12, false, true},
		".com.bqq.hotload/baidu/set/path123":    {12, true, false},
		".com.bqq.hotload/baidu/set/path1/test": {12, true, false},

		".cn.bqq.huxing/baidu/set/path1":      {13, false, false},
		".cn.bqq.huxing/baidu/set/path1/":     {13, false, true},
		".cn.bqq.huxing/baidu/set/path123":    {13, true, false},
		".cn.bqq.huxing/baidu/set/path1/test": {13, true, false},

		"/mytest/my":      {14, false, false},
		"/mytest/myname":  {14, true, false},
		"/mytest/my/":     {15, false, false},
		"/mytest/my/name": {15, true, false},
	}
	for path, value := range testIndexVal {
		nodeVal, err := n.getValue(path, nil, nil, false)
		test.MustBeValue(t, nil, err, fmt.Sprintf("err of n.getValue(%s, nil, nil, false)", path))
		test.Must(t, nodeVal.handler != nil, "nodeVal.handler for %s should not be nil, but go nil", path)
		test.MustBeValue(t, value.tsr, nodeVal.tsr, fmt.Sprintf("tsr for %s", path))
		test.MustBeValue(t, value.pre, nodeVal.pre, fmt.Sprintf("pre for %s", path))
		if nodeVal.handler != nil {
			nodeVal.handler.Handler(nil)
			test.MustBeValue(t, routeValMap[value.index].val, index, fmt.Sprintf("index for %s", path))
			test.MustBeValue(t, routeValMap[value.index].route, routeStr, fmt.Sprintf("routeStr for %s", path))
		}
	}
}

func TestRouterParams(t *testing.T) {
	n := &node{}
	var index int
	var routeStr string
	routeValMap := []routeVal{
		// 测试参数配置的样例
		// 一个参数且在后面
		{17, "/args/args1/:name", "", ""}, // 0 完全匹配：/args/args1/wang1, /args/args1/li, 重定向：/args/args1/chen/, 前缀：/args/args1/cheng/test
		// 匹配静态，不匹配模糊
		{18, "/args/args1/wang", "", ""}, // 1 /args/args1/wang, /args/args1/wang/, /args/args1/wang22/test /args/args1/wang/test
		{19, "/params/:user/", "", ""},   // 2 完全匹配：/params/user1/, /params/liu/, 重定向：/params/zhang, 前缀：/params/cheng/test
		// // // 匹配静态，不匹配模糊
		{20, "/params/lang/", "", ""}, // 3 /params/lang/, /params/lang,  /params/lang/33
		// // // 结尾添加/和不添加的匹配
		{21, "/room/room1/room2/:user", "", ""},  // 4
		{22, "/room/room1/room2/:user/", "", ""}, // 5

		// // 一个参数在最开头
		// {12, "/:path/args1/args2/args3", "", ""},
		// // 一个参数在中间
		// {12, "/path/params:/path2", "", ""},
		// //两个参数相连，且在后面
		// {13, "/sumit/user/:user/:value", "", ""},
		// //两个参数相连，且在前面
		// {13, "/:dir/:path/args1/args2"},
		// //两个参数相连，且在中间
		// {13, "/student/:id/:sex/room"},
		// // 两个参数且不相连，一前一后
		// // 两个参数且不相连，一前一非后
		// // 两个参数且不相连，一后一非前
		// // 两个参数且不相连，非前非后

		// // 相同的位置出现不同的参数
		// {12, "/load/:user/test1", "", ""}, // 11
		// {15, "/load/:name/test2", "", ""},
		// //不同的位置出现相同的参数
		// {12, "/pos/pos1/:user/pos2", "", ""}, // 11
		// {15, "/pos/:user/pos1/pos2", "", ""},
	}
	for i := range routeValMap {
		err := n.addRouter(routeValMap[i].route, http.HandleFunc(func() func(*http.Context) http.HttpPhase {
			v := routeValMap[i]
			return func(c *http.Context) http.HttpPhase {
				index = v.val
				routeStr = v.route
				return http.HttpFinish
			}
		}()))
		test.MustBeValue(t, nil, err, fmt.Sprintf("err of n.addRouter(%s)", routeValMap[i].route))
	}
	testIndexVal := map[string]resultInfo{
		// "/args/args1/wang":        {1, false, false},
		// "/args/args1/wang/":       {1, false, true},
		// "/args/args1/wang22/test": {1, true, false},
		// "/args/args1/wang/test":   {1, true, false},

		"/args/args1/wang1": {0, true, false},
		"/args/args1/li":    {0, false, false},
		// "/args/args1/chen/":      {0, false, true},
		// "/args/args1/cheng/test": {0, true, false},
	}
	for path, value := range testIndexVal {
		nodeVal, err := n.getValue(path, nil, nil, false)
		test.MustBeValue(t, nil, err, fmt.Sprintf("err of n.getValue(%s, nil, nil, false)", path))
		test.Must(t, nodeVal.handler != nil, "nodeVal.handler for %s should not be nil, but go nil", path)
		// test.MustBeValue(t, value.tsr, nodeVal.tsr, fmt.Sprintf("tsr for %s", path))
		// test.MustBeValue(t, value.pre, nodeVal.pre, fmt.Sprintf("pre for %s", path))
		if nodeVal.handler != nil {
			nodeVal.handler.Handler(nil)
			test.MustBeValue(t, routeValMap[value.index].val, index, fmt.Sprintf("index for %s", path))
			test.MustBeValue(t, routeValMap[value.index].route, routeStr, fmt.Sprintf("routeStr for %s", path))
		}
	}
}

func TestRouterCatchAll(t *testing.T) {
	n := &node{}
	var index int
	var routeStr string
	routeValMap := []routeVal{
		// 测试捕获全部的路由
		{1, "/upload/name/*val", "", ""}, // 0 完全匹配：/upload/name/test11val, 重定向：/upload/name/1testval/, 前缀：/upload/name/tttval/ttt, /upload/name/tttvalttt
		// {2, "/upload/name/testval", "", ""}, // 0 完全匹配：/upload/name/testval, 重定向：/upload/name/testval/, 前缀：/upload/name/testvala, /upload/name/testval/aa

	}
	for i := range routeValMap {
		err := n.addRouter(routeValMap[i].route, http.HandleFunc(func() func(*http.Context) http.HttpPhase {
			v := routeValMap[i]
			return func(c *http.Context) http.HttpPhase {
				index = v.val
				routeStr = v.route
				return http.HttpFinish
			}
		}()))
		test.MustBeValue(t, nil, err, fmt.Sprintf("err of n.addRouter(%s)", routeValMap[i].route))
	}
	testIndexVal := map[string]resultInfo{
		"/upload/name/testval":    {0, false, false},
		"/upload/name/testval/":   {0, false, true},
		"/upload/name/testvala":   {0, true, false},
		"/upload/name/testval/aa": {0, true, false},
	}
	for path, value := range testIndexVal {
		nodeVal, err := n.getValue(path, nil, nil, false)
		test.MustBeValue(t, nil, err, fmt.Sprintf("err of n.getValue(%s, nil, nil, false)", path))
		test.Must(t, nodeVal.handler != nil, "nodeVal.handler for %s should not be nil, but go nil", path)
		// test.MustBeValue(t, value.tsr, nodeVal.tsr, fmt.Sprintf("tsr for %s", path))
		// test.MustBeValue(t, value.pre, nodeVal.pre, fmt.Sprintf("pre for %s", path))
		if nodeVal.handler != nil {
			nodeVal.handler.Handler(nil)
			test.MustBeValue(t, routeValMap[value.index].val, index, fmt.Sprintf("index for %s", path))
			test.MustBeValue(t, routeValMap[value.index].route, routeStr, fmt.Sprintf("routeStr for %s", path))
		}
	}
}

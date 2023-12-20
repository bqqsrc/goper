//  Copyright (C) 晓白齐齐,版权所有.

package httpT

import (
	"fmt"
	"github.com/bqqsrc/bqqg/test"
	"github.com/bqqsrc/goper/config"
	"github.com/bqqsrc/goper/core"
	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/object"
	"strings"
	"testing"
)

type BaseTestHttp struct {
	callbackBuilder strings.Builder
}

func (b *BaseTestHttp) writeHead(keyValue object.ConfigPair, pc *object.Cycle) error {
	b.callbackBuilder.WriteString(keyValue.Key)
	b.callbackBuilder.WriteString(":{\n")
	return nil
}

func (b *BaseTestHttp) writeTail(keyValue object.ConfigPair, pc *object.Cycle) error {
	b.callbackBuilder.WriteString("};\n")
	return nil
}

type HttpTestExample1 struct {
	http.HttpComponent
	BaseTestHttp
	val1 string // httpCompt1Val1Test
	val2 int    // 999739
	val3 struct {
		Val1 string  `gson:"HttpTestExample1"` // HttpTestExample2Test
		Val2 float64 `gson:"HttpTestExample2"` //838383.8883
	}
}

type example1Server1 struct {
	Val1 bool   `gson:"Val1"`
	Val2 string `gson:"Val2"`
	Val3 int    `gson:"Val3"`
}

func (e *example1Server1) Equal(val any) bool {
	var v example1Server1
	ok := false
	if v, ok = val.(example1Server1); !ok {
		var vv *example1Server1
		if vv, ok = val.(*example1Server1); ok {
			v = *vv
		}
	}
	if ok {
		return e.Val1 == v.Val1 && e.Val2 == v.Val2 && e.Val3 == v.Val3
	}
	return false
}

type example1Location1 struct {
	Val1 bool
	Val2 string
}

func (e *example1Location1) Equal(val any) bool {
	var v example1Location1
	ok := false
	if v, ok = val.(example1Location1); !ok {
		var vv *example1Location1
		if vv, ok = val.(*example1Location1); ok {
			v = *vv
		}
	}
	if ok {
		return e.Val1 == v.Val1 && e.Val2 == v.Val2
	}
	return false
}

func (h *HttpTestExample1) CreateMainConfig(key string) http.HttpCommands {
	return http.HttpCommands{
		{
			object.ConfigPair{
				"httpCompt1Val1",
				&h.val1,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"httpCompt1Val2",
				&h.val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"httpCompt1Val3",
				&h.val3,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample1) CreateSrvConfig(key string) http.HttpCommands {
	var val1 bool
	var val3 int
	return http.HttpCommands{
		{
			object.ConfigPair{"srvCompt1Val1",
				&val1,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"srvCompt1Val2",
				&example1Server1{},
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"srvCompt1Val3",
				&val3,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample1) CreateLocConfig(key string) http.HttpCommands {
	var val2 float64
	return http.HttpCommands{
		{
			object.ConfigPair{"locCompt1Val1",
				&example1Location1{},
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"locCompt1Val2",
				&val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample1) CreateMthConfig(key string) http.HttpCommands {
	var val1 bool
	var val2 string
	return http.HttpCommands{
		{
			object.ConfigPair{"mthCompt1Val1",
				&val1,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"mthCompt1Val2",
				&val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

type HttpTestExample2 struct {
	http.HttpComponent
	BaseTestHttp
	val1       float64 // 7516.666
	val2       bool    // true
	val3       string  // HttpTestExample2Val3
	keywordStr []string
}

type example2Server1 struct {
	Val1 string `gson:"val1"`
	Val2 int    `gson:"val2"`
}

func (e *example2Server1) Equal(val any) bool {
	var v example2Server1
	ok := false
	if v, ok = val.(example2Server1); !ok {
		var vv *example2Server1
		if vv, ok = val.(*example2Server1); ok {
			v = *vv
		}
	}
	if ok {
		return e.Val1 == v.Val1 && e.Val2 == v.Val2
	}
	return false
}

type example2Location1 struct {
	Val1 bool
	Val2 string
}

func (e *example2Location1) Equal(val any) bool {
	var v example2Location1
	ok := false
	if v, ok = val.(example2Location1); !ok {
		var vv *example2Location1
		if vv, ok = val.(*example2Location1); ok {
			v = *vv
		}
	}
	if ok {
		return e.Val1 == v.Val1 && e.Val2 == v.Val2
	}
	return false
}

func (h *HttpTestExample2) CreateMainConfig(key string) http.HttpCommands {
	var keyword string
	return http.HttpCommands{
		{
			object.ConfigPair{"httpCompt2Val1",
				&h.val1,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"httpCompt2Val2",
				&h.val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"httpKeywords",
				&keyword,
			},
			true,
			nil,
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				if h.keywordStr == nil {
					h.keywordStr = make([]string, 0)
				}
				h.keywordStr = append(h.keywordStr, keyword)
				return nil
			},
		},
		{
			object.ConfigPair{"httpCompt2Val3",
				&h.val3,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample2) CreateSrvConfig(key string) http.HttpCommands {
	var val2 int
	return http.HttpCommands{
		{
			object.ConfigPair{"srvCompt2Val1",
				&example2Server1{},
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"srvCompt2Val2",
				&val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample2) CreateLocConfig(key string) http.HttpCommands {
	var val2 float64
	var val3 string
	return http.HttpCommands{
		{
			object.ConfigPair{"locCompt2Val1",
				&example2Location1{},
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"locCompt2Val2",
				&val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"locCompt2Val3",
				&val3,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample2) CreateMthConfig(key string) http.HttpCommands {
	var val1 int
	var val2 float64
	var val3 string
	return http.HttpCommands{
		{
			object.ConfigPair{"mthCompt2Val1",
				&val1,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"mthCompt2Val2",
				&val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"mthCompt2Val3",
				&val3,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

type HttpTestExample3 struct {
	http.HttpComponent
	BaseTestHttp
	val1 float64 // 943.666
	val2 int     // 9536
	val3 struct {
		Val1 int    // 388
		Val2 string //Val2Val2
		Val3 bool   // false
	}
}

type example3Server1 struct {
	Val1 string `gson:"val1"`
	Val2 int    `gson:"val2"`
	Val3 bool   `gson:"val3"`
}

func (e *example3Server1) Equal(val any) bool {
	var v example3Server1
	ok := false
	if v, ok = val.(example3Server1); !ok {
		var vv *example3Server1
		if vv, ok = val.(*example3Server1); ok {
			v = *vv
		}
	}
	if ok {
		return e.Val1 == v.Val1 && e.Val2 == v.Val2 && e.Val3 == v.Val3
	}
	return false
}

func (h *HttpTestExample3) CreateMainConfig(key string) http.HttpCommands {
	return http.HttpCommands{
		{
			object.ConfigPair{"httpCompt3Val1",
				&h.val1,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"httpCompt3Val2",
				&h.val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"httpCompt3Val3",
				&h.val3,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample3) CreateSrvConfig(key string) http.HttpCommands {
	return http.HttpCommands{
		{
			object.ConfigPair{"srvCompt3Val1",
				&example3Server1{},
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample3) CreateLocConfig(key string) http.HttpCommands {
	var val1 float64
	var val2 int
	return http.HttpCommands{
		{
			object.ConfigPair{"locCompt3Val1",
				&val1,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"locCompt3Val2",
				&val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

func (h *HttpTestExample3) CreateMthConfig(key string) http.HttpCommands {
	var val1 float64
	var val2 int
	var val3, val4 string
	return http.HttpCommands{
		{
			object.ConfigPair{"mthCompt3Val1",
				&val1,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"mthCompt3Val2",
				&val2,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"mthCompt3Val3",
				&val3,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
		{
			object.ConfigPair{"mthCompt3Val4",
				&val4,
			},
			false,
			h.writeHead,
			h.writeTail,
		},
	}
}

type Equaler interface {
	Equal(any) bool
}

func TestConfigStart(t *testing.T) {
	c := &config.Config{}
	ht1 := &HttpTestExample1{}
	ht2 := &HttpTestExample2{}
	ht3 := &HttpTestExample3{}
	pc := &object.Cycle{
		ConfFile: "./http.gs",
		Compts:   []object.Componenter{&core.Core{}, c, &http.Http{}, ht1, ht2, ht3},
	}
	err := pc.Compts[0].Start(pc)
	test.MustBeValue(t, nil, err, "err of c.Start(pc)")
	test.MustBeValue(t, "httpCompt1Val1Test", ht1.val1, "ht1.val1")
	test.MustBeValue(t, 999739, ht1.val2, "ht1.val2")
	test.MustBeValue(t, "HttpTestExample2Test", ht1.val3.Val1, "ht1.val3.Val1")
	test.MustBeValue(t, 838383.8883, ht1.val3.Val2, "ht1.val3.Val2")
	test.MustBeValue(t, "httpCompt1Val1:{\n};\nhttpCompt1Val2:{\n};\nhttpCompt1Val3:{\n};\nsrvCompt1Val1:{\n};\nsrvCompt1Val2:{\n};\nsrvCompt1Val3:{\n};\nlocCompt1Val1:{\n};\nmthCompt1Val1:{\n};\nmthCompt1Val2:{\n};\nlocCompt1Val2:{\n};\nlocCompt1Val1:{\n};\nlocCompt1Val2:{\n};\nsrvCompt1Val1:{\n};\nlocCompt1Val1:{\n};\nlocCompt1Val2:{\n};\nmthCompt1Val1:{\n};\nmthCompt1Val2:{\n};\nmthCompt1Val1:{\n};\nmthCompt1Val2:{\n};\nmthCompt1Val1:{\n};\nmthCompt1Val2:{\n};\n",
		ht1.callbackBuilder.String(), "ht1.callbackBuilder.String()")
	test.MustBeStringSlice(t, []string{"httpKeywords1", "httpKeywords2", "httpKeywords3", "httpKeywords4",
		"httpKeywords5", "httpKeywords6-1", "httpKeywords7", "httpKeywords8", "httpKeywords9", "httpKeywords10",
		"httpKeywords11", "httpKeywords11-1", "httpKeywords12", "httpKeywords14", "httpKeywords15", "httpKeywords15-1"}, ht2.keywordStr, "ht1.keywordStr")

	test.MustBeValue(t, 7516.666, ht2.val1, "ht2.val1")
	test.MustBeValue(t, true, ht2.val2, "ht2.val2")
	test.MustBeValue(t, "HttpTestExample2Val3", ht2.val3, "ht2.val3")
	test.MustBeValue(t, "httpCompt2Val1:{\n};\nhttpCompt2Val2:{\n};\nhttpCompt2Val3:{\n};\nsrvCompt2Val1:{\n};\nsrvCompt2Val2:{\n};\nmthCompt2Val1:{\n};\nmthCompt2Val3:{\n};\nmthCompt2Val1:{\n};\nmthCompt2Val2:{\n};\nlocCompt2Val1:{\n};\nlocCompt2Val2:{\n};\nlocCompt2Val3:{\n};\nmthCompt2Val1:{\n};\nlocCompt2Val1:{\n};\nlocCompt2Val3:{\n};\nlocCompt2Val2:{\n};\nmthCompt2Val2:{\n};\nmthCompt2Val1:{\n};\nmthCompt2Val3:{\n};\nsrvCompt2Val1:{\n};\nsrvCompt2Val2:{\n};\nsrvCompt2Val2:{\n};\nmthCompt2Val1:{\n};\nmthCompt2Val2:{\n};\nmthCompt2Val3:{\n};\n",
		ht2.callbackBuilder.String(), "ht2.callbackBuilder.String()")

	test.MustBeValue(t, 943.666, ht3.val1, "ht3.val1")
	test.MustBeValue(t, 9536, ht3.val2, "ht3.val2")
	test.MustBeValue(t, 388, ht3.val3.Val1, "ht3.val3.Val1")
	test.MustBeValue(t, "Val2Val2", ht3.val3.Val2, "ht3.val3.Val2")
	test.MustBeValue(t, false, ht3.val3.Val3, "ht3.val3.Val3")
	test.MustBeValue(t, "httpCompt3Val1:{\n};\nmthCompt3Val3:{\n};\nmthCompt3Val1:{\n};\nmthCompt3Val1:{\n};\nmthCompt3Val2:{\n};\nlocCompt3Val1:{\n};\nlocCompt3Val2:{\n};\nlocCompt3Val1:{\n};\nlocCompt3Val2:{\n};\nmthCompt3Val4:{\n};\nmthCompt3Val1:{\n};\nmthCompt3Val3:{\n};\nmthCompt3Val2:{\n};\nlocCompt3Val2:{\n};\nlocCompt3Val1:{\n};\nlocCompt3Val2:{\n};\nlocCompt3Val1:{\n};\nmthCompt3Val1:{\n};\nmthCompt3Val2:{\n};\nlocCompt3Val1:{\n};\nlocCompt3Val2:{\n};\nhttpCompt3Val2:{\n};\nhttpCompt3Val3:{\n};\nsrvCompt3Val1:{\n};\nlocCompt3Val1:{\n};\nlocCompt3Val2:{\n};\nmthCompt3Val1:{\n};\nmthCompt3Val4:{\n};\nmthCompt3Val3:{\n};\nmthCompt3Val2:{\n};\n",
		ht3.callbackBuilder.String(), "ht3.callbackBuilder.String()")

	test.Must(t, ht1.Configs == ht2.Configs, "ht1.Configs must equal ht2.Configs, but got ht1.Configs != ht2.Configs")
	test.Must(t, ht1.Configs == ht3.Configs, "ht1.Configs must equal ht3.Configs, but got ht1.Configs != ht3.Configs")

	trueVal := true
	falseVal := false
	int958 := 958
	int9563 := 9563
	inte78423 := -78423
	int7482 := 7482
	float983323 := 98.3323
	float459641336 := 459641.336
	int1569321 := 1569321
	float7589934 := 75899.336
	strsrv1loc2comp2 := "locCompt2Val3"
	float123666 := 12.3666
	int12854 := 12854
	float84233 := 842.33
	inte1569 := -1569
	floate96256 := -96.256
	strsrv2loc1comp2 := "locCompt2Val3Value"
	float9533 := 95.33
	inte78523 := -78523
	floate153666 := -153.666
	float746333 := 746.333
	int95423 := 95423
	floate7853666 := -785.3666
	var floate753 float64 = -753
	inte458 := -458
	strsrv1loc1mth1comp3 := "mthCompt3Val3125"
	strsrv1loc1mth1comp1 := "mthCompt1Val200"
	inte9696 := -9696
	float87336 := 87.336
	floate86366 := -86.366
	int8524 := 8524
	strsrv1loc1mth2comp2 := "mthCompt2Val3333"
	int9756 := 9756
	floate87263 := -87.263
	strsrv1loc2mth2comp3_1 := "mthCompt3Val4156"
	strsrv1loc2mth2comp3_2 := "mthCompt3Val3123"
	float78369 := 78.369
	inte48569 := -48569
	int7526 := 7526
	floate8536412 := -85.36412
	floate88369 := -88.369
	int7598624 := 7598624
	inte96954 := -96954
	strsrv2loc1mth1comp2 := "mthCompt2Val3111"
	strsrv3loc1mth1comp1 := "mthCompt1Val211"
	inte7896 := -7896
	float85336 := 85.336
	strsrv3loc1mth1comp2 := "mthCompt2Val3222"
	strsrv3loc1mth2comp3_1 := "mthCompt3Val4853"
	strsrv3loc1mth2comp3_2 := "mthCompt3Val3854"
	int96312 := 96312
	floate9633 := -96.33
	strsrv3loc1mth2comp1 := "mthCompt1Val222"
	strsrv3loc1mth3comp1 := "mthCompt1Val233"

	configRet := &http.AllHttpConfigs{
		[]http.HttpConfigs{
			{
				{
					Key:   "httpCompt1Val1",
					Value: &ht1.val1,
				},
				{Key: "httpCompt1Val2",
					Value: &ht1.val2,
				},
				{
					Key:   "httpCompt1Val3",
					Value: &ht1.val3,
				},
			},
			{
				{Key: "httpCompt2Val1",
					Value: &ht2.val1,
				},
				{
					Key:   "httpCompt2Val2",
					Value: &ht2.val2,
				},
				{
					Key:   "httpCompt2Val3",
					Value: &ht2.val3,
				},
			},
			{
				{Key: "httpCompt3Val1",
					Value: &ht3.val1,
				},
				{
					Key:   "httpCompt3Val2",
					Value: &ht3.val2,
				},
				{
					Key:   "httpCompt3Val3",
					Value: &ht3.val3,
				},
			},
		},
		[][]http.HttpConfigs{
			{
				{
					{
						Key:   "srvCompt1Val1",
						Value: &trueVal,
					},
					{
						Key:   "srvCompt1Val2",
						Value: &example1Server1{true, "srvCompt1Val1Val2", 9823},
					},
					{
						Key:   "srvCompt1Val3",
						Value: &int958,
					},
				},
				{
					{
						Key:   "srvCompt2Val1",
						Value: &example2Server1{"srvCompt2Val1val1", 78596},
					},
					{
						Key:   "srvCompt2Val2",
						Value: &int9563,
					},
				},
				nil,
			},
			{
				nil,
				{
					{
						Key:   "srvCompt2Val1",
						Value: &example2Server1{"srvCompt2Val1val1", 8323},
					},
					{
						Key:   "srvCompt2Val2",
						Value: &inte78423,
					},
				},
				nil,
			},
			{
				{
					{
						Key:   "srvCompt1Val1",
						Value: &falseVal,
					},
				},
				{
					{
						Key:   "srvCompt2Val2",
						Value: &int7482,
					},
				},
				{
					{

						Key:   "srvCompt3Val1",
						Value: &example3Server1{"srvCompt3Val1val1", 888, true},
					},
				},
			},
		},
		[][][]http.HttpConfigs{
			{
				{
					{
						{
							Key:   "locCompt1Val1",
							Value: &example1Location1{false, "locCompt1Val1"},
						},
						{
							Key:   "locCompt1Val2",
							Value: &float983323,
						},
					},
					nil,
					{
						{
							Key:   "locCompt3Val1",
							Value: &float459641336,
						},
						{
							Key:   "locCompt3Val2",
							Value: &int1569321,
						},
					},
				},
				{
					nil,
					{
						{
							Key:   "locCompt2Val1",
							Value: &example2Location1{true, "locCompt2Val1Val2"},
						},
						{
							Key:   "locCompt2Val2",
							Value: &float7589934,
						},
						{
							Key:   "locCompt2Val3",
							Value: &strsrv1loc2comp2,
						},
					},
					{
						{
							Key:   "locCompt3Val1",
							Value: &float123666,
						},
						{
							Key:   "locCompt3Val2",
							Value: &int12854,
						},
					},
				},
				{
					nil,
					nil,
					{
						{
							Key:   "locCompt3Val2",
							Value: &inte1569,
						},
						{
							Key:   "locCompt3Val1",
							Value: &float84233,
						},
					},
				},
			},
			{
				{
					nil,
					{
						{
							Key:   "locCompt2Val1",
							Value: &example2Location1{false, "locCompt2Val1"},
						},
						{
							Key:   "locCompt2Val3",
							Value: &strsrv2loc1comp2,
						},
						{
							Key:   "locCompt2Val2",
							Value: &floate96256,
						},
					},
					{
						{
							Key:   "locCompt3Val2",
							Value: &inte78523,
						},
						{
							Key:   "locCompt3Val1",
							Value: &float9533,
						},
					},
				},
				{
					{
						{
							Key:   "locCompt1Val1",
							Value: &example1Location1{true, "locCompt1Val1Val2"},
						},
						{
							Key:   "locCompt1Val2",
							Value: &floate153666,
						},
					},
					nil,
					{
						{
							Key:   "locCompt3Val1",
							Value: &float746333,
						},
						{
							Key:   "locCompt3Val2",
							Value: &int95423,
						},
					},
				},
			},
			{
				{
					{
						{
							Key:   "locCompt1Val1",
							Value: &example1Location1{true, "locCompt1Val1serverloc1"},
						},
						{
							Key:   "locCompt1Val2",
							Value: &floate7853666,
						},
					},
					nil,
					{
						{
							Key:   "locCompt3Val1",
							Value: &floate753,
						},
						{
							Key:   "locCompt3Val2",
							Value: &inte458,
						},
					},
				},
			},
		},
		[][][][]http.HttpConfigs{
			{
				{
					{
						{
							{
								Key:   "mthCompt1Val1",
								Value: &trueVal,
							},
							{
								Key:   "mthCompt1Val2",
								Value: &strsrv1loc1mth1comp1,
							},
						},
						{
							{
								Key:   "mthCompt2Val1",
								Value: &inte9696,
							},
						},
						{
							{
								Key:   "mthCompt3Val3",
								Value: &strsrv1loc1mth1comp3,
							},
							{
								Key:   "mthCompt3Val1",
								Value: &float87336,
							},
						},
					},
					{
						nil,
						{
							{
								Key:   "mthCompt2Val3",
								Value: &strsrv1loc1mth2comp2,
							},
							{
								Key:   "mthCompt2Val1",
								Value: &int9756,
							},
							{
								Key:   "mthCompt2Val2",
								Value: &floate87263,
							},
						},
						{
							{
								Key:   "mthCompt3Val1",
								Value: &floate86366,
							},
							{
								Key:   "mthCompt3Val2",
								Value: &int8524,
							},
						},
					},
				},
				{
					{
						nil,
						{
							{
								Key:   "mthCompt2Val1",
								Value: &inte48569,
							},
						},
						{
							{
								Key:   "mthCompt3Val4",
								Value: &strsrv1loc2mth2comp3_1,
							},
							{
								Key:   "mthCompt3Val1",
								Value: &float78369,
							},
							{
								Key:   "mthCompt3Val3",
								Value: &strsrv1loc2mth2comp3_2,
							},
							{
								Key:   "mthCompt3Val2",
								Value: &int7526,
							},
						},
					},
				},
				nil,
			},
			{
				{
					{
						nil,
						{
							{
								Key:   "mthCompt2Val2",
								Value: &floate8536412,
							},
							{
								Key:   "mthCompt2Val1",
								Value: &int7598624,
							},
							{
								Key:   "mthCompt2Val3",
								Value: &strsrv2loc1mth1comp2,
							},
						},
						{
							{
								Key:   "mthCompt3Val1",
								Value: &floate88369,
							},
							{
								Key:   "mthCompt3Val2",
								Value: &inte96954,
							},
						},
					},
				},
				nil,
			},
			{
				{
					{
						{
							{
								Key:   "mthCompt1Val1",
								Value: &trueVal,
							},
							{
								Key:   "mthCompt1Val2",
								Value: &strsrv3loc1mth1comp1,
							},
						},
						{
							{
								Key:   "mthCompt2Val1",
								Value: &inte7896,
							},
							{
								Key:   "mthCompt2Val2",
								Value: &float85336,
							},
							{
								Key:   "mthCompt2Val3",
								Value: &strsrv3loc1mth1comp2,
							},
						},
						nil,
					},
					{
						{
							{
								Key:   "mthCompt1Val1",
								Value: &falseVal,
							},
							{
								Key:   "mthCompt1Val2",
								Value: &strsrv3loc1mth2comp1,
							},
						},
						nil,
						{
							{
								Key:   "mthCompt3Val1",
								Value: &floate9633,
							},
							{
								Key:   "mthCompt3Val4",
								Value: &strsrv3loc1mth2comp3_1,
							},
							{
								Key:   "mthCompt3Val3",
								Value: &strsrv3loc1mth2comp3_2,
							},
							{
								Key:   "mthCompt3Val2",
								Value: &int96312,
							},
						},
					},
					{
						{
							{
								Key:   "mthCompt1Val1",
								Value: &falseVal,
							},
							{
								Key:   "mthCompt1Val2",
								Value: &strsrv3loc1mth3comp1,
							},
						},
						nil,
						nil,
					},
				},
			},
		},
	}

	isCommandEqual := func(ret1, ret2 object.ConfigPair) bool {
		if ret1.Key != ret2.Key {
			return false
		}
		return test.Equal(ret1.Value, ret2.Value)
	}

	isHttpCommandsEqual := func(ret1, ret2 http.HttpConfigs) bool {
		if ret1 == nil && ret2 == nil {
			return true
		}
		if ret1 == nil && ret2 != nil {
			return false
		}
		if ret1 != nil && ret2 == nil {
			return false
		}
		len1 := len(ret1)
		len2 := len(ret2)
		if len1 != len2 {
			return false
		}
		for index := 0; index < len1; index++ {
			if !isCommandEqual(ret1[index], ret2[index]) {
				return false
			}
		}
		return true
	}

	isHttpCommandsSliceEqual := func(ret1, ret2 any) bool {
		if ret1 == nil && ret2 == nil {
			return true
		}
		if ret1 == nil && ret2 != nil {
			return false
		}
		if ret1 != nil && ret2 == nil {
			return false
		}
		comands1, ok := ret1.([]http.HttpConfigs)
		if !ok {
			return false
		}
		comands2, ok := ret2.([]http.HttpConfigs)
		if !ok {
			return false
		}
		len1 := len(comands1)
		len2 := len(comands2)
		if len1 != len2 {
			return false
		}
		for index := 0; index < len1; index++ {
			if !isHttpCommandsEqual(comands1[index], comands2[index]) {
				return false
			}
		}
		return true
	}

	test.MustEqual(t, configRet.MainConfig, ht1.Configs.MainConfig,
		"ht1.Configs.MainConfig", "configRet.MainConfig", isHttpCommandsSliceEqual)

	for index := range ht1.Configs.SrvConfig {
		test.MustEqual(t, configRet.SrvConfig[index], ht1.Configs.SrvConfig[index],
			fmt.Sprintf("ht1.Configs.SrvConfig[%d]", index), fmt.Sprintf("configRet.SrvConfig[%d]", index), isHttpCommandsSliceEqual)
	}
	for i := 0; i < len(ht1.Configs.LocConfig); i++ {
		for j := 0; j < len(ht1.Configs.LocConfig[i]); j++ {
			test.MustEqual(t, configRet.LocConfig[i][j], ht1.Configs.LocConfig[i][j],
				fmt.Sprintf("ht1.Configs.LocConfig[%d][%d]", i, j), fmt.Sprintf("configRet.LocConfig[%d][%d]", i, j), isHttpCommandsSliceEqual)

		}
	}
	for i := 0; i < len(ht1.Configs.MthConfig); i++ {
		for j := 0; j < len(ht1.Configs.MthConfig[i]); j++ {
			for k := 0; k < len(ht1.Configs.MthConfig[i][j]); k++ {
				test.MustEqual(t, configRet.MthConfig[i][j][k], ht1.Configs.MthConfig[i][j][k],
					fmt.Sprintf("ht1.Configs.MthConfig[%d][%d][%d]", i, j, k), fmt.Sprintf("configRet.MthConfig[%d][%d][%d]", i, j, k), isHttpCommandsSliceEqual)
			}
		}
	}
}

func BenchmarkConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {

		pc := &object.Cycle{
			ConfFile: "./http.gs",
			Compts:   []object.Componenter{&core.Core{}, &config.Config{}, &http.Http{}, &HttpTestExample1{}, &HttpTestExample2{}, &HttpTestExample3{}},
		}
		pc.Compts[0].Start(pc)
	}
}

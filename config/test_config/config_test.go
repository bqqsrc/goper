//  Copyright (C) 晓白齐齐,版权所有.

package configT

import (
	"fmt"
	"github.com/bqqsrc/bqqg/test"
	"github.com/bqqsrc/goper/config"
	"github.com/bqqsrc/goper/core"
	"github.com/bqqsrc/goper/object"
	"strings"
	"testing"
)

type ConfigTestExample struct {
	object.BaseComponent
	val1 string
	val2 int
	val3 struct {
		Val1 string  `gson:"TestExampleVal1"`
		Val2 float64 `gson:"TestExampleVal2"`
	}
	val4            []float64
	val5            bool
	callbackBuilder strings.Builder
}

func (t *ConfigTestExample) CreateConfig(key string) []object.Command {
	return []object.Command{
		{
			object.ConfigPair{"TestExample1val1", // 要关注的key1
				&t.val1, // 要解析的值
			},
			false,
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				t.callbackBuilder.WriteString("key:")
				t.callbackBuilder.WriteString(keyValue.Key)
				return nil
			},
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				t.callbackBuilder.WriteString("\nvalue:")
				t.callbackBuilder.WriteString(*keyValue.Value.(*string))
				return nil
			},
		},
		{
			object.ConfigPair{"TestExample1val2",
				&t.val2,
			},
			false,
			nil,
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				t.callbackBuilder.WriteString("\nkey:")
				t.callbackBuilder.WriteString(keyValue.Key)
				t.callbackBuilder.WriteString("\nvalue:")
				t.callbackBuilder.WriteString(fmt.Sprintf("%v", *keyValue.Value.(*int)))
				return nil
			},
		},
		{
			object.ConfigPair{"TestExample1val3",
				&t.val3,
			},
			false,
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				t.callbackBuilder.WriteString("\nkey:")
				t.callbackBuilder.WriteString(keyValue.Key)
				return nil
			},
			nil,
		},
		{
			object.ConfigPair{"TestExample1val4",
				&t.val4,
			},
			false,
			nil,
			nil,
		},
		{
			object.ConfigPair{"TestExample1val5",
				&t.val5,
			},
			false,
			nil,
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				t.callbackBuilder.WriteString("\nkey:")
				t.callbackBuilder.WriteString(keyValue.Key)
				t.callbackBuilder.WriteString("\nvalue:")
				t.callbackBuilder.WriteString(fmt.Sprintf("%v", *keyValue.Value.(*bool)))
				return nil
			},
		},
	}
}

func (t *ConfigTestExample) InitConfig(pc *object.Cycle) error {
	t.callbackBuilder.WriteString("\nInitConfig")
	return nil
}

type ConfigTestExample2 struct {
	object.BaseComponent
	val1       uint
	val2       int32
	val3       []string
	val4       bool
	newStrVal4 []string
}

func (t *ConfigTestExample2) CreateConfig(key string) []object.Command {
	return []object.Command{
		{
			object.ConfigPair{"TestExample2val1", // 要关注的key1
				&t.val1, // 要解析的值
			},
			false,
			nil,
			nil,
		},
		{
			object.ConfigPair{"TestExample2val2",
				&t.val2,
			},
			false,
			nil,
			nil,
		},
		{
			object.ConfigPair{"TestExample2val3",
				&t.val3,
			},
			false,
			nil,
			nil,
		},
		{
			object.ConfigPair{"TestExample2val4",
				&t.val4,
			},
			false,
			nil,
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				if !t.val4 {
					t.newStrVal4 = append(t.val3, fmt.Sprintf("%d", t.val1), fmt.Sprintf("%d", t.val2), fmt.Sprintf("%v", t.val4))
				}
				return nil
			},
		},
	}
}

type ConfigTestExample3 struct {
	object.BaseComponent
	Val1     uint8   `gson:"Val1"`
	Val2     float32 `gson:"Val2"`
	Val3     []bool  `gson:"Val3"`
	Val4     string  `gson:"Val4"`
	keyword  string
	keywords []string
}

func (t *ConfigTestExample3) CreateConfig(key string) []object.Command {
	return []object.Command{
		{
			object.ConfigPair{"TestExample3", // 要关注的key1
				t, // 要解析的值
			},
			false,
			nil,
			nil,
		},
		{
			object.ConfigPair{"KeyWords", // 要关注的key1
				&t.keyword, // 要解析的值
			},
			true,
			nil,
			func(keyValue object.ConfigPair, pc *object.Cycle) error {
				if t.keywords == nil {
					t.keywords = make([]string, 0)
				}
				t.keywords = append(t.keywords, *keyValue.Value.(*string))
				return nil
			},
		},
	}
}

type ConfigTestExample4 struct {
	object.BaseComponent
	countTestExample4val1 int
	isInBlock             bool
	blockKeyBuild         strings.Builder
	isInTestExample1val3  bool
	structVal1            struct {
		Value1 int    // 333
		Value2 string // value2
		Value3 bool   // true
	}
	val1 float64 //99.88
	val2 string  // TestExampleVal2
	val3 float64 //8563.333
}

func (t *ConfigTestExample4) CreateConfig(key string) []object.Command {
	writeHead := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		if t.isInBlock {
			t.blockKeyBuild.WriteString(keyValue.Key)
			t.blockKeyBuild.WriteString(":{")
		}
		return nil
	}
	writeTail := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		if t.isInBlock {
			t.blockKeyBuild.WriteString("}")
		}
		return nil
	}
	switch key {
	case "", "TestExample4val1":
		return []object.Command{
			{
				object.ConfigPair{"TestExample4val1", // 要关注的key1
					nil, // 要解析的值
				},
				false,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					writeHead(keyValue, pc)
					t.countTestExample4val1++
					if t.countTestExample4val1 > 0 {
						t.isInBlock = true
					}
					return nil
				},
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					t.countTestExample4val1--
					if t.countTestExample4val1 == 0 {
						t.isInBlock = false
					}
					writeTail(keyValue, pc)
					return nil
				},
			},
		}
	case "TestValue1":
		return []object.Command{
			{
				object.ConfigPair{"TestValue1", // 要关注的key1
					&t.structVal1, // 要解析的值
				},
				false,
				writeHead,
				writeTail,
			},
		}
	case "TestValue2":
		return []object.Command{
			{
				object.ConfigPair{"TestValue2", // 要关注的key1
					&t.val1, // 要解析的值
				},
				false,
				writeHead,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					t.blockKeyBuild.WriteString(fmt.Sprintf("%f", *keyValue.Value.(*float64)))
					writeTail(keyValue, pc)
					return nil
				},
			},
		}
	case "TestExampleVal1":
		if !t.isInTestExample1val3 {
			return []object.Command{
				{
					object.ConfigPair{"TestExampleVal1", // 要关注的key1
						&t.val2, // 要解析的值
					},
					false,
					writeHead,
					func(keyValue object.ConfigPair, pc *object.Cycle) error {
						t.blockKeyBuild.WriteString(t.val2)
						writeTail(keyValue, pc)
						return nil
					},
				},
			}
		} else {
			return nil
		}
	case "TestExampleVal2":
		if t.isInTestExample1val3 {
			return []object.Command{
				{
					object.ConfigPair{"TestExampleVal2", // 要关注的key1
						&t.val3, // 要解析的值
					},
					false,
					writeHead,
					func(keyValue object.ConfigPair, pc *object.Cycle) error {
						t.blockKeyBuild.WriteString(fmt.Sprintf("%f", t.val3))
						writeTail(keyValue, pc)
						return nil
					},
				},
			}
		} else {
			return nil
		}
	case "TestExample1val3":
		return []object.Command{
			{
				object.ConfigPair{"TestExample1val3", // 要关注的key1
					nil, // 要解析的值
				},
				false,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					writeHead(keyValue, pc)
					t.isInTestExample1val3 = true
					return nil
				},
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					writeTail(keyValue, pc)
					t.isInTestExample1val3 = false
					return nil
				},
			},
		}
	default:
		return nil
	}
}

func (t *ConfigTestExample4) InitConfig(pc *object.Cycle) error {
	t.blockKeyBuild.WriteString("InitConfig")
	return nil
}

type ConfigTestInclude struct {
	object.BaseComponent
	isInBlock      bool
	blockKeyBuild  strings.Builder
	includeValue1  int     // 988
	includeValue2  float32 // 8523.66
	includeValue3  string  // includeVal3
	testValueVal1  uint    // 965821
	testValueVal2  string  // TestValueVal2_3
	isInTestValue1 bool
	structVal1     struct {
		Value1 int    // 952
		Value2 string // TestValue3Val2
		Value3 bool   // true
	}
	structVal2 struct {
		Value1 int     `gson:"Val1"` // -897
		Value2 float32 `gson:"Val2"` // -853.66
		Value3 string  `gson:"Val3"` // Value3Val
	}
}

func (t *ConfigTestInclude) CreateConfig(key string) []object.Command {
	writeHead := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		if t.isInBlock {
			t.blockKeyBuild.WriteString(keyValue.Key)
			t.blockKeyBuild.WriteString(":{")
		}
		return nil
	}
	writeTail := func(keyValue object.ConfigPair, pc *object.Cycle) error {
		if t.isInBlock {
			t.blockKeyBuild.WriteString("};")
		}
		return nil
	}

	switch key {
	case "", "TestInclude1", "includeVal1", "includeVal2", "includeVal3":
		return []object.Command{
			{
				object.ConfigPair{"TestInclude1", // 要关注的key1
					nil, // 要解析的值
				},
				false,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					t.isInBlock = true
					writeHead(keyValue, pc)
					return nil
				},
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					writeTail(keyValue, pc)
					t.isInBlock = false
					return nil
				},
			},
			{
				object.ConfigPair{"includeVal1", // 要关注的key1
					&t.includeValue1, // 要解析的值
				},
				false,
				nil,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					t.blockKeyBuild.WriteString(keyValue.Key)
					t.blockKeyBuild.WriteString(":")
					t.blockKeyBuild.WriteString(fmt.Sprintf("%d", t.includeValue1))
					t.blockKeyBuild.WriteString(";")
					return nil
				},
			},
			{
				object.ConfigPair{"includeVal2", // 要关注的key1
					&t.includeValue2, // 要解析的值
				},
				false,
				nil,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					t.blockKeyBuild.WriteString(keyValue.Key)
					t.blockKeyBuild.WriteString(":")
					t.blockKeyBuild.WriteString(fmt.Sprintf("%f", t.includeValue2))
					t.blockKeyBuild.WriteString(";")
					return nil
				},
			},
			{
				object.ConfigPair{"includeVal3", // 要关注的key1
					&t.includeValue3, // 要解析的值
				},
				false,
				nil,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					t.blockKeyBuild.WriteString(keyValue.Key)
					t.blockKeyBuild.WriteString(":")
					t.blockKeyBuild.WriteString(t.includeValue3)
					t.blockKeyBuild.WriteString(";")
					return nil
				},
			},
		}
	case "TestValue3":
		return []object.Command{
			{
				object.ConfigPair{"TestValue1", // 要关注的key1
					&t.structVal1, // 要解析的值
				},
				false,
				writeHead,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					return writeTail(keyValue, pc)

				},
			},
		}
	case "TestValue1":
		return []object.Command{
			{
				object.ConfigPair{"TestValue1", // 要关注的key1
					nil, // 要解析的值
				},
				false,
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					t.isInTestValue1 = true
					writeHead(keyValue, pc)
					return nil
				},
				func(keyValue object.ConfigPair, pc *object.Cycle) error {
					t.isInTestValue1 = false
					writeTail(keyValue, pc)
					return nil
				},
			},
		}
	case "TestValue2":
		return []object.Command{
			{
				object.ConfigPair{"TestValue2", // 要关注的key1
					nil, // 要解析的值
				},
				false,
				writeHead,
				writeTail,
			},
		}
	case "TestValueVal1":
		var returnVal any
		if t.isInTestValue1 {
			returnVal = &t.testValueVal1
		}
		return []object.Command{
			{
				object.ConfigPair{"TestValueVal1", // 要关注的key1
					returnVal, // 要解析的值
				},
				false,
				writeHead,
				writeTail,
			},
		}
	case "TestValueVal2":
		var returnVal any
		if !t.isInTestValue1 {
			returnVal = &t.testValueVal2
		}
		return []object.Command{
			{
				object.ConfigPair{"TestValueVal2", // 要关注的key1
					returnVal, // 要解析的值
				},
				false,
				writeHead,
				writeTail,
			},
		}
	case "TestValue4":
		return []object.Command{
			{
				object.ConfigPair{"TestValue4", // 要关注的key1
					&t.structVal2, // 要解析的值
				},
				false,
				writeHead,
				writeTail,
			},
		}
	default:
		return nil
	}
}

func (t *ConfigTestInclude) InitConfig(pc *object.Cycle) error {
	t.blockKeyBuild.WriteString("InitConfig")
	return nil
}

func TestConfigStart(t *testing.T) {
	c := &config.Config{}
	ct := &ConfigTestExample{}
	ct1 := &ConfigTestExample2{}
	ct2 := &ConfigTestExample3{}
	ct3 := &ConfigTestExample4{}
	ct4 := &ConfigTestInclude{}
	pc := &object.Cycle{
		ConfFile: "./config.gs",
		Compts:   []object.Componenter{&core.Core{}, c, ct, ct1, ct2, ct3, ct4},
	}
	err := pc.Compts[0].Start(pc)

	test.MustBeValue(t, nil, err, "err of c.Start(pc)")

	test.MustBeValue(t, "TestExample1", ct.val1, "ct.val1")
	test.MustBeValue(t, 999823, ct.val2, "ct.val2")
	test.MustBeValue(t, "TestExampleVal1", ct.val3.Val1, "ct.val3.Val1")
	test.MustBeValue(t, 8563.333, ct.val3.Val2, "ct.val3.Val2")
	test.MustBeFloat64Slice(t, []float64{986.33, 36.66, 589.33, 4256.33}, ct.val4, "ct.val4")
	test.MustBeValue(t, true, ct.val5, "ct.val5")
	ctBuilderShouldBe := `key:TestExample1val1
value:TestExample1
key:TestExample1val2
value:999823
key:TestExample1val3
key:TestExample1val5
value:true
InitConfig`
	ctBuilderBe := ct.callbackBuilder.String()
	test.MustBeValue(t, ctBuilderShouldBe, ctBuilderBe, "ct.callbackBuilder.String()")

	test.MustBeValue(t, uint(896359), ct1.val1, "ct1.val1")
	test.MustBeValue(t, int32(-96354), ct1.val2, "ct1.val2")
	ct1Val3ShouldBe := []string{"TestExample2val3_1", "TestExample2val3_2", "TestExample2val3_3"}
	test.MustBeStringSlice(t, ct1Val3ShouldBe, ct1.val3, "ct1.val3")
	test.MustBeValue(t, false, ct1.val4, "ct1.val4")
	ct1Val3ShouldBe = append(ct1Val3ShouldBe, "896359", "-96354", "false")
	test.MustBeStringSlice(t, ct1Val3ShouldBe, ct1.newStrVal4, "ct1.newStrVal4")

	test.MustBeValue(t, uint8(110), ct2.Val1, "ct2.Val1")
	test.MustBeValue(t, float32(33.8563), ct2.Val2, "ct2.Val2")
	test.MustBeBoolSlice(t, []bool{true, false, false, true, false, true}, ct2.Val3, "ct2.Val3")
	test.MustBeValue(t, "ConfigTestExample3", ct2.Val4, "ct2.Val4")
	test.MustBeValue(t, "KeyWords8", ct2.keyword, "ct2.keyword")
	test.MustBeStringSlice(t, []string{"KeyWords1", "KeyWords3", "KeyWords4", "KeyWords6", "KeyWords7", "KeyWords8"}, ct2.keywords, "ct2.keywords")

	test.MustBeValue(t, 333, ct3.structVal1.Value1, "ct3.structVal1.Value1")
	test.MustBeValue(t, "value2", ct3.structVal1.Value2, "ct3.structVal1.Value2")
	test.MustBeValue(t, true, ct3.structVal1.Value3, "ct3.structVal1.Value3")
	test.MustBeValue(t, 99.88, ct3.val1, "ct3.val1")
	test.MustBeValue(t, "TestExampleVal2", ct3.val2, "ct3.val2")
	test.MustBeValue(t, 8563.333, ct3.val3, "ct3.val3")
	blockStr := fmt.Sprintf("TestExample4val1:{TestValue1:{}TestValue2:{%f}}TestExampleVal1:{%s}TestExample1val3:{TestExampleVal2:{%f}}InitConfig", 99.88, "TestExampleVal2", 8563.333)
	test.MustBeValue(t, blockStr, ct3.blockKeyBuild.String(), "ct3.blockKeyBuild.String()")

	test.MustBeValue(t, 988, ct4.includeValue1, "ct4.includeValue1")
	test.MustBeValue(t, float32(8523.66), ct4.includeValue2, "ct4.includeValue2")
	test.MustBeValue(t, "includeVal3", ct4.includeValue3, "ct4.includeValue3")
	test.MustBeValue(t, uint(965821), ct4.testValueVal1, "ct4.testValueVal1")
	test.MustBeValue(t, "TestValueVal2_3", ct4.testValueVal2, "ct4.testValueVal2")
	test.MustBeValue(t, 952, ct4.structVal1.Value1, "ct4.structVal1.Value1")
	test.MustBeValue(t, "TestValue3Val2", ct4.structVal1.Value2, "ct4.structVal1.Value2")
	test.MustBeValue(t, true, ct4.structVal1.Value3, "ct4.structVal1.Value3")
	test.MustBeValue(t, -897, ct4.structVal2.Value1, "ct4.structVal2.Value1")
	test.MustBeValue(t, float32(-853.66), ct4.structVal2.Value2, "ct4.structVal2.Value2")
	test.MustBeValue(t, "Value3Val", ct4.structVal2.Value3, "ct4.structVal2.Value3")

	blockStr = fmt.Sprintf("includeVal1:%d;includeVal2:%f;includeVal3:%s;TestInclude1:{TestValue1:{TestValueVal1:{};TestValueVal2:{};};TestValue2:{TestValueVal1:{};TestValueVal2:{};};TestValue3:{};TestValue4:{};};InitConfig", ct4.includeValue1, ct4.includeValue2, ct4.includeValue3)
	test.MustBeValue(t, blockStr, ct4.blockKeyBuild.String(), "ct4.blockKeyBuild.String()")
}

func BenchmarkConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {

		pc := &object.Cycle{
			ConfFile: "./config.gs",
			Compts:   []object.Componenter{&core.Core{}, &config.Config{}, &ConfigTestExample{}, &ConfigTestExample2{}, &ConfigTestExample3{}, &ConfigTestExample4{}, &ConfigTestInclude{}},
		}
		pc.Compts[0].Start(pc)
	}
}

//Test Invalid
//Test 重复定义关键字，重复定义键，关键字和键同名，同组件、不同组件

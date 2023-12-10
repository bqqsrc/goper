//  Copyright (C) 晓白齐齐,版权所有.

package core

import (
	"github.com/bqqsrc/goper/object"
	"strings"
	"testing"
	"time"
)

type CoreTestExample struct {
	object.BaseComponent
	TestResult strings.Builder
}

func (c *CoreTestExample) Start(pc *object.Cycle) error {
	c.TestResult.WriteString("TestResult Start\n")
	return nil
}

func (c *CoreTestExample) Update(pc *object.Cycle) (time.Duration, error) {
	c.TestResult.WriteString("TestResult Update\n")
	return 0, nil
}

func (c *CoreTestExample) Awake() error {
	c.TestResult.WriteString("TestResult Awake\n")
	return nil
}

func TestCore(t *testing.T) {
	ct := &CoreTestExample{}
	ct1 := &CoreTestExample{}
	ct2 := &CoreTestExample{}
	pc := &object.Cycle{
		Compts: []object.Componenter{&Core{}, ct, ct1, ct2},
	}
	err := pc.Compts[0].Start(pc)
	must(t, err == nil, "err of c.Start(pc) should be nil, but got %v", err)
	testMustResult := `TestResult Awake
TestResult Start
TestResult Update
`
	testGotResult := ct.TestResult.String()
	must(t, testGotResult == testMustResult, "ct.TestResult should be %s, but got %s", testMustResult, testGotResult)
	testGotResult = ct1.TestResult.String()
	must(t, testGotResult == testMustResult, "ct1.TestResult should be %s, but got %s", testMustResult, testGotResult)
	testGotResult = ct2.TestResult.String()
	must(t, testGotResult == testMustResult, "ct2.TestResult should be %s, but got %s", testMustResult, testGotResult)
}

func must(t *testing.T, ret bool, errFormat string, args ...any) {
	if !ret {
		t.Errorf(errFormat, args...)
	}
}

func BenchmarkCore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pc := &object.Cycle{}
		pc.Compts = make([]object.Componenter, 0, 1001)
		pc.Compts = append(pc.Compts, &Core{})
		for i := 0; i < 1000; i++ {
			pc.Compts = append(pc.Compts, &CoreTestExample{})
		}
		pc.Compts[0].Start(pc)
	}
}

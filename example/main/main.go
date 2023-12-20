//  Copyright (C) 晓白齐齐,版权所有.

package main

import (
	"fmt"

	"github.com/bqqsrc/goper"
)

func main() {
	if err := goper.Launch(); err != nil {
		fmt.Printf("err is %v, %T", err, err)
	} else {
		fmt.Println("err is nil")
	}
}

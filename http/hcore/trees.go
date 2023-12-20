// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

//  Copyright (C) 晓白齐齐,版权所有. 2023

package hcore

import (
	"fmt"
	"strings"

	"github.com/bqqsrc/goper/http"
)

type completeTrees struct {
	domainTrees  MethodTrees
	generalTrees MethodTrees
}

func (d *completeTrees) addRouter(domain, path, method string, handler http.HttpHandler) error {
	var err error
	if path == "" {
		return fmt.Errorf("addRouter can't add a empty path")
	}
	if domain == "" {
		if d.generalTrees, err = addRouter(path, method, handler, d.generalTrees); err != nil {
			return err
		}
	} else {
		path = getWholePath(domain, path)
		if d.domainTrees, err = addRouter(path, method, handler, d.domainTrees); err != nil {
			return err
		}
	}
	return nil
}

func getWholePath(domain, path string) string {
	if domain == "" {
		return path
	}
	var newPath strings.Builder
	domainSlice := strings.Split(domain, ".")
	domainSepStr := string(domainSep)
	for index := len(domainSlice) - 1; index >= 0; index-- {
		newPath.WriteString(domainSepStr)
		newPath.WriteString(domainSlice[index])
		// if index != 0 {
		// 	newPath.WriteString(".")
		// }
	}
	// newPath.WriteString("@")
	newPath.WriteString(path)
	return newPath.String()
}

func (d *completeTrees) getHandler(domain, path, method string) (nodeValue, error) {
	if domain != "" && d.domainTrees != nil {
		domainPath := getWholePath(domain, path)
		if value, err := getHandler(domainPath, method, d.domainTrees); err == nil && value.handler != nil {
			return value, err
		}
	}
	if d.generalTrees != nil {
		return getHandler(path, method, d.generalTrees)
	}
	return nodeValue{}, fmt.Errorf("handler not found")
}

type MethodTree struct {
	method string
	root   *node
}

type MethodTrees = []MethodTree

func findMethod(method string, trees MethodTrees) (MethodTree, int) {
	for i, v := range trees {
		if v.method == method {
			return v, i
		}
	}
	return MethodTree{}, -1
}

func addRouter(path, method string, handler http.HttpHandler, trees MethodTrees) (MethodTrees, error) {
	if trees == nil {
		trees = make(MethodTrees, 0)
	}
	tree, index := findMethod(method, trees)
	if tree.root == nil {
		tree.root = &node{}
	}
	tree.root.addRouter(path, handler)
	if index < 0 {
		tree.method = method
		trees = append(trees, tree)
	} else {
		trees[index] = tree
	}
	return trees, nil
}

func getHandler(path, method string, trees MethodTrees) (nodeValue, error) {
	if tree, index := findMethod(method, trees); index >= 0 {
		return tree.root.getValue(path, nil, nil, false)
	}
	return nodeValue{}, fmt.Errorf("method not found")
}

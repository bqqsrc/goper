// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// at https://github.com/julienschmidt/httprouter/blob/master/LICENSE

package http

import (
	httpp "net/http"
	"net/url"

	"github.com/bqqsrc/bqqg/errors"
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

// Get returns the value of the first Param which key matches the given name and a boolean true.
// If no matching Param is found, an empty string is returned and a boolean false .
func (ps Params) Get(key string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return "", false
}

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByKey(key string) (va string) {
	va, _ = ps.Get(key)
	return
}

type Context struct {
	Request    *httpp.Request
	Writer     httpp.ResponseWriter
	ParamsData Params
	QueryData  url.Values
	FormData   url.Values
	Keys       map[string]any
	Response   Response
	Errors     errors.ErrorGroup
}

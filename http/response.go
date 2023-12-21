//  Copyright (C) 晓白齐齐,版权所有.

package http

import "encoding/json"

type Response interface {
	Bytes() []byte
	SetData(string, any)
	MoveData(string)
}

type responseData struct {
	Data map[string]any
}

func (r *responseData) SetData(key string, value any) {
	if r.Data == nil {
		r.Data = make(map[string]any)
	}
	r.Data[key] = value
}

func (r *responseData) MoveData(key string) {
	if r.Data != nil {
		delete(r.Data, key)
	}
}

func (r *responseData) Bytes() []byte {
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return []byte(err.Error())
	}
	return ret
}

func ResponseData() *responseData {
	return &responseData{}
}

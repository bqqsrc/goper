//  Copyright (C) 晓白齐齐,版权所有.

package http

type Response interface {
	Bytes() []byte
	SetData(string, any)
	MoveData(string)
}

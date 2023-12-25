//  Copyright (C) 晓白齐齐,版权所有.

package http

const (
	HttpLogic HttpPhase = iota
	HttpNotFound
	HttpError
	HttpFinish
	HttpNext
)

const HttpPhaseCount = 3

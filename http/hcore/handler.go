//  Copyright (C) 晓白齐齐,版权所有.

package hcore

import (
	"github.com/bqqsrc/bqqg/errors"
	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/log"
)

var phaseOrder = [...]http.HttpPhase{http.HttpLogic}

type HandlersChain []http.HttpHandler

type phaseHandlers struct {
	phase    http.HttpPhase
	handlers HandlersChain
}

type httpPhaseHandlers map[http.HttpPhase]phaseHandlers

func (h httpPhaseHandlers) Handler(context *http.Context) http.HttpPhase {
	if h == nil || len(h) == 0 {
		log.Errorln("httpPhaseHandlers is empty")
		return http.HttpError
	}
	var errGroup errors.ErrorGroup
	index := 0
walk:
	for index < len(phaseOrder) {
		phase := phaseOrder[index]
	next:
		if phaseHandlers, ok := h[phase]; ok {
			handlers := phaseHandlers.handlers
			if handlers != nil && len(handlers) > 0 {
				hIndex := 0
				for hIndex < len(handlers) {
					if handler := handlers[hIndex]; handler != nil {
						nextPhase := handler.Handler(context)
						switch nextPhase {
						case http.HttpNext:
							hIndex++
						case http.HttpFinish:
							break walk
						default:
							index = 0
							for index < len(phaseOrder) {
								if phaseOrder[index] == nextPhase {
									phase = nextPhase
									goto next
								}
							}
							phase = nextPhase
							if _, ok = h[phase]; !ok {
								errGroup = errGroup.AddErrorf("not found handlers of HttpPhase")
							}
							goto next
						}
					}
				}
			}
		}
		index++
	}
	if errGroup != nil {
		log.Errorf("Handler error: %s", errGroup)
		return http.HttpError
	}
	return http.HttpFinish
}

func (h httpPhaseHandlers) addHandlers(phase http.HttpPhase, handlers ...http.HttpHandler) httpPhaseHandlers {
	if handlers == nil || len(handlers) == 0 {
		return h
	}
	if h == nil {
		h = make(httpPhaseHandlers, http.HttpPhaseCount)
		h[phase] = phaseHandlers{phase, handlers}
	} else {
		if h[phase].handlers == nil {
			h[phase] = phaseHandlers{phase, handlers}
		} else {
			h[phase] = phaseHandlers{phase, append(h[phase].handlers, handlers...)}
		}
	}
	return h
}

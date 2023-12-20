//  Copyright (C) 晓白齐齐,版权所有.

package hcore

import (
	"github.com/bqqsrc/goper/http"
)

type rest struct {
	methods  []string
	handlers httpPhaseHandlers
}
type router struct {
	domain string
	path   string
	rests  []rest
}

type server struct {
	listen
	routers []router
}

func (s *server) addHander(routerIndex, restIndex int, phase http.HttpPhase, handlers ...http.HttpHandler) {
	if s.routers == nil {
		s.routers = make([]router, routerIndex+1)
	} else if len(s.routers) <= routerIndex {
		appRouters := make([]router, routerIndex+1-len(s.routers))
		s.routers = append(s.routers, appRouters...)
	}
	if s.routers[routerIndex].rests == nil {
		s.routers[routerIndex].rests = make([]rest, restIndex+1)
	} else if len(s.routers[routerIndex].rests) <= restIndex {
		appRests := make([]rest, restIndex+1-len(s.routers[routerIndex].rests))
		s.routers[routerIndex].rests = append(s.routers[routerIndex].rests, appRests...)
	}
	s.routers[routerIndex].rests[restIndex].handlers = s.routers[routerIndex].rests[restIndex].handlers.addHandlers(phase, handlers...)
}

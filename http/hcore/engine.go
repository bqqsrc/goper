//  Copyright (C) 晓白齐齐,版权所有.

package hcore

import (
	"fmt"
	httpp "net/http"
	"strings"
	"sync"

	"github.com/bqqsrc/goper/http"
)

type engine struct {
	listenInfo listen
	pool       sync.Pool
	trees      *completeTrees
}

func new() *engine {
	e := &engine{}
	e.pool.New = func() any { return &http.Context{} }
	e.trees = &completeTrees{}
	return e
}

func (e *engine) setListen(l listen) error {
	if err := e.checkListen(l); err != nil {
		return err
	}
	e.listenInfo = l
	return nil
}

func (e *engine) checkListen(l listen) error {
	if e.listenInfo.protocol.Name == "" {
		return nil
	}
	if e.listenInfo.protocol.Name != l.protocol.Name {
		return fmt.Errorf("listen should has only one protocol, but find different protocol for same port, port: %d, protocol: %s, %s ", e.listenInfo.port, e.listenInfo.protocol.Name, l.protocol.Name)
	}
	if e.listenInfo.protocol.CrtFile != l.protocol.CrtFile || e.listenInfo.protocol.KeyFile != l.protocol.KeyFile {
		return fmt.Errorf("listen should has only one crt-key pair, but find different crt-key pair for same port, port: %d, crt: %s, %s, key: %s, %s ",
			e.listenInfo.port, e.listenInfo.protocol.CrtFile, l.protocol.CrtFile, e.listenInfo.protocol.KeyFile, l.protocol.KeyFile)
	}
	return nil
}

func (e *engine) addRouter(domain, path, method string, handler http.HttpHandler) error {
	return e.trees.addRouter(domain, path, method, handler)
}

func (e *engine) getHandler(domain, path, method string) (http.HttpHandler, *http.Params) {
	if value, err := e.trees.getHandler(domain, path, method); err == nil && value.handler != nil {
		return value.handler, value.params
	}
	return nil, nil
}

func (e *engine) handleHTTPRequest(c *http.Context) {
	meth := c.Request.Method
	domain := strings.Split(c.Request.Host, ":")[0]
	path := c.Request.URL.Path
	if e.trees != nil {
		if handler, params := e.getHandler(domain, path, meth); handler != nil {
			if params != nil {
				c.ParamsData = *params
			}
			phase := handler.Handler(c)
			switch phase {
			case http.HttpError:
				//TODO handler call err
			case http.HttpFinish:
				//TODO handler call finish
			}
		} else {
			//TODO: handler nil, not found
		}
	} else {
		//TODO:trees empty

	}
}

func (e *engine) ServeHTTP(w httpp.ResponseWriter, req *httpp.Request) {
	c := e.pool.Get().(*http.Context)
	c.Writer = w
	c.Request = req

	e.handleHTTPRequest(c)

	e.pool.Put(c)
}

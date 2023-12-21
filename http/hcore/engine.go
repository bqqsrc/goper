//  Copyright (C) 晓白齐齐,版权所有.

package hcore

import (
	"errors"
	"fmt"
	"io/ioutil"
	httpp "net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/bqqsrc/goper/http"
	"github.com/bqqsrc/goper/log"
)

const defaultMultipartMemory = 32 << 20

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

func (e *engine) getHandler(domain, path, method string) (*httpPhaseHandlers, *http.Params) {
	if value, err := e.trees.getHandler(domain, path, method); err == nil && value.handler != nil {
		if handler, ok := value.handler.(*httpPhaseHandlers); ok {
			return handler, value.params
		} else {
			return nil, nil
		}
		// return value.handler, value.params
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
			c.QueryData = c.Request.URL.Query()
			c.FormData = make(url.Values)
			if err := c.Request.ParseMultipartForm(defaultMultipartMemory); err != nil {
				if !errors.Is(err, httpp.ErrNotMultipart) {
					log.Errorf("error on parse multipart form array: %v", err)
					return
				}
			}
			c.FormData = c.Request.PostForm
			c.Body, _ = ioutil.ReadAll(c.Request.Body)

			phase := handler.Handler(c)
			if c.NotResponse {
				return
			}
			switch phase {
			case http.HttpError:
				//TODO handler call err
				c.Writer.Write([]byte(c.Errors.Error()))
			case http.HttpFinish:
				//TODO handler call finish
				c.Writer.Write(c.Response.Bytes())
			}
		} else {
			//TODO: handler nil, not found
			c.Writer.Write([]byte("404 not found"))
		}
	} else {
		//TODO:trees empty
		c.Writer.Write([]byte("404 not found"))
	}
}

func (e *engine) ServeHTTP(w httpp.ResponseWriter, req *httpp.Request) {
	c := e.pool.Get().(*http.Context)
	c.Reset()
	c.Writer = w
	c.Request = req

	e.handleHTTPRequest(c)

	e.pool.Put(c)
}

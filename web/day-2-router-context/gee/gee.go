package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc 函数类型，input与output一致，便为同一类型
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	// map记录path与handler的映射关系
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	// URL
	key := method + "-" + pattern
	// handler 实现url与handler的映射
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

func (engine *Engine) Run(add string) (err error) {
	return http.ListenAndServe(add, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path

	if handler, ok := engine.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

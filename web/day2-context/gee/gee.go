package gee

import "net/http"

// HandlerFunc 函数类型，input与output一致，便为同一类型
type HandlerFunc func(c *Context)

type Engine struct {
	// map记录path与handler的映射关系
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	// URL
	engine.router.addRouter(method, pattern, handler)
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

	c := newContext(w, r)
	engine.router.handle(c)
}

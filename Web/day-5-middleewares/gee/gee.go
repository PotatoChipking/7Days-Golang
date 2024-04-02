package gee

import (
	"log"
	"net/http"
	"strings"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance共享一个Engine实例
}

type Engine struct {
	// 一个指针
	*RouterGroup

	router *router
	groups []*RouterGroup // store all groups
}

func New() *Engine {
	//1. 创建一个新的引擎对象。
	//2. 初始化引擎对象的路由器（router）为一个新的路由器对象。
	engine := &Engine{
		router: newRouter(),
	}
	//初始化引擎对象的路由组（RouterGroup）为一个新的路由组对象，其中包含了当前引擎对象的引用。
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	//将初始化的路由组对象添加到引擎对象的路由组数组中。
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// engine上层加了一层封装，便可使用同一个engine实例进行分组
// 每个group记录prefix，作为url的一部分，原本调用enginge传入pattern时，这里需要group内部的prefix+pattern

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	// 共享一个engine实例
	engine := group.engine

	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup

}

// HandlerFunc 函数类型，input与output一致，便为同一类型
type HandlerFunc func(c *Context)

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// middleware 新增method
func (group *RouterGroup) Use(middleware ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middleware...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {

		// 由于Use将某个中间件应用与某个group中
		// 针对请求的url，判断prefix是否与某个group相同
		// 这里将对应的group的中间件放入context中，用于handle
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

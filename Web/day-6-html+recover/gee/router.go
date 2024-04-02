package gee

import (
	"net/http"
	"strings"
)

// 路由模型，负责将url映射到handle，处理请求
type router struct {
	// trie 用于保存url，作为pattern，共同前缀为part存储在节点中
	// pattern存储于 iswild的节点中
	// 相对于map，可以节省key的存储空间
	roots map[string]*node
	// 将url映射到handler，key为 method + url
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 业务逻辑，将url写入至trie，router中
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 对于abc/deg/*/sdaf/ss  parts中存储的为["abc", "deg", "*"]
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// 将handle加入队列中
// 相对于之前，由router定位到hanfle，调用handle(c)，的串行执行
// 通过队列 +  index的方式可以允许乱序执行，通过C.next调用下一个handler
// 等待下一个handle执行完成后，返回执行当前handle的next后面的语句
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		// 对于一个动态路由/as/:name/pro
		// Params记录了map[name] = mike
		c.Params = params
		// 将当前handler加入队列中
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// 调用下一个
	c.Next()
}

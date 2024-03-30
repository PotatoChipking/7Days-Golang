package main

import (
	"day3-routerTrie/gee"
	"net/http"
)

// 动态路由，理解为：/*可以匹配到多个参数，比如：name，这里name为形参，实参可以是mike，rose。
// 静态路由，只能匹配到写好的/mike，/rose

// 这里调用过程：r.get方法将pattern与HandlerFunc传递至engine中，enging通过router.add方法

// 将pattern写入trie中，调用handle监听网页时，router对传入的pattern进行检索，调用对应的HandlerFunc
func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}

package main

import (
	"gee"
	"net/http"
)

// Engine的map成员记录url-handler的映射
// 相较于switch的方案，可以通过⭐外部调用⭐新增映射

// 对于handlerfun，将router与http的处理/返回分离
// 将http的res获取封装到context中，用于数据更细粒度的交互

// 将router拆分，只负责存储/调用/新增  URL与hander

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}

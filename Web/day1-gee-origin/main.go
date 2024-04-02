package main

import (
	"fmt"
	"gee"
	"net/http"
)

// Engine的map成员记录url-handler的映射
// 相较于switch的方案，可以通过⭐外部调用⭐新增映射

func main() {
	r := gee.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9999")
}

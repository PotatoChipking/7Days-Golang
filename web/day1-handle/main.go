package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND %s\n", r.URL)
	}
}

// 观察main函数，由url调用handle的部分集成至engine中
// engine封装了respnseWriter和 request
// 内部通过switch语句实现⭐url与⭐handler的映射

func main() {
	engine := new(Engine)

	log.Fatal(http.ListenAndServe("localhost:8080", engine))
}

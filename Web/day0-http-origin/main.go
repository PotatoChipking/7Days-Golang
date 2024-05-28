package main

import (
	"fmt"
	"log"
	"net/http"
)

// 对于每个url，调用一个handle处理
func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.Path)
}

func counter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "count", r.URL.Path)
}

//func main() {
//	fmt.Println("Go版本为：", runtime.Version())
//}

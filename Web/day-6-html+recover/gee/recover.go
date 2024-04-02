package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}

func trace(message string) string {
	var pcs [32]uintptr

	// 返回调用栈的程序计数器，[0]: Callers本身，[1]:trace， [2]: defer func
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback")
	for _, pc := range pcs[:n] {

		// 通过 runtime.FuncForPC(pc) 获取对应的函数，在通过 fn.FileLine(pc) 获取到调用该函数的文件名和行号，打印在日志中。
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()

}

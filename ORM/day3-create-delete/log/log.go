package log

import (
	"io"
	"log"
	"os"
	"sync"
)

// 1、日志分级
// 2、颜色区分
// 3、显示打印日志的文件名和行号
var (
	// 注册不同的log， 第三个参数支持显示文件名和行号
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)

	loggers = []*log.Logger{errorLog, infoLog}
	mu      sync.Mutex
)

// log methods
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	// 迭代类型
	InfoLevel = iota
	ErrorLevel
	Disable
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	// 这里<号，需要先判断较大的值，>反过来
	if ErrorLevel < level {
		errorLog.SetOutput(io.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}

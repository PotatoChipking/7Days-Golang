package codec

import "io"

// 请求头，用于定义方法名称、消息num
type Header struct {
	ServiceMethod string
	Seq           uint64
	Error         string
}

// 抽象出编解码器的接口，便于在调用接口返回值不变的情况下，更换实例-json/binary
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// Codec construct

// NewCodecFunc 类似HandleFunc
type NewCodecFunc func(closer io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

// Type与Func的映射,根据类型调用对应的函数类型实现
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}

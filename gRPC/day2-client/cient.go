package day2_client

import (
	d1 "7days-golang-learn/gRPC/day1-decode"
	"7days-golang-learn/gRPC/day1-decode/codec"
	"errors"
	"io"
	"sync"
)

// 远程调用函数要求：
// 1、函数类型可以外部访问
// 2、函数可以外部访问
// 3、函数有两个参数，都是可以外部访问的
// 4、函数第二个参数是指针
// 5、函数返回值为error
// 例：
// func (t *T) MethodName(argType T1, replyType *T2) error

// 承载一次RPC调用所需信息
type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call
}

func (call *Call) done() {
	call.Done <- call
}

type Client struct {
	cc      codec.Codec
	opt     *d1.Option
	sending sync.Mutex
	header  codec.Header
	mu      sync.Mutex
	seq     uint64

	pending  map[uint64]*Call
	closing  bool
	shutdown bool
}

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("connection is shutd")

func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

func (clinet *Client) IsAvaliable() bool {
	clinet.mu.Unlock()
	defer clinet.mu.Unlock()
	return !clinet.shutdown && !clinet.closing
}

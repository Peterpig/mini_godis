package protocol

import (
	"bytes"
	"fmt"
	"strconv"
)

// redis协议
// https://www.redis.com.cn/topics/protocol.html

var (
	CRLF = "\r\n"
)

// 我们简单的将 Reply 分为两类:
// 单行: StatusReply, IntReply, ErrorReply
// 多行: BulkReply, MultiBulkReply

type StatusReplay struct {
	Status string
}

func MakeStatusReplay(status string) *StatusReplay {
	return &StatusReplay{Status: status}
}

func (r *StatusReplay) ToBytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

type StandardErrReply struct {
	Status string
}

func MakeErrorReplay(status string) *StandardErrReply {
	return &StandardErrReply{Status: status}
}

func (r *StandardErrReply) ToBytes() []byte {
	return []byte("-" + r.Status + CRLF)
}

type IntReply struct {
	Code int64
}

func MakeIntReplay(code int64) *IntReply {
	return &IntReply{Code: code}
}

func (r *IntReply) ToBytes() []byte {
	return []byte(":" + fmt.Sprintf("%d", r.Code) + CRLF)
}

// 字符串数组
type MultiBulkReply struct {
	Args [][]byte
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: args}
}

func (r *MultiBulkReply) ToBytes() []byte {
	argsLen := len(r.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argsLen) + CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString("$-1" + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"runtime/debug"
	"strconv"

	"github.com/Peterpig/mini_godis/interface/redis"
	"github.com/Peterpig/mini_godis/lib/logger"
	"github.com/Peterpig/mini_godis/redis/protocol"
)

type Playload struct {
	Data redis.Reply
	Err  error
}

func ParseStream(reader io.Reader) <-chan *Playload {
	ch := make(chan *Playload)
	go parse(reader, ch)
	return ch
}

func parse(rawReader io.Reader, ch chan<- *Playload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("%v, %v", err, debug.Stack())
		}
	}()

	reader := bufio.NewReader(rawReader)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			ch <- &Playload{Err: err}
			close(ch)
			return
		}
		logger.Info("read data, %s", string(line))

		length := len(line)

		if length <= 2 || line[length-2] != '\r' {
			continue
		}

		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		switch line[0] {
		case '+':
			content := string(line[1:])
			ch <- &Playload{Data: protocol.MakeStatusReplay(content)}
		case '-':
			content := string(line[1:])
			ch <- &Playload{Data: protocol.MakeErrorReplay(content)}
		case ':':
			value, err := strconv.ParseInt(string(line[1:]), 10, 64)
			if err != nil {
				protocolError(ch, "invalid integer"+string(line[1:]))
			}
			ch <- &Playload{Data: protocol.MakeIntReplay(value)}
		case '*':
			err := parseArray(line, reader, ch)
			if err != nil {
				protocolError(ch, "invalid array"+string(line))
				close(ch)
				return
			}
		}
	}
}

func protocolError(ch chan<- *Playload, msg string) {
	err := errors.New("protocol error: " + msg)
	ch <- &Playload{Err: err}
}

func parseArray(line []byte, reader *bufio.Reader, ch chan<- *Playload) error {
	nStrs, err := strconv.ParseInt(string(line[1:]), 10, 64)
	if err != nil {
		protocolError(ch, "invalid array length"+string(line[1:]))
		return nil
	} else if nStrs == 0 {
		ch <- &Playload{Data: protocol.MakeEmptyMultiBulkReply()}
		return nil
	}
	lines := make([][]byte, 0, nStrs)
	for i := int64(0); i < nStrs; i++ {
		var line []byte
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		length := len(line)

		if length < 4 || line[length-2] != '\r' || line[0] != '$' {
			protocolError(ch, "invalid bulk string "+string(line))
			break
		}

		strLen, err := strconv.ParseInt(string(line[1:length-2]), 10, 64)
		if err != nil || strLen < -1 {
			protocolError(ch, "invalid bulk string length "+string(line))
			break
		} else if strLen == -1 {
			lines = append(lines, []byte{})
		} else {
			body := make([]byte, strLen)
			_, err := io.ReadFull(reader, body)
			if err != nil {
				return err
			}
			lines = append(lines, body[:len(body)-2])
		}
	}

	ch <- &Playload{Data: protocol.MakeMultiBulkReply(lines)}
	return nil
}

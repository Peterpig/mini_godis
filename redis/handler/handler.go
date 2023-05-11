package handler

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/Peterpig/mini_godis/interface/db"
	"github.com/Peterpig/mini_godis/lib/logger"
	"github.com/Peterpig/mini_godis/redis/parser"
)

type Handler struct {
	activeConn sync.Map
	db         db.DB
	closing    bool
}

func MakeHandler() *Handler {
	return &Handler{}
}

// func (h *Handler) closeClieng()

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing {
		conn.Close()
	}

	ch := parser.ParseStream(conn)
	for playload := range ch {
		if playload.Err == io.EOF || playload.Err == io.ErrUnexpectedEOF || strings.Contains(playload.Err.Error(), "use of closed network connection") {
			logger.Info("Connection closed: ")
		}
		fmt.Println("playload = ", playload)
	}
}

func (h *Handler) Close() error {
	return nil
}

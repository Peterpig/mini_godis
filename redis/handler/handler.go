package handler

import (
	"context"
	"net"
	"sync"

	"github.com/Peterpig/mini_godis/interface/db"
)

type Handler struct {
	activeConn sync.Map
	db         db.DB
	closing    bool
}

func MakeHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {}

func (h *Handler) Close() error {
	return nil
}

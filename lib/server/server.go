package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Peterpig/mini_godis/interface/tcp"
	"github.com/Peterpig/mini_godis/lib/logger"
)

type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max-connect"`
	Timeout    time.Duration `yaml:"timeout"`
}

func LinstenAndServerWithSignal(cfg *Config, handler tcp.Handle) (err error) {
	signCh := make(chan os.Signal, 1)
	errCh := make(chan error)

	var linster net.Listener

	signal.Notify(signCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case <-signCh:
			logger.Info("Get exit signal")
		case err := <-errCh:
			logger.Error("Accept error: %v", err)
		}
		logger.Info("Shutting down ....")
		linster.Close()
		handler.Close()
		os.Exit(0)
	}()

	linster, err = net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Fatal(fmt.Sprintf("listen err: %v", err))
		return
	}

	logger.Info("bind: %s, start listening...", cfg.Address)
	defer linster.Close()
	defer handler.Close()
	ctx := context.Background()
	var wait sync.WaitGroup

	connCh := make(chan struct{}, cfg.MaxConnect)

	for {
		conn, err := linster.Accept()
		if err != nil {
			errCh <- err
			break
		}
		connCh <- struct{}{}
		wait.Add(1)
		logger.Info("Accept link")

		go func() {
			defer func() {
				<-connCh
				wait.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}

	wait.Wait()
	return
}

package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"test/faraway/client/config"
	"test/faraway/client/internal/controller/handler"
	"test/faraway/client/internal/pow"
	"test/faraway/client/pkg/logger"
	"test/faraway/client/pkg/protocol"
	"test/faraway/client/pkg/tcpclient"
	"time"
)

func Run(cfg *config.Config) {

	// set log level
	logger.SetLevel(cfg.Level)

	log := logger.Get()

	protocol.SetMaxReadSize(cfg.MaxReadSize)

	// proof of work for client
	pw := pow.New(pow.MaxAttemptsCount(cfg.MaxAttemptsCount))

	// new handler
	h := handler.New(pw)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// new TCP client that calls handler if connection successful
	cli := tcpclient.New(
		ctx,
		h,
		tcpclient.Port(cfg.Port),
		tcpclient.DialTimeout(time.Duration(int64(cfg.DialTimeoutSecond)*int64(time.Second))),
		tcpclient.ConnTimeout(time.Duration(int64(cfg.ConnTimeoutSecond)*int64(time.Second))),
		tcpclient.MaxDailCount(cfg.MaxDailCount),
		tcpclient.DelayDail(time.Duration(int64(cfg.DelayDailMillisecond)*int64(time.Millisecond))),
		tcpclient.ReadTimeout(time.Duration(int64(cfg.ReadTimeoutSecond)*int64(time.Second))),
	)

	// for interrupting from keyboard
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app run tcp client signal: " + s.String())
	case err := <-cli.Error():
		log.Error("app run tcp client error: %s", err)
	case <-cli.Done():
		log.Info("app run tcp client done: bye")
	}
}

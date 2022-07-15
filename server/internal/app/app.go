package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"test/faraway/server/config"
	"test/faraway/server/internal/controller/handler"
	"test/faraway/server/internal/pow"
	"test/faraway/server/internal/usecase"
	"test/faraway/server/pkg/logger"
	"test/faraway/server/pkg/protocol"
	"test/faraway/server/pkg/tcpserver"
	"time"
)

func Run(cfg *config.Config) {

	// set log level
	logger.SetLevel(cfg.Level)

	log := logger.Get()

	protocol.SetMaxReadSize(cfg.MaxReadSize)

	// new word of wisdom use case
	uc := usecase.NewWordOfWisdom()

	// new proof of work
	pw := pow.New(pow.LeadingZeroCount(cfg.LeadingZeroCount))

	// new handler
	h := handler.New(uc, pw)

	// context to stop tcp server on keyboard interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// new TCP server that calls a handler for each connection
	srv := tcpserver.New(
		ctx,
		h,
		tcpserver.Port(cfg.Port),
		tcpserver.AcceptDelay(time.Duration(int64(cfg.AcceptDelayMillisecond)*int64(time.Millisecond))),
		tcpserver.ConnTimeout(time.Duration(int64(cfg.ConnTimeoutSecond)*int64(time.Second))),
		tcpserver.ReadTimeout(time.Duration(int64(cfg.ReadTimeoutSecond)*int64(time.Second))),
	)

	// for interrupting from keyboard
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app run tcp server signal: " + s.String())
	case err := <-srv.Error():
		log.Error("app run tcp server error: %s", err)
	}

}

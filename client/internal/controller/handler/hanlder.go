package handler

import (
	"context"
	"io"
	"test/faraway/client/pkg/logger"
	"test/faraway/client/pkg/protocol"
)

type powI interface {
	Solve(context.Context, io.ReadWriter) error
}

type Handler struct {
	pow powI
}

func New(pow powI) Handler {
	h := Handler{
		pow: pow,
	}
	return h
}

func (h Handler) Handle(ctx context.Context, rw io.ReadWriter) (err error) {

	// get logger
	log := logger.Get()

	// solve a challenge
	err = h.pow.Solve(ctx, rw)
	if err != nil {
		log.Error("handler pow.Solve error: %s", err)
		return
	}

	log.Info("challenge solve success")

	// wrap io.ReadWriter for convenience
	byProtocol := protocol.Get()

	var wordOfWisdom []byte
	// read word of wisdom
	wordOfWisdom, err = byProtocol.Read(rw)
	if err != nil {
		log.Error("handler read err: %s", err)
		return
	}

	log.Info("read success: %s", string(wordOfWisdom))
	return
}

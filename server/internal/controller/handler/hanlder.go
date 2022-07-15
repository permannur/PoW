package handler

import (
	"context"
	"io"
	"test/faraway/server/pkg/logger"
	"test/faraway/server/pkg/protocol"
)

type uc interface {
	GetWordOfWisdom() string
}

type powI interface {
	Verify(ctx context.Context, writer io.ReadWriter) error
}

type Handler struct {
	useCase uc
	pow     powI
}

func New(u uc, pow powI) Handler {
	h := Handler{
		useCase: u,
		pow:     pow,
	}
	return h
}

func (h Handler) Handle(ctx context.Context, rw io.ReadWriter) {

	// get logger
	log := logger.Get()

	// verify proof of work
	err := h.pow.Verify(ctx, rw)
	if err != nil {
		log.Error("error verifying pow, err: %s", err)
		return
	}

	log.Info("pow verifying success")

	// wrap io.ReadWriter for convenience
	byProtocol := protocol.Get()

	// get word of wisdom
	wordOfWisdom := h.useCase.GetWordOfWisdom()

	// write word of wisdom
	err = byProtocol.Write(rw, []byte(wordOfWisdom))
	if err != nil {
		log.Error("error write to conn: %s", err)
		return
	}

	log.Info("send word of wisdom")
}

package pow

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"io"
	"test/faraway/server/pkg/logger"
	"test/faraway/server/pkg/protocol"
)

const (
	_defaultLeadingZeroCount = 20
)

var _defaultGenerator = uuid.New

type Pow struct {
	leadingZeroCount byte
	generator        func() []byte
}

// New proof of work for server
func New(opts ...Option) Pow {
	p := Pow{
		leadingZeroCount: _defaultLeadingZeroCount,
		generator: func() []byte {
			bt16 := _defaultGenerator()
			return bt16[:]
		},
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// Verify proof of work by interacting with the client
func (p Pow) Verify(ctx context.Context, rw io.ReadWriter) (err error) {

	log := logger.Get()

	// wrap io.ReadWriter for convenience
	byProtocol := protocol.Get()

	log.Info("leading zero count: %d", p.leadingZeroCount)

	// write a leading zero count
	err = byProtocol.Write(rw, []byte{p.leadingZeroCount})
	if err != nil {
		log.Error("error in writing leading zero nonce: %s", err)
		return
	}

	// generate 16 byte uuid as a challenge
	id := p.generator()

	challenge := id[:]

	// write a challenge
	err = byProtocol.Write(rw, challenge)
	if err != nil {
		log.Error("error in writing challenge: %s", err)
		return
	}

	log.Info("challenge sent")

	// read nonce in goroutine
	var nonce []byte
	ch := make(chan struct{})
	go func() {
		nonce, err = byProtocol.Read(rw)
		ch <- struct{}{}
	}()

	// if the context is closed, there is no need to wait
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
	}

	if err != nil {
		log.Error("error in reading nonce: %s", err)
		return
	}

	// check nonce length
	if len(nonce) != 4 {
		return fmt.Errorf("format error nonce, %d != 4", len(nonce))
	}

	log.Info("nonce received, verifying")

	// append nonce bytes to challenge
	challenge = append(challenge, nonce...)

	leadingZeroCount := getLeadingZeroCountFromHash(challenge)
	// check leading zero count
	if leadingZeroCount < p.leadingZeroCount {
		return fmt.Errorf("leading zero nonce is less than needed, %d < %d", leadingZeroCount, p.leadingZeroCount)
	}

	// calc nonce for log
	var k, nonceUint32 uint32 = 1, 0
	for i := 3; i >= 0; i-- {
		nonceUint32 += k * uint32(nonce[i])
		k *= 256
	}

	log.Info("ok nonce: %d", nonceUint32)

	return
}

// Hashes the content and returns a leading zero count
func getLeadingZeroCountFromHash(content []byte) (zeroCount byte) {
	var hash [32]byte

	hash = sha256.Sum256(content)
loopI:
	for i := 0; i < 32; i++ {
		for j := 7; j >= 0; j-- {
			if (hash[i]>>j)&1 == 1 {
				break loopI
			}
			zeroCount++
		}
	}
	return
}

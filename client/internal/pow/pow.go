package pow

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"math"
	"test/faraway/client/pkg/logger"
	"test/faraway/client/pkg/protocol"
)

const (
	_defaultMaxAttemptsCount = math.MaxUint32
)

type Pow struct {
	maxAttemptsCount uint32
}

// New proof of work for client
func New(opts ...Option) Pow {
	p := Pow{
		maxAttemptsCount: _defaultMaxAttemptsCount,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// Solve a challenge
func (p Pow) Solve(ctx context.Context, rw io.ReadWriter) (err error) {

	log := logger.Get()

	log.Info("Start solving a challenge")

	// create new protocol from io.ReadWriter
	byProtocol := protocol.Get()

	var btLeadingZeroCount []byte

	// read leading zero count
	btLeadingZeroCount, err = byProtocol.Read(rw)
	if err != nil {
		log.Error("error in reading leading zero count: %s", err)
		return
	}

	// check leading zero count length
	if len(btLeadingZeroCount) != 1 {
		return fmt.Errorf("format error in leading zero count, %d != 1", len(btLeadingZeroCount))
	}

	leadingZeroCount := btLeadingZeroCount[0]

	log.Info("leading zero count: %d", leadingZeroCount)

	var challenge []byte

	// read a challenge
	challenge, err = byProtocol.Read(rw)
	if err != nil {
		log.Error("error in reading challenge: %s", err)
		return
	}

	// check challenge length
	if len(challenge) != 16 {
		return fmt.Errorf("format error in challenge, %d != 16", len(challenge))
	}

	log.Info("challenge received")

	var nonce []byte
	var nonceUint32 uint32
	// bruteForce the challenge
	nonce, nonceUint32, err = p.bruteForce(ctx, leadingZeroCount, challenge)
	if err != nil {
		return
	}

	log.Info("challenge successfully solved, nonce: %d", nonceUint32)

	// write a nonce
	err = byProtocol.Write(rw, nonce)
	if err != nil {
		log.Error("error in writing nonce: %s", err)
		return
	}

	return
}

// bruteForce - brute force for solving challenge
func (p Pow) bruteForce(ctx context.Context, leadingZeroCount byte, challenge []byte) ([]byte, uint32, error) {

	// append nonce bytes to challenge
	challenge = append(challenge, 0, 0, 0, 0)

	challengeLen := len(challenge)

	var attemptsCount uint32

	for attemptsCount < p.maxAttemptsCount {
		select {
		case <-ctx.Done():
			return nil, attemptsCount, ctx.Err()
		default:
			// get leading zero count
			var zeroCount byte
			zeroCount = getLeadingZeroCountFromHash(challenge)

			if zeroCount < leadingZeroCount {
				// increment attempts count
				attemptsCount++
				k := attemptsCount
				// map attempts count to challenge bytes
				for i := 1; i <= 4; i++ {
					challenge[challengeLen-i] = byte(k % 256)
					k /= 256
				}
			} else {
				// return attempts count as 4 bytes
				return challenge[challengeLen-4:], attemptsCount, nil
			}
		}
	}
	return nil, attemptsCount, fmt.Errorf("could not find a matching hash")
}

// Hashes the content and returns a leading zero count
func getLeadingZeroCountFromHash(challenge []byte) (zeroCount byte) {
	var hash [32]byte
	hash = sha256.Sum256(challenge)
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

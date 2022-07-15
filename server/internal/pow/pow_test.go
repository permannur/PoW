package pow

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"test/faraway/server/pkg/protocol"
	"testing"
	"time"
)

func TestPow_Solve(t *testing.T) {
	testCases := []struct {
		name             string
		leadingZeroCount byte
		uuid             string
		nonce            uint32
		success          bool
	}{
		{"test case 1", 20, "37cff247-b772-404f-9daa-2199ff622b8c", 762604, true},
		{"test case 2", 20, "37cff247-b772-404f-9daa-2199ff622b8c", 762605, false},
		{"test case 3", 15, "bd0f3b97-1b53-4dfa-ad5c-e4ecc7c9ee44", 9184, true},
		{"test case 4", 15, "bd0f3b97-1b53-4dfa-ad5c-e4ecc7c9ee44", 9185, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, client := net.Pipe()
			pow := New(
				LeadingZeroCount(tc.leadingZeroCount),
				Generator(
					func() []byte {
						id, _ := uuid.Parse(tc.uuid)
						return id[:]
					},
				),
			)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				assert.Equal(t, pow.Verify(ctx, server) == nil, tc.success)
				wg.Done()
			}()

			byProtocol := protocol.Get()

			btLeadingZeroCount, err := byProtocol.Read(client)
			assert.NoError(t, err)

			assert.Equal(t, len(btLeadingZeroCount), 1)

			var challenge []byte

			challenge, err = byProtocol.Read(client)
			assert.NoError(t, err)

			assert.Equal(t, len(challenge), 16)

			nonce := make([]byte, 4)
			k := tc.nonce
			for i := 3; i >= 0; i-- {
				nonce[i] = byte(k % 256)
				k /= 256
			}

			assert.NoError(t, byProtocol.Write(client, nonce))
			wg.Wait()
		})
	}
}

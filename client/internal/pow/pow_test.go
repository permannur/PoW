package pow

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"test/faraway/client/pkg/protocol"
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
		{"test case 1", 20, "37cff247-b772-404f-9daa-2199ff622b8c", 762605, false},
		{"test case 2", 15, "bd0f3b97-1b53-4dfa-ad5c-e4ecc7c9ee44", 9184, true},
		{"test case 3", 5, "c5990de7-7080-449d-832c-0aac304e237d", 21, true},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc1 := tc
			server, client := net.Pipe()

			pow := New()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func(ctx context.Context) {
				assert.NoError(t, pow.Solve(ctx, client))
				wg.Done()
			}(ctx)

			byProtocol := protocol.Get()

			assert.NoError(t, byProtocol.Write(server, []byte{tc1.leadingZeroCount}))

			id, err := uuid.Parse(tc1.uuid)
			assert.NoError(t, err)

			challenge := id[:]

			err = byProtocol.Write(server, challenge)

			var nonce []byte
			nonce, err = byProtocol.Read(server)

			if len(nonce) != 4 {
				assert.Error(t, err)
			}

			assert.EqualValues(t, 4, len(nonce))

			challenge = append(challenge, nonce...)

			assert.LessOrEqual(t, tc1.leadingZeroCount, getLeadingZeroCountFromHash(challenge))

			var k, nonceUint32 uint32 = 1, 0
			for i := 3; i >= 0; i-- {
				nonceUint32 += k * uint32(nonce[i])
				k *= 256
			}
			assert.Equal(t, tc1.nonce == nonceUint32, tc.success)
			wg.Wait()
			t.Deadline()
		})
	}
}

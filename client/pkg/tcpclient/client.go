package tcpclient

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	_defaultAddr         = ":6677"
	_defaultDialTimeout  = 3 * time.Second
	_defaultConnTimeout  = 5 * time.Second
	_defaultMaxDailCount = 10
	_defaultDelayDail    = 10 * time.Millisecond
	_defaultReadTimeout  = 3 * time.Second
)

// handler interface
type handler interface {
	Handle(context.Context, io.ReadWriter) error
}

// Client - tcp client
type Client struct {
	conn         net.Conn
	handler      handler
	addr         string
	dialTimeout  time.Duration
	connTimeout  time.Duration
	maxDailCount byte
	delayDail    time.Duration
	readTimeout  time.Duration
	error        chan error
	done         chan struct{}
	stopAtomic   int32
}

func New(ctx context.Context, handler handler, opts ...Option) *Client {
	c := &Client{
		handler:      handler,
		addr:         _defaultAddr,
		dialTimeout:  _defaultDialTimeout,
		connTimeout:  _defaultConnTimeout,
		maxDailCount: _defaultMaxDailCount,
		delayDail:    _defaultDelayDail,
		readTimeout:  _defaultReadTimeout,
		error:        make(chan error, 1),
		done:         make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(c)
	}

	go c.start(ctx)

	return c
}

// start tcp client that tries to connect to server
func (c *Client) start(ctx context.Context) {

	var err error

	defer func() {
		if err != nil {
			c.error <- err
		} else {
			c.done <- struct{}{}
		}
		close(c.error)
		close(c.done)
	}()

	c.conn, err = c.dail(ctx)
	if err != nil {
		return
	}
	defer func() {
		_ = c.conn.Close()
	}()

	ctxChild, cancel := context.WithTimeout(ctx, c.connTimeout)
	defer cancel()

	err = c.handler.Handle(ctxChild, c.conn)
}

func (c *Client) dail(ctx context.Context) (net.Conn, error) {

	for i := 0; i < int(c.maxDailCount); i++ {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			return nil, err
		default:
			// dail with time out
			conn, err := net.DialTimeout("tcp", c.addr, c.dialTimeout)
			if err == nil {
				// if connection success, then set read deadline
				err = conn.SetReadDeadline(time.Now().Add(c.readTimeout))
				if err == nil {
					return conn, nil
				}
			}
			// delay between unsuccessful attempts
			time.Sleep(c.delayDail)
		}
	}

	err := fmt.Errorf("reached max attempt count")
	return nil, err
}

// Error - for access error channel from outside
func (c Client) Error() <-chan error {
	return c.error
}

// Done - for access done channel from outside
func (c Client) Done() <-chan struct{} {
	return c.done
}

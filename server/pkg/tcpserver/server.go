package tcpserver

import (
	"context"
	"io"
	"net"
	"time"
)

const (
	_defaultAddr        = ":6677"
	_defaultAcceptDelay = 10 * time.Millisecond
	_defaultConnTimeout = 5 * time.Second
	_defaultReadTimeout = 3 * time.Second
)

// handler interface
type handler interface {
	Handle(ctx context.Context, writer io.ReadWriter)
}

// Server - tcp server
type Server struct {
	listener    net.Listener
	handler     handler
	addr        string
	acceptDelay time.Duration
	connTimeout time.Duration
	readTimeout time.Duration
	error       chan error
}

func New(ctx context.Context, handler handler, opts ...Option) *Server {
	s := &Server{
		handler:     handler,
		addr:        _defaultAddr,
		acceptDelay: _defaultAcceptDelay,
		connTimeout: _defaultConnTimeout,
		readTimeout: _defaultReadTimeout,
		error:       make(chan error, 1),
	}

	for _, opt := range opts {
		opt(s)
	}

	go s.start(ctx)

	return s
}

// starts listener and call handler for each accepted connection
func (s *Server) start(ctx context.Context) {

	var err error

	// start a listener with context
	lsCfg := net.ListenConfig{KeepAlive: -1}
	s.listener, err = lsCfg.Listen(ctx, "tcp", s.addr)

	// if error acquired, then notify app and return
	if err != nil {
		s.error <- err
		close(s.error)
		return
	}

	// accept connection until done not closed
	for {
		var conn net.Conn
		conn, err = s.listener.Accept()
		if err != nil {
			// delay between unsuccessful accept
			time.Sleep(s.acceptDelay)
			continue
		}

		// if accept successful, handle it before connection timeout
		go func(conn net.Conn) {
			defer conn.Close()

			// set read timeout
			err1 := conn.SetReadDeadline(time.Now().Add(s.readTimeout))
			if err1 != nil {
				return
			}

			// create a context with timeout
			ctxChild, cancel := context.WithTimeout(ctx, s.connTimeout)
			defer cancel()

			// call handler with context
			s.handler.Handle(ctxChild, conn)
		}(conn)
	}
}

// Error - for access notify channel from outside
func (s *Server) Error() <-chan error {
	return s.error
}

package tcpserver

import (
	"net"
	"time"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.addr = net.JoinHostPort("", port)
	}
}

func AcceptDelay(delay time.Duration) Option {
	return func(s *Server) {
		s.acceptDelay = delay
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.connTimeout = timeout
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

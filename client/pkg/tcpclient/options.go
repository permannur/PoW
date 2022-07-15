package tcpclient

import (
	"net"
	"time"
)

type Option func(client *Client)

func Port(port string) Option {
	return func(c *Client) {
		c.addr = net.JoinHostPort("", port)
	}
}

func DialTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.dialTimeout = timeout
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.connTimeout = timeout
	}
}

func MaxDailCount(count byte) Option {
	return func(c *Client) {
		c.maxDailCount = count
	}
}

func DelayDail(delay time.Duration) Option {
	return func(c *Client) {
		c.delayDail = delay
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.readTimeout = timeout
	}
}

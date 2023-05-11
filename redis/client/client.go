package client

import (
	"net"
	"sync"
)

type Client struct {
	conn net.Conn

	waitingReplay sync.WaitGroup

	mu sync.Mutex
}

func MakeClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) Close() error {
	c.conn.Close()
	return nil
}

func (c *Client) Write(b []byte) (err error) {
	if len(b) == 0 {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	_, err = c.conn.Write(b)
	return
}

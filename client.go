package funclib

import "net"

type (
	Client struct {
		conn net.Conn
	}
)

// TODO something more flexible with error handling will be need to compliment WithConn

// Standard dial with standard params like net.Dial.
func Dial(network, addr string) (*Client, error) {
	conn, err := net.Dial(network, addr)

	if err != nil {
		return nil, err
	}

	return &Client{conn}, nil
}

// Lend the created connection and do something.
func (c *Client) WithConn(handler func(conn net.Conn)) {
	handler(c.conn)
}

// Close the underlying net.Conn.
func (c *Client) Close() error {
	return c.conn.Close()
}

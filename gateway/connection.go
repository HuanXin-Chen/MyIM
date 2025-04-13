package gateway

import "net"

// 连接对象
type connection struct {
	fd   int
	conn *net.TCPConn
}

func (c *connection) Close() {
	err := c.conn.Close()
	panic(err)
}

func (c *connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

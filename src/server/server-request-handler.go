package server

import (
	"bufio"
	"net"
	"strconv"
)

// ServerRequestHandler ...
type ServerRequestHandler struct {
	port        int
	listener    net.Listener
	connection  net.Conn
	outToClient *bufio.Reader
	inToClient  *bufio.Writer
	remoteAddr  string
}

func newServerRequestHandler(port int) *ServerRequestHandler {
	tcpSRH := new(ServerRequestHandler)
	tcpSRH.listener, _ = net.Listen("tcp", ":"+strconv.Itoa(port))

	return tcpSRH
}

func (c *ServerRequestHandler) accept() {
	// log.Println("Listen on", tcpSRH.listener.Addr().String())
	c.connection, _ = c.listener.Accept()
	// log.Println("Accept a connection request from", conn.RemoteAddr())
	c.remoteAddr = c.connection.RemoteAddr().String()
	c.inToClient = bufio.NewWriter(c.connection)
	c.outToClient = bufio.NewReader(c.connection)
}

func (c *ServerRequestHandler) send(msg []byte) {
	c.inToClient.Write(msg)
	c.inToClient.Flush()
}

func (c *ServerRequestHandler) receive() ([]byte, error) {
	data := make([]byte, 4096)
	n, err := c.outToClient.Read(data)
	return data[:n], err
}

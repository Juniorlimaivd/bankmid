package server

import (
	"bufio"
	"encoding/gob"
	"net"
	"strconv"
)

// ServerRequestHandler ...
type ServerRequestHandler struct {
	port        int
	listener    net.Listener
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
	conn, _ := c.listener.Accept()
	// log.Println("Accept a connection request from", conn.RemoteAddr())
	c.remoteAddr = conn.RemoteAddr().String()
	c.inToClient = bufio.NewWriter(conn)
	c.outToClient = bufio.NewReader(conn)
}

func (c *ServerRequestHandler) send(msg []byte) {
	encoder := gob.NewEncoder(c.inToClient)

	encoder.Encode(msg)

	c.inToClient.Flush()
}

func (c *ServerRequestHandler) receive() []byte {
	decoder := gob.NewDecoder(c.outToClient)

	var data []byte

	decoder.Decode(&data)

	return data
}

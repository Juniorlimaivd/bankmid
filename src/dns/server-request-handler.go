package dns

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"strconv"
)

// ServerRequestHandler ...
type ServerRequestHandler struct {
	port        int
	listener    net.Listener
	conn        net.Conn
	outToClient *bufio.Reader
	inToClient  *bufio.Writer
	remoteAddr  string
}

func newServerRequestHandler(port int) *ServerRequestHandler {
	tcpSRH := new(ServerRequestHandler)
	tcpSRH.listener, _ = net.Listen("tcp", ":"+strconv.Itoa(port))

	log.Println("DNS Listen on", tcpSRH.listener.Addr().String())
	tcpSRH.conn, _ = tcpSRH.listener.Accept()

	log.Println("Accept a connection request from", tcpSRH.conn.RemoteAddr())
	tcpSRH.remoteAddr = tcpSRH.conn.RemoteAddr().String()
	tcpSRH.inToClient = bufio.NewWriter(tcpSRH.conn)
	tcpSRH.outToClient = bufio.NewReader(tcpSRH.conn)

	return tcpSRH
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

func (c *ServerRequestHandler) close() {
	c.conn.Close()
	c.listener.Close()
}

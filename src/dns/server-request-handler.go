package dns

import (
	"bufio"
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
	c.inToClient.Write(msg)
	c.inToClient.Flush()
}

func (c *ServerRequestHandler) receive() []byte {
	data := make([]byte, 4096)
	n, _ := c.outToClient.Read(data)
	return data[:n]
}

func (c *ServerRequestHandler) close() {
	c.conn.Close()
	c.listener.Close()
}

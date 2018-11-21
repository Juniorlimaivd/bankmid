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
	remoteIP    string
	remotePort  int
}

func newServerRequestHandler(port int) (*ServerRequestHandler, error) {
	tcpSRH := new(ServerRequestHandler)
	tcpSRH.listener, _ = net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))

	log.Printf("DNS Listen on %s", tcpSRH.listener.Addr().String())
	tcpSRH.conn, _ = tcpSRH.listener.Accept()
	remoteAddr := tcpSRH.conn.RemoteAddr().String()

	log.Printf("Accept a connection request from %s", remoteAddr)

	remoteIP, remotePort, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil, err
	}

	tcpSRH.remoteIP = remoteIP
	tcpSRH.remotePort, _ = strconv.Atoi(remotePort)

	tcpSRH.inToClient = bufio.NewWriter(tcpSRH.conn)
	tcpSRH.outToClient = bufio.NewReader(tcpSRH.conn)

	return tcpSRH, nil
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

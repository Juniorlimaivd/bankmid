package server

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
	connection  net.Conn
	outToClient *bufio.Reader
	inToClient  *bufio.Writer
	remoteAddr  string
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func newServerRequestHandler(port int) *ServerRequestHandler {
	tcpSRH := new(ServerRequestHandler)
	tcpSRH.listener, _ = net.Listen("tcp", ":"+strconv.Itoa(port))
	tcpSRH.remoteAddr = GetOutboundIP()
	log.Println("Server IP is : ", tcpSRH.remoteAddr)

	return tcpSRH
}

func (c *ServerRequestHandler) accept() {
	// log.Println("Listen on", tcpSRH.listener.Addr().String())
	c.connection, _ = c.listener.Accept()
	// log.Println("Accept a connection request from", conn.RemoteAddr())
	c.remoteAddr = GetOutboundIP()
	log.Println("Server IP is : ", c.remoteAddr)
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

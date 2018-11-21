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

func newServerRequestHandler(port int) (*ServerRequestHandler, error) {
	var err error
	tcpSRH := new(ServerRequestHandler)
	tcpSRH.listener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("It was not possible to create the server: %s", err)
	}
	tcpSRH.remoteAddr = GetOutboundIP()
	log.Printf("Server IP is : %s", tcpSRH.remoteAddr)

	return tcpSRH, nil
}

func (c *ServerRequestHandler) accept() {
	var err error
	// log.Printf("Listen on", tcpSRH.listener.Addr().String())
	c.connection, err = c.listener.Accept()
	if err != nil {
		log.Printf("It was not possible to accept %s", err)
	}
	// log.Printf("Accept a connection request from", conn.RemoteAddr())
	c.remoteAddr = GetOutboundIP()
	log.Printf("Server IP is : %s", c.remoteAddr)
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

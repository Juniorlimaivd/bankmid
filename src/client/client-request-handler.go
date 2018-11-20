package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

// ClientRequestHandler handles tcp connections
type ClientRequestHandler struct {
	host               string
	port               int
	sentMessageSize    int
	receiveMessageSize int
	conn               net.Conn
	rw                 *bufio.ReadWriter
}

func newClientRequestHandler(host string, port int) *ClientRequestHandler {
	return &ClientRequestHandler{
		host: host,
		port: port,
	}
}

func (c *ClientRequestHandler) connect() error {
	addr := c.host + ":" + strconv.Itoa(c.port)
	var err error
	c.conn, err = net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	c.rw = bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
	return nil
}

func (c *ClientRequestHandler) send(data []byte) error {
	c.rw.Write(data)
	c.rw.Flush()
	return nil
}

func (c *ClientRequestHandler) receive() []byte {
	data := make([]byte, 4096)
	n, err := c.rw.Read(data)

	if err != nil {
		fmt.Println(err)
		return make([]byte, 0)
	}
	return data[:n]
}

func (c *ClientRequestHandler) close() error {
	err := c.rw.Flush()
	if err != nil {
		return err
	}
	err = c.conn.Close()
	return err
}

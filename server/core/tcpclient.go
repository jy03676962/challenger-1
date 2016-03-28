package core

import (
	"bufio"
	"encoding/json"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net"
	"strings"
)

var _ = log.Printf

type TCPClient struct {
	*Hub
	id     string
	conn   *net.TCPConn
	reader *bufio.Reader
	server *TCPServer
	ch     chan *HubMap
	doneCh chan struct{}
}

func NewTCPClient(conn *net.TCPConn, server *TCPServer) *TCPClient {

	reader := bufio.NewReader(conn)
	ch := make(chan *HubMap, channelBufSize)
	doneCh := make(chan struct{})
	id := uuid.NewV4().String()

	return &TCPClient{server.Hub, id, conn, reader, server, ch, doneCh}
}

func (c *TCPClient) Addr() string {
	return c.conn.RemoteAddr().String()
}

func (c *TCPClient) ID() string {
	return c.id
}

func (c *TCPClient) Write(msg *HubMap) {
	select {
	case c.ch <- msg:
	default:
		close(c.doneCh)
	}
}

func (c *TCPClient) Listen() {
	if c.add() {
		go c.listenWrite()
		c.listenRead()
		c.conn.Close()
	}
}

func (c *TCPClient) listenWrite() {
	for {
		select {
		case msg := <-c.ch:
			c.writeMessage(msg)
		case <-c.doneCh:
			c.del()
			return
		}
	}
}

func (c *TCPClient) listenRead() {
	for {
		select {
		case <-c.doneCh:
			return
		default:
			msg, err := c.readMessage()
			if err == io.EOF {
				close(c.doneCh)
				return
			} else if err != nil {
				c.err(err)
				close(c.doneCh)
			} else if msg != nil {
				c.msg(msg)
			}
		}
	}
}

func (c *TCPClient) add() bool {
	e := TCPOutput{}
	e.ID = c.ID()
	e.Type = S_Add
	e.Client = c
	e.Addr = c.Addr()
	return c.send(&e)
}

func (c *TCPClient) msg(msg *HubMap) bool {
	e := TCPOutput{}
	e.ID = c.ID()
	e.Type = S_Msg
	e.Client = c
	e.Message = msg
	e.Addr = c.Addr()
	return c.send(&e)
}

func (c *TCPClient) del() bool {
	e := TCPOutput{}
	e.ID = c.ID()
	e.Type = S_Del
	e.Client = c
	e.Addr = c.Addr()
	return c.send(&e)
}

func (c *TCPClient) err(err error) bool {
	e := TCPOutput{}
	e.ID = c.ID()
	e.Type = S_Err
	e.Client = c
	e.Addr = c.Addr()
	e.Error = err
	return c.send(&e)
}

func (c *TCPClient) send(output *TCPOutput) bool {
	select {
	case c.TCPOutputCh <- output:
		return true
	case <-c.TCPServerQuitCh:
		return false
	}
}

func (c *TCPClient) writeMessage(message *HubMap) error {
	d, err := json.Marshal(message.Data())
	if err != nil {
		return err
	}
	buf := make([]byte, len(d)+2)
	for i := 1; i < len(buf)-1; i++ {
		buf[i] = d[i-1]
	}
	buf[0] = 60
	buf[len(buf)-1] = 62
	_, r := c.conn.Write(buf)
	return r
}

func (c *TCPClient) readMessage() (*HubMap, error) {
	b, err := c.reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != 60 {
		log.Println("hardware message must start with <")
		return nil, nil
	}
	msg := make([]byte, 0)
	for {
		b, err := c.reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == 62 {
			break
		}
		msg = append(msg, b)
	}
	if len(msg) == 0 {
		log.Println("got empty hardware message")
		return nil, nil
	}
	msgStr := string(msg)
	log.Println("got tcp message:", msgStr)
	if strings.HasPrefix(msgStr, "[UR]") {
		ret := NewHubMap()
		ret.SetCmd("upload_rev")
		ret.Set("data", msgStr[4:])
		return ret, nil
	} else {
		var f interface{}
		err := json.Unmarshal(msg, &f)
		if err != nil {
			return nil, err
		}
		return NewHubMapWithData(f.(map[string]interface{})), nil
	}
}

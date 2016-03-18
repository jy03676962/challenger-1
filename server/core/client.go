package core

// the skeleton of this file is borrowed from https://github.com/golang-samples/websocket

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxId int = 0

type Client struct {
	id       int
	username string
	ws       *websocket.Conn
	server   *Server
	ch       chan map[string]interface{}
	doneCh   chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan map[string]interface{}, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, "", ws, server, ch, doneCh}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) SetUsername(name string) {
	c.username = name
}

func (c *Client) GetUsername() string {
	return c.username
}

func (c *Client) Write(msg map[string]interface{}) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		case msg := <-c.ch:
			websocket.JSON.Send(c.ws, msg)

		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true
			return
		}
	}
}

func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true
			return

		default:
			var msg interface{}
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.messageCh <- &SocketEvent{msg.(map[string]interface{}), c}
			}
		}
	}
}

package core

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
)

var _ = log.Printf

const channelBufSize = 100

var maxId int = 0

type Client struct {
	*Hub
	id     int
	group  SocketGroupType
	ws     *websocket.Conn
	server *Server
	ch     chan *HubMap
	doneCh chan struct{}
}

func NewClient(ws *websocket.Conn, server *Server, group SocketGroupType) *Client {

	maxId++
	ch := make(chan *HubMap, channelBufSize)
	doneCh := make(chan struct{})

	return &Client{server.Hub, maxId, group, ws, server, ch, doneCh}
}

func (c *Client) Write(msg *HubMap) {
	select {
	case c.ch <- msg:
	default:
		close(c.doneCh)
	}
}

func (c *Client) Listen() {
	if c.add() {
		go c.listenWrite()
		c.listenRead()
		c.ws.Close()
	}
}

func (c *Client) listenWrite() {
	for {
		select {
		case msg := <-c.ch:
			log.Println("write socket message:", msg.Data())
			websocket.JSON.Send(c.ws, msg.Data())
		case <-c.doneCh:
			c.del()
			return
		}
	}
}

func (c *Client) listenRead() {
	for {
		select {
		case <-c.doneCh:
			return
		default:
			var msg interface{}
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				close(c.doneCh)
				return
			} else if err != nil {
				c.err(err)
			} else {
				c.msg(msg)
			}
		}
	}
}

func (c *Client) add() bool {
	e := SocketOutput{}
	e.Type = S_Add
	e.ID = c.id
	e.Group = c.group
	e.Client = c
	return c.send(&e)
}

func (c *Client) msg(msg interface{}) bool {
	e := SocketOutput{}
	e.Type = S_Msg
	e.ID = c.id
	e.Group = c.group
	e.Client = c
	e.SocketMessage = NewHubMapWithData(msg.(map[string]interface{}))
	if c.group == SG_Game {
		e.SocketMessage.Set("cid", c.id)
	}
	return c.send(&e)
}

func (c *Client) del() bool {
	e := SocketOutput{}
	e.Type = S_Del
	e.ID = c.id
	e.Group = c.group
	e.Client = c
	return c.send(&e)
}

func (c *Client) err(err error) bool {
	e := SocketOutput{}
	e.Type = S_Err
	e.ID = c.id
	e.Group = c.group
	e.Error = err
	e.Client = c
	return c.send(&e)
}

func (c *Client) send(output *SocketOutput) bool {
	select {
	case c.SocketOutputCh <- output:
		return true
	case <-c.ServerQuitCh:
		return false
	}
}

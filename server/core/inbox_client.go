package core

import (
	"log"
)

var _ = log.Println

type InboxClient struct {
	conn   InboxConnection
	doneCh chan bool
	id     int
	inbox  *Inbox
}

func NewInboxClient(conn InboxConnection, inbox *Inbox, id int) *InboxClient {
	client := InboxClient{}
	client.conn = conn
	client.id = id
	client.doneCh = make(chan bool)
	client.inbox = inbox
	return &client
}

func (c *InboxClient) Listen() {
	c.listenRead()
	c.conn.Close()
	c.inbox.RemoveClient(c.id)
}

func (c *InboxClient) Write(msg *InboxMessage, addr InboxAddress) {
	if !c.conn.Accept(addr) {
		return
	}
	go func() {
		e := c.conn.WriteJSON(msg)
		if e != nil {
			log.Printf("send message error:%v\n", e.Error())
			close(c.doneCh)
		}
	}()
}

func (c *InboxClient) listenRead() {
	for {
		select {
		case <-c.doneCh:
			return
		default:
			m := NewInboxMessage()
			e := c.conn.ReadJSON(m)
			if e != nil {
				log.Printf("read message error:%v\n", e.Error())
				return
			}
			c.inbox.ReceiveMessage(m)
		}
	}

}

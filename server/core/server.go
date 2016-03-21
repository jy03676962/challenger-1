package core

import (
	"golang.org/x/net/websocket"

	"log"
)

type Server struct {
	*Hub
	clients map[int]*Client
	match   *Match
}

func NewServer() *Server {
	s := Server{}
	s.Hub = NewHub()
	s.clients = make(map[int]*Client)
	s.match = NewMatch(s.Hub)
	return &s
}

func (s *Server) Run() {
	go s.match.Run()
	for {
		select {
		case output := <-s.SocketOutputCh:
			s.handleSocketOutput(output)
		case input := <-s.SocketInputCh:
			s.handleSocketInput(input)
		case msg := <-s.MatchOutputCh:
			s.handleMatchOutput(msg)
		case <-s.ServerQuitCh:
			return
		}
	}
}
func (s *Server) OnConnected(ws *websocket.Conn) {
	client := NewClient(ws, s)
	client.Listen()
}
func (s *Server) handleMatchOutput(msg *HubMap) {
	s.sendAll(msg)
}

func (s *Server) handleSocketInput(i *SocketInput) {
	if i.Broadcast {
		for _, client := range s.clients {
			client.Write(i.SocketMessage)
		}
	} else {
		if client, ok := s.clients[i.DestID]; ok {
			client.Write(i.SocketMessage)
		}
	}
}

func (s *Server) handleSocketOutput(e *SocketOutput) {
	switch e.Type {
	case S_Add:
		log.Printf("add client:%v\n", e.ID)
		s.clients[e.ID] = e.Client
	case S_Del:
		delete(s.clients, e.ID)
		hm := NewHubMap()
		hm.SetCmd("disconnect")
		hm.Set("cid", e.ID)
		s.MatchInputCh <- hm
	case S_Msg:
		s.handleSocketMessage(e)
	case S_Err:
		log.Println("Error:", e.Error.Error())
	}
}

func (s *Server) handleSocketMessage(e *SocketOutput) {
	msg := e.SocketMessage
	if msg.GetCmd() == "login" {
		data := NewHubMap()
		data.SetCmd("options")
		data.Set("options", s.Options)
		s.send(data, e.ID)
	}
	s.MatchInputCh <- msg
}

func (s *Server) sendAll(msg *HubMap) {
	input := SocketInput{}
	input.SocketMessage = msg
	input.Broadcast = true
	go s.doSend(&input)
}

func (s *Server) send(msg *HubMap, cid int) {
	input := SocketInput{}
	input.SocketMessage = msg
	input.Broadcast = false
	input.DestID = cid
	go s.doSend(&input)
}

func (s *Server) doSend(input *SocketInput) {
	select {
	case s.SocketInputCh <- input:
	case <-s.ServerQuitCh:
	}
}

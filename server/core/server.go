package core

import (
	"golang.org/x/net/websocket"

	"log"
)

type Server struct {
	*Hub
	clients    map[int]*Client
	apiClients map[int]*Client
}

func NewServer(hub *Hub) *Server {
	s := Server{}
	s.Hub = hub
	s.clients = make(map[int]*Client)
	s.apiClients = make(map[int]*Client)
	return &s
}

func (s *Server) Run() {
	for {
		select {
		case output := <-s.SocketOutputCh:
			s.handleSocketOutput(output)
		case input := <-s.SocketInputCh:
			s.handleSocketInput(input)
		case <-s.ServerQuitCh:
			return
		case <-s.MainQuitCh:
			return
		}
	}
}
func (s *Server) OnConnected(ws *websocket.Conn) {
	client := NewClient(ws, s, SG_Game)
	client.Listen()
}

func (s *Server) OnApiConnected(ws *websocket.Conn) {
	client := NewClient(ws, s, SG_Api)
	client.Listen()
}

func (s *Server) OnClientConnected(ws *websocket.Conn) {
	client := NewClient(ws, s, SG_Client)
	client.Listen()
}

func (s *Server) handleSocketInput(i *SocketInput) {
	var it map[int]*Client
	if i.Group == SG_Game {
		it = s.clients
	} else {
		it = s.apiClients
	}
	if i.Broadcast {
		for _, client := range it {
			client.Write(i.SocketMessage)
		}
	} else {
		if client, ok := it[i.DestID]; ok {
			client.Write(i.SocketMessage)
		}
	}
}

func (s *Server) handleSocketOutput(e *SocketOutput) {
	switch e.Type {
	case S_Add:
		log.Printf("add client:%v\n", e.ID)
		if e.Group == SG_Game {
			s.clients[e.ID] = e.Client
		} else if e.Group == SG_Api {
			s.apiClients[e.ID] = e.Client
		}
	case S_Del:
		log.Printf("del client:%v\n", e.ID)
		if e.Group == SG_Game {
			delete(s.clients, e.ID)
			hm := NewHubMap()
			hm.SetCmd("disconnect")
			hm.Set("cid", e.ID)
			s.MatchInputCh <- hm
		} else if e.Group == SG_Api {
			delete(s.clients, e.ID)
		}
	case S_Msg:
		if e.Group == SG_Game {
			s.handleGameSocketMessage(e)
		} else {
			s.handleApiSocketMessage(e)
		}
	case S_Err:
		log.Println("Error:", e.Error.Error())
	}
}

func (s *Server) handleGameSocketMessage(e *SocketOutput) {
	msg := e.SocketMessage
	if msg.GetCmd() == "login" {
		data := NewHubMap()
		data.SetCmd("options")
		data.Set("options", s.Options)
		s.send(data, e.ID)
	}
	s.MatchInputCh <- msg
}

func (s *Server) handleApiSocketMessage(e *SocketOutput) {
	i := TCPInput{}
	i.Message = e.SocketMessage
	go func() {
		select {
		case s.TCPInputCh <- &i:
		case <-s.TCPServerQuitCh:
		}
	}()
}

func (s *Server) sendAll(msg *HubMap) {
	input := SocketInput{}
	input.SocketMessage = msg
	input.Broadcast = true
	input.Group = SG_Game
	go s.doSend(&input)
}

func (s *Server) send(msg *HubMap, cid int) {
	input := SocketInput{}
	input.SocketMessage = msg
	input.Broadcast = false
	input.DestID = cid
	input.Group = SG_Game
	go s.doSend(&input)
}

func (s *Server) doSend(input *SocketInput) {
	select {
	case s.SocketInputCh <- input:
	case <-s.ServerQuitCh:
	}
}

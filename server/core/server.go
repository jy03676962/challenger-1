package core

import (
	"golang.org/x/net/websocket"

	"log"
)

type Server struct {
	*Hub
	clients map[SocketGroupType]map[int]*Client
}

func NewServer(hub *Hub) *Server {
	s := Server{}
	s.Hub = hub
	s.clients = make(map[SocketGroupType]map[int]*Client)
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

func (s *Server) OnAdminConnected(ws *websocket.Conn) {
	client := NewClient(ws, s, SG_Admin)
	client.Listen()
}

func (s *Server) handleSocketInput(i *SocketInput) {
	clientList := s.clients[i.Group]
	if clientList == nil {
		return
	}
	if i.Broadcast {
		for _, client := range clientList {
			client.Write(i.SocketMessage)
		}
	} else {
		if client, ok := clientList[i.DestID]; ok {
			client.Write(i.SocketMessage)
		}
	}
}

func (s *Server) handleSocketOutput(e *SocketOutput) {
	switch e.Type {
	case S_Add:
		log.Printf("add client:%v, group:%v\n", e.ID, e.Group.String())
		clientList := s.clients[e.Group]
		if clientList == nil {
			clientList = make(map[int]*Client)
			s.clients[e.Group] = clientList
		}
		clientList[e.ID] = e.Client
	case S_Del:
		log.Printf("del client:%v, group:%v\n", e.ID, e.Group.String())
		if clientList := s.clients[e.Group]; clientList != nil {
			delete(clientList, e.ID)
		}
		if e.Group == SG_Game {
			hm := NewHubMap()
			hm.SetCmd("disconnect")
			hm.Set("cid", e.ID)
			s.MatchInputCh <- hm
		}
	case S_Msg:
		if e.Group == SG_Game {
			s.handleGameSocketMessage(e)
		} else if e.Group == SG_Api {
			s.handleApiSocketMessage(e)
		} else if e.Group == SG_Client {
			s.handleClientSocketMessage(e)
		} else if e.Group == SG_Admin {
			s.handleAdminSocketMessage(e)
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
		s.send(data, e.ID, SG_Game)
	}
	s.MatchInputCh <- msg
}

func (s *Server) handleApiSocketMessage(e *SocketOutput) {
	i := TCPInput{Message: e.SocketMessage}
	go func() {
		select {
		case s.TCPInputCh <- &i:
		case <-s.TCPServerQuitCh:
		}
	}()
}

func (s *Server) handleClientSocketMessage(e *SocketOutput) {
	msg := e.SocketMessage
	if msg.GetCmd() == "query" {
	}
}

func (s *Server) handleAdminSocketMessage(e *SocketOutput) {
	msg := e.SocketMessage
	log.Printf("got message:%v\n", msg)
	switch msg.GetCmd() {
	case "init":
	case "queryHallData":
		data := NewHubMap()
		data.SetCmd("HallData")
		if teams := GetAllTeamsFromQueueWithLock(); teams != nil {
			data.Set("data", teams)
		}
		s.send(data, e.ID, SG_Admin)
	case "teamCutLine":
		teamID := msg.GetStr("teamID")
		TeamCutLine(teamID)
	case "teamRemove":
		teamID := msg.GetStr("teamID")
		TeamRemove(teamID)
	case "teamChangeMode":
		teamID := msg.GetStr("teamID")
		mode := msg.GetStr("mode")
		TeamChangeMode(teamID, mode)
	case "teamDelay":
		teamID := msg.GetStr("teamID")
		TeamDelay(teamID)
	}

}

func (s *Server) sendAll(msg *HubMap, group SocketGroupType) {
	input := SocketInput{}
	input.SocketMessage = msg
	input.Broadcast = true
	input.Group = group
	go s.doSend(&input)
}

func (s *Server) send(msg *HubMap, cid int, group SocketGroupType) {
	input := SocketInput{}
	input.SocketMessage = msg
	input.Broadcast = false
	input.DestID = cid
	input.Group = group
	log.Printf("will send:%v\n", input)
	go s.doSend(&input)
}

func (s *Server) doSend(input *SocketInput) {
	select {
	case s.SocketInputCh <- input:
	case <-s.ServerQuitCh:
	}
}

package core

import (
	"log"
	"net"
	"os"
)

type TCPServer struct {
	*Hub
	host    string
	clients map[string]*TCPClient
}

func NewTCPServer(host string, hub *Hub) *TCPServer {
	s := TCPServer{}
	s.host = host
	s.Hub = hub
	s.clients = make(map[string]*TCPClient)
	return &s
}

func (s *TCPServer) Run() {
	ln, err := net.Listen("tcp", s.host)
	if err != nil {
		log.Println("Api start error:", err.Error())
		os.Exit(1)
	}
	log.Println("start listen tcp:", s.host)
	defer ln.Close()
	go s.Accept(ln)
	for {
		select {
		case output := <-s.TCPOutputCh:
			s.handleTCPOutput(output)
		case input := <-s.TCPInputCh:
			s.handleTCPInput(input)
		case <-s.MainQuitCh:
			return
		case <-s.TCPServerQuitCh:
			return
		}

	}
}

func (s *TCPServer) Accept(ln net.Listener) {
	for {
		select {
		case <-s.TCPServerQuitCh:
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Error accepting: ", err.Error())
			} else {
				log.Println("got hardware connection")
				TCPConn := conn.(*net.TCPConn)
				client := NewTCPClient(TCPConn, s)
				go client.Listen()
			}
		}
	}
}

func (s *TCPServer) handleTCPOutput(output *TCPOutput) {
	switch output.Type {
	case S_Add:
		log.Printf("add client:%v\n", output.Addr)
		s.clients[output.ID] = output.Client
		s.sendAdd(output.Addr)
	case S_Del:
		log.Printf("del client:%v\n", output.Addr)
		delete(s.clients, output.ID)
		s.sendDel(output.Addr)
	case S_Msg:
		s.handleTCPMessage(output)
	case S_Err:
		log.Println("Error:", output.Error.Error())
		s.sendErr(output.Error.Error(), output.Addr)
	}

}

func (s *TCPServer) handleTCPMessage(output *TCPOutput) {
	msg := output.Message
	msg.Set("addr", output.Addr)
	go s.doSend(msg)
}

func (s *TCPServer) handleTCPInput(input *TCPInput) {
	for _, c := range s.clients {
		c.Write(input.Message)
	}
}

func (s *TCPServer) sendAdd(addr string) {
	d := NewHubMap()
	d.SetCmd("addTCP")
	d.Set("addr", addr)
	go s.doSend(d)
}
func (s *TCPServer) sendDel(addr string) {
	d := NewHubMap()
	d.SetCmd("delTCP")
	d.Set("addr", addr)
	go s.doSend(d)
}

func (s *TCPServer) sendErr(errMsg string, addr string) {
	d := NewHubMap()
	d.SetCmd("errTCP")
	d.Set("addr", addr)
	d.Set("msg", errMsg)
	go s.doSend(d)
}

func (s *TCPServer) doSend(data *HubMap) {
	i := SocketInput{}
	i.Group = SG_Api
	i.Broadcast = true
	i.SocketMessage = data
	select {
	case s.SocketInputCh <- &i:
	case <-s.ServerQuitCh:
	}
}

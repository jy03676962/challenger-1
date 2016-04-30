package core

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
	"strconv"
)

var _ = log.Println

type Srv struct {
	inbox            *Inbox
	queue            *Queue
	match            *Match
	db               *DB
	inboxMessageChan chan *InboxMessage
	pDict            map[string]*PlayerController
}

func NewSrv() *Srv {
	s := Srv{}
	s.inbox = NewInbox(s)
	s.queue = NewQueue()
	s.match = NewMatch(s)
	s.db = NewDb()
	s.inboxMessageChan = make(chan *InboxMessage, 1)
	return &s
}

func (s *Srv) Run(tcpAddr string, udpAddr string, dbPath string) {
e:
	s.db.connect(dbPath)
	if e != nil {
		log.Printf("open database error:%v\n", r.Error())
		os.Exit(1)
	}
	//go s.inbox.Run()
	go s.match.Run()
	go s.listenTcp(tcpAddr)
	go s.listenUdp(udpAddr)
	s.mainLoop()
}

func (s *Srv) ListenWebSocket(conn *websocket.Conn) {
	inbox.ListenConnection(NewInboxWsConnection(conn))
}

// http interface

func (s *Srv) AddTeam(c echo.Context) error {
	count, _ := strconv.Atoi(c.FormValue("count"))
	mode := c.FormValue("mode")
	t, err := s.queue.AddTeamToQueue(count, mode)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (s *Srv) ResetQueue(c echo.Context) error {
	err := s.queue.ResetQueue()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

// MARK: internal

func (s *Srv) mainLoop() {
	for {
		select {
		case msg := <-s.inboxMessageChan:
			s.handleInboxMessage(msg)
		}
	}
}

func (s *Srv) listenTcp(address string) {
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Println("resolve tcp address error:", err.Error())
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Println("listen tcp error:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println("tcp listen error: ", err.Error())
		} else {
			go inbox.ListenConnection(NewInboxTcpConnection(conn))
		}
	}
}

func (s *Srv) listenUdp(address string) {
	udpAddress, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("resolve udp address error:", err.Error())
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		log.Println("udp listen error: ", err.Error())
		os.Exit(1)
	}
	inbox.ListenConnection(NewInboxUdpConnection(conn))
}

// nonblock
func (s *Srv) onInboxMessageArrived(msg *InboxMessage) {
	s.inboxMessageChan <- msg
}

// nonblock, 下发match数据
func (s *Srv) onMatchUpdated() {
	b, e := json.Marshal(s.match)
	if e != nil {
		log.Printf("marshal match error:%v\n", e.Error())
		return
	}
	s.sendMsg("updateMatch", string(b), InboxAddressTypeSimulatorDevice, "")
}

func (s *Srv) saveMatch(d *MatchData) {
	s.db.SaveMatch(d)
}

func (s *Srv) handleInboxMessage(msg *InboxMessage) {
	id := msg.GetStr("ID")
	if len(id) == 0 {
		log.Printf("message has no ID:%v\n", msg.Data)
		return
	}
	cmd := msg.GetCmd()
	if len(cmd) == 0 {
		log.Printf("message has no cmd:%v\n", msg.Data)
		return
	}

	t := msg.Get("TYPE").(InboxAddressType)
	addr := InboxAddress{t, id}
	if t == InboxAddressTypeSimulatorDevice || t == InboxAddressTypeWearableDevice {
		key := fmt.Sprintf("%v:%v", t, id)
		if _, ok := s.pDict[key]; !ok {
			controller = NewPlayerController(addr)
			s.pDict[addr.String()] = controller
			json, _ := json.Marshal(controller)
			addrs := []InboxAddress{
				InboxAddress{InboxAddressTypeAdminDevice, ""},
				InboxAddress{InboxAddressTypeSimulatorDevice, ""},
			}
			s.sendMsgs("updatePlayerController", json, addrs)
		}
	}
}

func (s *Srv) sendMsg(cmd string, data interface{}, t InboxAddressType, id string) {
	addr := InboxAddress{t, id}
	sendMsgs(cmd, data, []InboxAddress{addr})
}

func (s *Srv) sendMsgs(cmd string, data interface{}, addrs []InboxAddress) {
	msg := NewInboxMessage{}
	msg.SetCmd(cmd)
	msg.Set("data", data)
	go s.inbox.Send(msg, addrs)
}

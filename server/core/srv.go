package core

import (
	"encoding/json"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
	"os"
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
	s.inbox = NewInbox(&s)
	s.queue = NewQueue(&s)
	s.match = NewMatch(&s)
	s.db = NewDb()
	s.inboxMessageChan = make(chan *InboxMessage, 1)
	return &s
}

func (s *Srv) Run(tcpAddr string, udpAddr string, dbPath string) {
	e := s.db.connect(dbPath)
	if e != nil {
		log.Printf("open database error:%v\n", e.Error())
		os.Exit(1)
	}
	//go s.inbox.Run()
	go s.match.Run()
	go s.listenTcp(tcpAddr)
	go s.listenUdp(udpAddr)
	s.mainLoop()
}

func (s *Srv) ListenWebSocket(conn *websocket.Conn) {
	s.inbox.ListenConnection(NewInboxWsConnection(conn))
}

// http interface

func (s *Srv) AddTeam(c echo.Context) error {
	count, _ := strconv.Atoi(c.FormValue("count"))
	mode := c.FormValue("mode")
	s.queue.AddTeamToQueue(count, mode)
	return c.JSON(http.StatusOK, nil)
}

func (s *Srv) ResetQueue(c echo.Context) error {
	s.queue.ResetQueue()
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
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("tcp listen error: ", err.Error())
		} else {
			go s.inbox.ListenConnection(NewInboxTcpConnection(conn))
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
	s.inbox.ListenConnection(NewInboxUdpConnection(conn))
}

func (s *Srv) onInboxMessageArrived(msg *InboxMessage) {
	s.inboxMessageChan <- msg
}

// nonblock, 下发match数据
func (s *Srv) onMatchUpdated(matchData []byte) {
	s.sendMsg("updateMatch", string(matchData), InboxAddressTypeSimulatorDevice, "")
}

func (s *Srv) saveMatch(d *MatchData) {
	s.db.saveMatch(d)
}

// nonblock, 下发queue数据
func (s *Srv) onQueueUpdated(queueData []Team) {
	s.sendMsgs("updateQueue", queueData, InboxAddressTypeAdminDevice)
}

func (s *Srv) handleInboxMessage(msg *InboxMessage) {
	shouldUpdatePlayerController := false
	if msg.RemoveAddress != nil && msg.RemoveAddress.Type.IsPlayerControllerType() {
		delete(s.pDict, msg.RemoveAddress.String())
		shouldUpdatePlayerController = true
	}
	if msg.AddAddress != nil && msg.AddAddress.Type.IsPlayerControllerType() {
		s.pDict[msg.AddAddress.String()] = NewPlayerController(*msg.AddAddress)
		shouldUpdatePlayerController = true
	}
	if shouldUpdatePlayerController {
		json, _ := json.Marshal(s.pDict)
		s.sendMsgs("updatePlayerController", json, InboxAddressTypeAdminDevice, InboxAddressTypeSimulatorDevice)
	}

	if msg.Address == nil {
		log.Printf("message has no address:%v\n", msg.Data)
		return
	}
	cmd := msg.GetCmd()
	if len(cmd) == 0 {
		log.Printf("message has no cmd:%v\n", msg.Data)
		return
	}
	switch msg.AddAddress.Type {
	case InboxAddressTypeSimulatorDevice:
		s.handleSimulatorMessage(msg)
	case InboxAddressTypeArduinoTestDevice:
		s.handleArduinoTestMessage(msg)
	case InboxAddressTypeAdminDevice:
		s.handleAdminMessage(msg)
	}
}

func (s *Srv) handleSimulatorMessage(msg *InboxMessage) {
	cmd := msg.GetCmd()
	if cmd == "init" {
		s.sendMsgToAddresses("options", GetOptions(), []InboxAddress{*msg.Address})
	}
}

func (s *Srv) handleArduinoTestMessage(msg *InboxMessage) {
	s.send(msg, []InboxAddress{InboxAddress{InboxAddressTypeArduinoDevice, ""}})
}

func (s *Srv) handleAdminMessage(msg *InboxMessage) {
	switch msg.GetCmd() {
	case "queryHallData":
		s.queue.TeamQueryData()
	case "teamCutLine":
		teamID := msg.GetStr("teamID")
		s.queue.TeamCutLine(teamID)
	case "teamRemove":
		teamID := msg.GetStr("teamID")
		s.queue.TeamRemove(teamID)
	case "teamChangeMode":
		teamID := msg.GetStr("teamID")
		mode := msg.GetStr("mode")
		s.queue.TeamChangeMode(teamID, mode)
	case "teamDelay":
		teamID := msg.GetStr("teamID")
		s.queue.TeamDelay(teamID)
	case "teamAddPlayer":
		teamID := msg.GetStr("teamID")
		s.queue.TeamAddPlayer(teamID)
	case "teamRemovePlayer":
		teamID := msg.GetStr("teamID")
		s.queue.TeamRemovePlayer(teamID)
	case "teamPrepare":
		teamID := msg.GetStr("teamID")
		s.queue.TeamPrepare(teamID)
	case "teamStart":
		teamID := msg.GetStr("teamID")
		s.queue.TeamStart(teamID)
	case "teamCall":
		teamID := msg.GetStr("teamID")
		s.queue.TeamCall(teamID)
	}
}

func (s *Srv) sendMsg(cmd string, data interface{}, t InboxAddressType, id string) {
	addr := InboxAddress{t, id}
	s.sendMsgToAddresses(cmd, data, []InboxAddress{addr})
}

func (s *Srv) sendMsgs(cmd string, data interface{}, types ...InboxAddressType) {
	addrs := make([]InboxAddress, len(types))
	for i, t := range types {
		addrs[i] = InboxAddress{t, ""}
	}
	s.sendMsgToAddresses(cmd, data, addrs)
}

func (s *Srv) sendMsgToAddresses(cmd string, data interface{}, addrs []InboxAddress) {
	msg := NewInboxMessage()
	msg.SetCmd(cmd)
	msg.Set("data", data)
	s.send(msg, addrs)
}

func (s *Srv) send(msg *InboxMessage, addrs []InboxAddress) {
	go s.inbox.Send(msg, addrs)
}

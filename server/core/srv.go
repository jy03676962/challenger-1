package core

import (
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var _ = log.Println

type pendingMatch struct {
	ids  []string
	mode string
}

type Srv struct {
	inbox            *Inbox
	queue            *Queue
	db               *DB
	inboxMessageChan chan *InboxMessage
	mChan            chan MatchEvent
	prepareDoneChan  chan bool
	pDict            map[string]*PlayerController
	aDict            map[string]*ArduinoController
	arduinoMode      ArduinoMode
	pm               *pendingMatch
	m                *Match
}

func NewSrv() *Srv {
	s := Srv{}
	s.inbox = NewInbox(&s)
	s.queue = NewQueue(&s)
	s.db = NewDb()
	s.inboxMessageChan = make(chan *InboxMessage, 1)
	s.mChan = make(chan MatchEvent)
	s.prepareDoneChan = make(chan bool)
	s.pDict = make(map[string]*PlayerController)
	s.aDict = make(map[string]*ArduinoController)
	s.arduinoMode = ArduinoModeFree
	s.pm = nil
	s.m = nil
	return &s
}

func (s *Srv) Run(tcpAddr string, udpAddr string, dbPath string) {
	e := s.db.connect(dbPath)
	if e != nil {
		log.Printf("open database error:%v\n", e.Error())
		os.Exit(1)
	}
	//go s.inbox.Run()
	go s.listenTcp(tcpAddr)
	go s.listenUdp(udpAddr)
	s.mainLoop()
}

func (s *Srv) ListenWebSocket(conn *websocket.Conn) {
	log.Println("got new ws connection")
	s.inbox.ListenConnection(NewInboxWsConnection(conn))
}

// http interface

func (s *Srv) AddTeam(c echo.Context) error {
	count, _ := strconv.Atoi(c.FormValue("count"))
	mode := c.FormValue("mode")
	id := s.queue.AddTeamToQueue(count, mode)
	d := map[string]interface{}{"id": id}
	return c.JSON(http.StatusOK, d)
}

func (s *Srv) ResetQueue(c echo.Context) error {
	id := s.queue.ResetQueue()
	d := map[string]interface{}{"id": id}
	return c.JSON(http.StatusOK, d)
}

// MARK: internal

func (s *Srv) mainLoop() {
	for {
		select {
		case msg := <-s.inboxMessageChan:
			s.handleInboxMessage(msg)
		case evt := <-s.mChan:
			s.handleMatchEvent(evt)
		case <-s.prepareDoneChan:
			pm := s.pm
			s.pm = nil
			s.startNewMatch(pm.ids, pm.mode)
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
	log.Println("listen tcp:", address)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("tcp listen error: ", err.Error())
		} else {
			log.Println("got new tcp connection")
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
	log.Println("listen udp:", address)
	s.inbox.ListenConnection(NewInboxUdpConnection(conn))
}

func (s *Srv) onInboxMessageArrived(msg *InboxMessage) {
	s.inboxMessageChan <- msg
}

func (s *Srv) onMatchEvent(evt MatchEvent) {
	s.mChan <- evt
}

func (s *Srv) saveMatch(d *MatchData) {
	s.db.saveMatch(d)
}

// nonblock, 下发queue数据
func (s *Srv) onQueueUpdated(queueData []Team) {
	log.Println("on queue updated")
	s.sendMsgs("HallData", queueData, InboxAddressTypeAdminDevice)
}

func (s *Srv) handleMatchEvent(evt MatchEvent) {
	switch evt.Type {
	case MatchEventTypeEnd:
		if evt.ID == s.m.ID {
			s.m = nil
		}
		for _, p := range s.pDict {
			if p.MatchID == evt.ID {
				p.MatchID = 0
			}
		}
		s.sendMsgs("matchStop", evt.ID, InboxAddressTypeSimulatorDevice, InboxAddressTypeAdminDevice)
	case MatchEventTypeUpdate:
		s.sendMsgs("updateMatch", evt.Data, InboxAddressTypeSimulatorDevice, InboxAddressTypeAdminDevice)
	}
}

func (s *Srv) handleInboxMessage(msg *InboxMessage) {
	shouldUpdatePlayerController := false
	if msg.RemoveAddress != nil && msg.RemoveAddress.Type.IsPlayerControllerType() {
		cid := msg.RemoveAddress.String()
		pc := s.pDict[cid]
		if pc.MatchID > 0 {
			s.m.OnMatchCmdArrived(msg)
		}
		delete(s.pDict, cid)
		shouldUpdatePlayerController = true
	}
	if msg.AddAddress != nil && msg.AddAddress.Type.IsPlayerControllerType() {
		pc := NewPlayerController(*msg.AddAddress)
		s.pDict[pc.ID] = pc
		shouldUpdatePlayerController = true
	}
	if shouldUpdatePlayerController {
		s.sendMsgs("ControllerData", s.getControllerData(), InboxAddressTypeAdminDevice, InboxAddressTypeSimulatorDevice)
	}

	if msg.RemoveAddress != nil && msg.RemoveAddress.Type.IsArduinoControllerType() {
		id := msg.RemoveAddress.String()
		delete(s.aDict, id)
	}

	if msg.AddAddress != nil && msg.AddAddress.Type.IsArduinoControllerType() {
		ac := NewArduinoController(*msg.AddAddress)
		s.aDict[ac.ID] = ac
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
	switch msg.Address.Type {
	case InboxAddressTypeSimulatorDevice:
		s.handleSimulatorMessage(msg)
	case InboxAddressTypeArduinoTestDevice:
		s.handleArduinoTestMessage(msg)
	case InboxAddressTypeAdminDevice:
		s.handleAdminMessage(msg)
	case InboxAddressTypeArduinoDevice:
		s.handleArduinoMessage(msg)
	}
}

func (s *Srv) handleArduinoMessage(msg *InboxMessage) {
	cmd := msg.GetCmd()
	switch cmd {
	case "confirm_mode_change":
		if ac, ok := s.aDict[msg.Address.String()]; ok {
			ac.Mode = msg.Get("mode").(ArduinoMode)
		}
	case "confirm_status_change":
		if ac, ok := s.aDict[msg.Address.String()]; ok {
			ac.Status = msg.Get("status").(ArduinoStatus)
		}
	}
}

func (s *Srv) handleSimulatorMessage(msg *InboxMessage) {
	cmd := msg.GetCmd()
	switch cmd {
	case "init":
		d := map[string]interface{}{
			"options": GetOptions(),
			"ID":      msg.Address.ID,
		}
		s.sendMsgToAddresses("init", d, []InboxAddress{*msg.Address})
	case "startMatch":
		mode := msg.GetStr("mode")
		ids := make([]string, 0)
		for _, pc := range s.pDict {
			if pc.Address.Type == InboxAddressTypeSimulatorDevice {
				ids = append(ids, pc.ID)
			}
		}
		s.startNewMatch(ids, mode)
	case "stopMatch", "playerMove", "playerStop":
		mid := uint(msg.Get("matchID").(float64))
		if s.m != nil && s.m.ID == mid {
			s.m.OnMatchCmdArrived(msg)
		}
	}
}

func (s *Srv) handleArduinoTestMessage(msg *InboxMessage) {
	s.send(msg, []InboxAddress{InboxAddress{InboxAddressTypeArduinoDevice, ""}})
}

func (s *Srv) handleAdminMessage(msg *InboxMessage) {
	switch msg.GetCmd() {
	case "init":
		s.sendMsg("init", nil, msg.Address.Type, msg.Address.ID)
	case "queryHallData":
		s.queue.TeamQueryData()
	case "queryControllerData":
		s.sendMsg("ControllerData", s.getControllerData(), msg.Address.Type, msg.Address.ID)
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
		mode := msg.GetStr("mode")
		ids := msg.Get("ids").(string)
		controllerIDs := strings.Split(ids, ",")
		s.queue.TeamStart(teamID)
		s.startNewMatch(controllerIDs, mode)
	case "teamCall":
		teamID := msg.GetStr("teamID")
		s.queue.TeamCall(teamID)
	case "arduinoModeChange":
	}
}

func (s *Srv) startNewMatch(controllerIDs []string, mode string) {
	if s.pm != nil || s.m != nil {
		return
	}
	if !s.canStartMatch() {
		s.pm = &pendingMatch{controllerIDs, mode}
		go s.prepare()
		return
	}
	mid := s.db.saveMatch(&MatchData{})
	for _, id := range controllerIDs {
		s.pDict[id].MatchID = mid
	}
	s.m = NewMatch(s, controllerIDs, mid, mode)
	go s.m.Run()
	log.Println("will send newMatch")
	s.sendMsgs("newMatch", mid, InboxAddressTypeAdminDevice, InboxAddressTypeSimulatorDevice)
}

func (s *Srv) canStartMatch() bool {
	for _, ac := range s.aDict {
		if ac.Mode != ArduinoModeFree || ac.Status != ArduinoStatusNormal {
			return false
		}
	}
	return true
}

func (s *Srv) prepare() {
	for {
		if s.canStartMatch() {
			s.prepareDoneChan <- true
			break
		}
		for _, ac := range s.aDict {
			if ac.Mode != ArduinoModeFree {
				msg := NewInboxMessage()
				msg.SetCmd("mode_change")
				msg.Set("mode", string(ArduinoModeFree))
				s.send(msg, []InboxAddress{ac.Address})
			} else if ac.Status != ArduinoStatusNormal {
				msg := NewInboxMessage()
				msg.SetCmd("status_change")
				msg.Set("status", string(ArduinoStatusNormal))
				s.send(msg, []InboxAddress{ac.Address})
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func (s *Srv) getControllerData() []PlayerController {
	r := make([]PlayerController, len(s.pDict))
	i := 0
	for _, pc := range s.pDict {
		r[i] = *pc
		i += 1
	}
	return r
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
	if data != nil {
		msg.Set("data", data)
	}
	s.send(msg, addrs)
}

func (s *Srv) send(msg *InboxMessage, addrs []InboxAddress) {
	s.inbox.Send(msg, addrs)
}

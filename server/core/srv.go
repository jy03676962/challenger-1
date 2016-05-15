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
	pDict            map[string]*PlayerController
	aDict            map[string]*ArduinoController
	mDict            map[uint]*Match
}

func NewSrv() *Srv {
	s := Srv{}
	s.inbox = NewInbox(&s)
	s.queue = NewQueue(&s)
	s.db = NewDb()
	s.inboxMessageChan = make(chan *InboxMessage, 1)
	s.mChan = make(chan MatchEvent)
	s.pDict = make(map[string]*PlayerController)
	s.aDict = make(map[string]*ArduinoController)
	s.mDict = make(map[uint]*Match)
	s.initArduinoControllers()
	return &s
}

func (s *Srv) Run(tcpAddr string, udpAddr string, dbPath string) {
	e := s.db.connect(dbPath)
	if e != nil {
		log.Printf("open database error:%v\n", e.Error())
		os.Exit(1)
	}
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

func (s *Srv) GetHistory(c echo.Context) error {
	d := s.db.getHistory(10)
	return c.JSON(http.StatusOK, d)
}

func (s *Srv) MatchStartAnswer(c echo.Context) error {
	mid, _ := strconv.Atoi(c.FormValue("mid"))
	d := s.db.startAnswer(mid, c.FormValue("eid"))
	s.sendMsgs("startAnswer", *d, InboxAddressTypePostgameDevice)
	return c.JSON(http.StatusOK, d)
}

func (s *Srv) MatchStopAnswer(c echo.Context) error {
	mid, _ := strconv.Atoi(c.FormValue("mid"))
	s.db.stopAnswer(mid)
	s.sendMsgs("stopAnswer", nil, InboxAddressTypePostgameDevice)
	return c.JSON(http.StatusOK, mid)
}

func (s *Srv) GetSurvey(c echo.Context) error {
	return c.JSON(http.StatusOK, GetSurvey())
}

func (s *Srv) UpdateQuestionInfo(c echo.Context) error {
	pid, _ := strconv.Atoi(c.FormValue("pid"))
	p := s.db.updateQuestionInfo(pid, c.FormValue("qid"), c.FormValue("aid"))
	s.sendMsgs("updatePlayerData", *p, InboxAddressTypeAdminDevice)
	return c.JSON(http.StatusOK, nil)
}

func (s *Srv) UpdatePlayerData(c echo.Context) error {
	pid, _ := strconv.Atoi(c.FormValue("pid"))
	p := s.db.updatePlayerData(pid, c.FormValue("name"), c.FormValue("eid"))
	s.sendMsgs("updatePlayerData", *p, InboxAddressTypeAdminDevice, InboxAddressTypePostgameDevice)
	return c.JSON(http.StatusOK, nil)
}

func (s *Srv) UpdateMatchData(c echo.Context) error {
	mid, _ := strconv.Atoi(c.FormValue("mid"))
	s.db.updateMatchData(mid, c.FormValue("eid"))
	return c.JSON(http.StatusOK, nil)
}

func (s *Srv) GetAnsweringMatchData(c echo.Context) error {
	d := s.db.getAnsweringMatchData()
	ret := make(map[string]interface{})
	if d == nil {
		ret["code"] = 1
	} else {
		ret["code"] = 0
		ret["data"] = d
	}
	return c.JSON(http.StatusOK, ret)
}

// internal

func (s *Srv) mainLoop() {
	for {
		select {
		case msg := <-s.inboxMessageChan:
			s.handleInboxMessage(msg)
		case evt := <-s.mChan:
			s.handleMatchEvent(evt)
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

// nonblock, 下发queue数据
func (s *Srv) onQueueUpdated(queueData []Team) {
	log.Println("on queue updated")
	s.sendMsgs("HallData", queueData, InboxAddressTypeAdminDevice)
}

func (s *Srv) handleMatchEvent(evt MatchEvent) {
	switch evt.Type {
	case MatchEventTypeEnd:
		delete(s.mDict, evt.ID)
		for _, p := range s.pDict {
			if p.MatchID == evt.ID {
				p.MatchID = 0
			}
		}
		d := evt.Data.(map[string]interface{})
		d["matchID"] = evt.ID
		s.queue.TeamFinishMatch(d["teamID"].(string))
		s.db.saveMatchData(d["matchData"].(*MatchData))
		s.sendMsgs("matchStop", d, InboxAddressTypeSimulatorDevice, InboxAddressTypeAdminDevice)
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
			s.mDict[pc.MatchID].OnMatchCmdArrived(msg)
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
		if controller := s.aDict[id]; controller != nil {
			controller.Online = false
			controller.ScoreUpdated = false
		}
		s.sendMsgs("removeTCP", msg.RemoveAddress, InboxAddressTypeArduinoTestDevice)
	}

	if msg.AddAddress != nil && msg.AddAddress.Type.IsArduinoControllerType() {
		if controller := s.aDict[msg.AddAddress.String()]; controller != nil {
			controller.Online = true
			if controller.NeedUpdateScore() {
				s.updateArduinoControllerScore(controller)
			}
		} else {
			log.Printf("Warning: get arduino connection not belong to list:%v\n", msg.AddAddress.String())
		}
		s.sendMsgs("addTCP", msg.AddAddress, InboxAddressTypeArduinoTestDevice)
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
	case InboxAddressTypeMainArduinoDevice, InboxAddressTypeSubArduinoDevice:
		s.handleArduinoMessage(msg)
	case InboxAddressTypePostgameDevice:
		s.handlePostGameMessage(msg)
	case InboxAddressTypeWearableDevice:
		s.handleWearableMessage(msg)
	}
}

func (s *Srv) handleWearableMessage(msg *InboxMessage) {
	msg.SetCmd("wearableLoc")
	for _, m := range s.mDict {
		m.OnMatchCmdArrived(msg)
	}
}

func (s *Srv) handleArduinoMessage(msg *InboxMessage) {
	cmd := msg.GetCmd()
	switch cmd {
	case "confirm_init_score":
		if controller := s.aDict[msg.Address.String()]; controller != nil {
			controller.ScoreUpdated = true
		}
	}
	if msg.GetCmd() != "init" {
		s.sends(msg, InboxAddressTypeArduinoTestDevice)
	}
}

func (s *Srv) handleSimulatorMessage(msg *InboxMessage) {
	cmd := msg.GetCmd()
	switch cmd {
	case "init":
		d := map[string]interface{}{
			"options": GetOptions(),
			"ID":      msg.Address.String(),
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
		s.startNewMatch(ids, mode, "", true)
	case "stopMatch", "playerMove", "playerStop":
		mid := uint(msg.Get("matchID").(float64))
		if match := s.mDict[mid]; match != nil {
			match.OnMatchCmdArrived(msg)
		}
	}
}

func (s *Srv) handleArduinoTestMessage(msg *InboxMessage) {
	s.sends(msg, InboxAddressTypeSubArduinoDevice, InboxAddressTypeMainArduinoDevice, InboxAddressTypeArduinoTestDevice)
}

func (s *Srv) handlePostGameMessage(msg *InboxMessage) {
	switch msg.GetCmd() {
	case "init":
		s.sendMsg("init", nil, msg.Address.ID, msg.Address.Type)
	}

}

func (s *Srv) handleAdminMessage(msg *InboxMessage) {
	switch msg.GetCmd() {
	case "init":
		s.sendMsg("init", nil, msg.Address.ID, msg.Address.Type)
	case "queryHallData":
		s.queue.TeamQueryData()
	case "queryControllerData":
		s.sendMsg("ControllerData", s.getControllerData(), msg.Address.ID, msg.Address.Type)
	case "queryQuestionCount":
		s.sendMsg("QuestionCount", len(GetSurvey().Questions), msg.Address.ID, msg.Address.Type)
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
		s.startNewMatch(controllerIDs, mode, teamID, false)
	case "teamCall":
		teamID := msg.GetStr("teamID")
		s.queue.TeamCall(teamID)
	case "arduinoModeChange":
		mode := ArduinoMode(msg.Get("mode").(float64))
		am := NewInboxMessage()
		am.SetCmd("mode_change")
		am.Set("mode", string(mode))
		s.sends(am, InboxAddressTypeMainArduinoDevice, InboxAddressTypeSubArduinoDevice)
	case "queryArduinoList":
		arduinolist := make([]ArduinoController, len(s.aDict))
		i := 0
		for _, controller := range s.aDict {
			arduinolist[i] = *controller
			i += 1
		}
		s.sendMsg("ArduinoList", arduinolist, msg.Address.ID, msg.Address.Type)
	case "stopMatch":
		mid := uint(msg.Get("matchID").(float64))
		if match := s.mDict[mid]; match != nil {
			match.OnMatchCmdArrived(msg)
		}
	}
}

func (s *Srv) startNewMatch(controllerIDs []string, mode string, teamID string, isSimulator bool) {
	md := s.db.newMatch()
	mid := md.ID
	for _, id := range controllerIDs {
		s.pDict[id].MatchID = mid
	}
	m := NewMatch(s, controllerIDs, md, mode, teamID, isSimulator)
	s.mDict[mid] = m
	go m.Run()
	s.sendMsgs("newMatch", mid, InboxAddressTypeAdminDevice, InboxAddressTypeSimulatorDevice)
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

func (s *Srv) sendMsg(cmd string, data interface{}, id string, t InboxAddressType) {
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

func (s *Srv) sends(msg *InboxMessage, types ...InboxAddressType) {
	addrs := make([]InboxAddress, len(types))
	for i, t := range types {
		addrs[i] = InboxAddress{t, ""}
	}
	s.send(msg, addrs)
}

// wall参数, 1主墙, 2小墙, 3二者同时
func (s *Srv) ledControl(wall int, mode string, ledT ...string) {
	if wall&1 > 0 {
		m := NewInboxMessage()
		m.SetCmd("led_ctrl")
		li := make([]map[string]string, 0)
		if ledT == nil {
			li = append(li, map[string]string{"wall": "M", "let_t": "1", "mode": mode})
		} else {
			for t := range ledT {
				li = append(li, map[string]string{"wall": "M", "let_t": t, "mode": mode})
			}
		}
		m.Set("led", li)
		s.sends(m, InboxAddressTypeMainArduinoDevice)
	}
	if wall&2 > 0 {
		m := NewInboxMessage()
		m.SetCmd("led_ctrl")
		li := make([]map[string]string, 0)
		li = append(li, map[string]string{"wall": "O1", "let_t": "1", "mode": mode})
		li = append(li, map[string]string{"wall": "O2", "let_t": "1", "mode": mode})
		li = append(li, map[string]string{"wall": "O3", "let_t": "1", "mode": mode})
		m.Set("led", li)
		s.sends(m, InboxAddressTypeSubArduinoDevice)
	}
}

func (s *Srv) ledControlByCell(x int, y int, mode string) {
	ids := GetOptions().mainArduinosByPos(x, y)
	m := NewInboxMessage()
	m.SetCmd("led_ctrl")
	m.Set("led", map[string]string{"wall": "M", "let_t": "1", "mode": mode})
	if len(ids) == 0 {
		return
	}
	addrs := make([]InboxAddress, len(ids))
	for i, id := range ids {
		addrs[i] = InboxAddress{InboxAddressTypeMainArduinoDevice, id}
	}
	s.send(m, addrs)
}

func (s *Srv) ledFlowEffect() {
	// TODO: 需要更新arduino表之后处理
}

func (s *Srv) send(msg *InboxMessage, addrs []InboxAddress) {
	s.inbox.Send(msg, addrs)
}

func (s *Srv) initArduinoControllers() {
	for _, main := range GetOptions().MainArduino {
		addr := InboxAddress{InboxAddressTypeMainArduinoDevice, main}
		controller := NewArduinoController(addr)
		s.aDict[addr.String()] = controller
	}
	for _, sub := range GetOptions().SubArduino {
		addr := InboxAddress{InboxAddressTypeSubArduinoDevice, sub}
		controller := NewArduinoController(addr)
		s.aDict[addr.String()] = controller
	}
}

func (s *Srv) updateArduinoControllerScore(controller *ArduinoController) {
	if !controller.NeedUpdateScore() {
		return
	}
	scoreInfo := GetScoreInfo()
	msg := NewInboxMessage()
	msg.SetCmd("init_score")
	msg.Set("score", scoreInfo)
	s.send(msg, []InboxAddress{controller.Address})
}

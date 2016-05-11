package core

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var _ = log.Printf

type MatchEventType int

const (
	MatchEventTypeEnd = iota
	MatchEventTypeUpdate
)

type MatchEvent struct {
	Type MatchEventType
	ID   uint
	Data interface{}
}

type Match struct {
	Member       []*Player       `json:"member"`
	Stage        string          `json:"stage"`
	TotalTime    float64         `json:"totalTime"`
	Elasped      float64         `json:"elasped"`
	WarmupTime   float64         `json:"warmupTime"`
	RampageTime  float64         `json:"rampageTime"`
	Mode         string          `json:"mode"`
	Gold         int             `json:"gold"`
	Energy       float64         `json:"energy"`
	OnButtons    map[string]bool `json:"onButtons"`
	RampageCount int             `json:"rampageCount"`
	Lasers       []*Laser        `json:"lasers"`
	ID           uint            `json:"id"`
	TeamID       string          `json:"teamID"`
	MaxEnergy    float64         `json:"maxEnergy"`

	offButtons    []string
	hiddenButtons map[string]float64
	goldDropTime  float64
	opt           *MatchOptions
	srv           *Srv
	msgCh         chan *InboxMessage
	closeCh       chan bool
	matchData     *MatchData
}

func NewMatch(s *Srv, controllerIDs []string, matchData *MatchData, mode string, teamID string) *Match {
	m := Match{}
	m.srv = s
	m.Member = make([]*Player, len(controllerIDs))
	for i, id := range controllerIDs {
		m.Member[i] = NewPlayer(id)
	}
	m.ID = matchData.ID
	m.matchData = matchData
	m.Stage = "before"
	m.opt = GetOptions()
	m.Mode = mode
	m.msgCh = make(chan *InboxMessage)
	m.closeCh = make(chan bool)
	m.TeamID = teamID
	m.MaxEnergy = GetOptions().MaxEnergy
	return &m
}

func (m *Match) Run() {
	dt := 33 * time.Millisecond
	tickChan := time.Tick(dt)
	m.Stage = "warmup"
	m.WarmupTime = m.opt.Warmup
	if m.Mode == "g" {
		m.TotalTime = m.opt.Mode1TotalTime
	} else {
		m.Gold = m.opt.Mode2InitGold[len(m.Member)-1]
	}
	for _, member := range m.Member {
		member.Pos = m.opt.RealPosition(m.opt.ArenaEntrance)
	}
	for {
		<-tickChan
		m.handleInputs()
		if m.Stage == "after" || m.Stage == "stop" {
			break
		}
		m.tick(dt)
		m.sync()
	}
	d := make(map[string]interface{})
	d["matchData"] = m.dumpMatchData()
	d["teamID"] = m.TeamID
	m.srv.onMatchEvent(MatchEvent{MatchEventTypeEnd, m.ID, d})
	close(m.closeCh)
}

func (m *Match) OnMatchCmdArrived(cmd *InboxMessage) {
	go func() {
		select {
		case m.msgCh <- cmd:
		case <-m.closeCh:
		}
	}()
}

func (m *Match) tick(dt time.Duration) {
	sec := dt.Seconds()
	m.Elasped += sec
	if m.Mode == "g" {
		m.TotalTime -= sec
	}
	if m.Stage == "warmup" {
		m.WarmupTime -= sec
		if m.WarmupTime <= 0 {
			m.WarmupTime = 0
			m.enterOngoing()
		}
	} else if m.Stage == "ongoing" {
		if m.Mode == "s" && m.goldDropTime > 0 {
			m.goldDropTime -= sec
			if m.goldDropTime <= 0 {
				m.Gold -= m.opt.Mode2GoldDropRate[len(m.Member)-1]
				m.goldDropTime = m.opt.Mode2GoldDropInterval
			}
		}
		for k, v := range m.hiddenButtons {
			v -= sec
			if v <= 0 {
				delete(m.hiddenButtons, k)
				m.OnButtons[k] = true
			}
		}
	}
	for _, player := range m.Member {
		m.playerTick(player, sec)
	}
	for _, laser := range m.Lasers {
		laser.Tick(sec)
	}
	m.checkRampage(sec)
	if m.Mode == "g" && m.TotalTime <= 0 || m.Mode == "s" && m.Gold <= 0 {
		m.enterAfter()
	}
}

func (m *Match) enterOngoing() {
	m.Stage = "ongoing"
	if m.Mode == "s" {
		m.goldDropTime = m.opt.Mode2GoldDropInterval
	}
	m.initLasers()
	m.initButtons()
}

func (m *Match) checkRampage(sec float64) {
	if m.RampageTime > 0 {
		m.RampageTime = math.Max(m.RampageTime-sec, 0)
		if m.RampageTime == 0 {
			m.leaveRampage()
		}
	} else if m.Energy >= m.opt.MaxEnergy {
		if len(m.Member) == 1 {
			m.enterRampage()
		} else {
			together := true
			p, pBool := m.opt.TilePosition(m.Member[0].Pos)
			if pBool {
				for i := 1; i < len(m.Member); i++ {
					pp, ppBool := m.opt.TilePosition(m.Member[i].Pos)
					if !ppBool || pp.X != p.X || pp.Y != p.Y {
						together = false
						break
					}
				}
				if together {
					m.enterRampage()
				}
			}
		}
	}
}

func (m *Match) enterAfter() {
	m.Stage = "after"
}

func (m *Match) sync() {
	b, _ := json.Marshal(m)
	m.srv.onMatchEvent(MatchEvent{MatchEventTypeUpdate, m.ID, string(b)})
}

func (m *Match) reset() {
	m.Member = make([]*Player, 0)
	m.Stage = "before"
	m.TotalTime = 0
	m.Elasped = 0
	m.WarmupTime = 0
	m.RampageTime = 0
	m.Mode = ""
	m.Gold = 0
	m.Energy = 0
	m.OnButtons = nil
	m.RampageCount = 0
	m.Lasers = nil
	m.offButtons = nil
	m.hiddenButtons = nil
	m.goldDropTime = 0
}

func (m *Match) handleInputs() bool {
	hasInputs := false
	for {
		select {
		case msg := <-m.msgCh:
			hasInputs = true
			m.handleInput(msg)
		default:
			return hasInputs
		}
	}
}

func (m *Match) handleInput(msg *InboxMessage) {
	if msg.RemoveAddress != nil {
		m.removePlayer(msg.RemoveAddress.String())
		return
	}
	cmd := msg.GetCmd()
	switch cmd {
	case "stopMatch":
		m.Stage = "stop"
	case "playerMove":
		if player := m.getPlayer(msg.Address.String()); player != nil {
			player.moving = true
			player.Direction = msg.GetStr("dir")
		}
	case "playerStop":
		if player := m.getPlayer(msg.Address.String()); player != nil {
			player.moving = false
		}
	}
}

func (m *Match) getPlayer(controllerID string) *Player {
	for _, player := range m.Member {
		if player.ControllerID == controllerID {
			return player
		}
	}
	return nil
}

func (m *Match) removePlayer(cid string) {
	destIdx := -1
	for idx, player := range m.Member {
		if player.ControllerID == cid {
			destIdx = idx
		}
	}
	if destIdx >= 0 {
		m.Member = append(m.Member[:destIdx], m.Member[destIdx+1:]...)
		idx := -1
		for i, laser := range m.Lasers {
			if laser.IsFollow(cid) {
				idx = i
			}
		}
		if idx >= 0 {
			m.Lasers = append(m.Lasers[:idx], m.Lasers[idx+1:]...)
		}
	}
	if len(m.Member) == 0 {
		m.Stage = "stop"
	}
}

func (m *Match) playerTick(player *Player, sec float64) {
	if player.InvincibleTime > 0 {
		player.InvincibleTime = math.Max(player.InvincibleTime-sec, 0)
	}
	moved := player.UpdatePos(sec, m.opt)
	if m.Stage != "ongoing" {
		return
	}
	if moved && player.Button != "" {
		m.consumeButton(player.Button, player)
	}
	if !moved {
		player.Stay(sec, m.opt, m.RampageTime > 0)
	}
}

func (m *Match) dumpMatchData() *MatchData {
	m.matchData.Mode = m.Mode
	m.matchData.Gold = m.Gold
	m.matchData.Elasped = m.Elasped
	m.matchData.Member = make([]PlayerData, 0)
	m.matchData.RampageCount = m.RampageCount
	m.matchData.AnswerType = MatchNotAnswer
	m.matchData.TeamID = m.TeamID
	m.matchData.ExternalID = ""
	for _, player := range m.Member {
		playerData := PlayerData{}
		playerData.Gold = player.Gold
		playerData.Energy = player.Energy
		playerData.LostGold = player.LostGold
		playerData.Combo = player.ComboCount
		strs := make([]string, 4)
		for i, c := range player.LevelData {
			strs[i] = strconv.Itoa(c)
		}
		playerData.LevelData = strings.Join(strs, ",")
		playerData.HitCount = player.HitCount
		playerData.Name = ""
		playerData.QuestionInfo = ""
		playerData.Answered = 0
		playerData.ExternalID = ""
		playerData.ControllerID = player.ControllerID
		m.matchData.Member = append(m.matchData.Member, playerData)
	}
	return m.matchData
}

func (m *Match) modeIndex() int {
	if m.Mode == "g" {
		return 0
	} else {
		return 1
	}
}

func (m *Match) initLasers() {
	m.Lasers = make([]*Laser, len(m.Member))
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	l := r.Perm(m.opt.ArenaWidth * m.opt.ArenaHeight)
	for i, player := range m.Member {
		loc := l[i]
		p := P{loc % m.opt.ArenaWidth, loc / m.opt.ArenaWidth}
		m.Lasers[i] = NewLaser(p, player, m)
		m.Lasers[i].Pause(m.opt.LaserAppearTime)
	}
}

func (m *Match) initButtons() {
	for _, player := range m.Member {
		player.Button = ""
		player.lastButton = ""
		player.ButtonLevel = 0
		player.ButtonTime = 0
	}
	l := len(m.opt.Buttons)
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	list := r.Perm(l)
	n := m.opt.InitButtonNum[len(m.Member)-1]
	m.OnButtons = make(map[string]bool)
	m.offButtons = make([]string, l-n)
	m.hiddenButtons = make(map[string]float64)
	for i, v := range list {
		id := strconv.Itoa(v)
		if i < n {
			m.OnButtons[id] = true
		} else {
			m.offButtons[i-n] = id
		}
	}
}

func (m *Match) consumeButton(btn string, player *Player) {
	player.LevelData[player.ButtonLevel] += 1
	if player.ButtonLevel > 0 {
		m.Gold += m.opt.GoldBonus[m.modeIndex()]
		player.Gold += 1
		if m.RampageTime <= 0 {
			sec := time.Since(player.lastHitTime).Seconds()
			player.lastHitTime = time.Now()
			var max float64
			if player.Combo == 0 {
				max = m.opt.FirstComboInterval[len(m.Member)-1]
			} else {
				max = m.opt.FirstComboInterval[len(m.Member)-1]
			}
			if sec <= max {
				player.Combo += 1
			} else {
				player.Combo = 0
			}
			extra := 0.0
			if player.Combo == 1 {
				extra = m.opt.FirstComboExtra
				player.ComboCount += 1
			} else if player.Combo > 1 {
				extra = m.opt.ComboExtra
			}
			delta := m.opt.EnergyBonus[player.ButtonLevel][len(m.Member)-1] + extra
			m.Energy = math.Min(m.opt.MaxEnergy, m.Energy+delta)
			player.Energy += delta
		}
	}
	player.lastButton = btn
	player.ButtonLevel = 0
	player.Button = ""
	player.ButtonTime = 0
	delete(m.OnButtons, btn)
	if m.RampageTime > 0 {
		m.offButtons = append(m.offButtons, btn)
		return
	}
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	i := r.Intn(len(m.offButtons))
	key := m.offButtons[i]
	m.offButtons[i] = btn
	t := m.opt.ButtonHideTime[m.modeIndex()]
	m.hiddenButtons[key] = t
}

func (m *Match) enterRampage() {
	m.RampageTime = m.opt.RampageTime[m.modeIndex()]
	for i := 0; i < len(m.opt.Buttons); i++ {
		k := strconv.Itoa(i)
		m.OnButtons[k] = true
	}
	for _, laser := range m.Lasers {
		laser.IsPause = true
		laser.pauseTime = m.RampageTime
	}
	m.offButtons = make([]string, 0)
	m.Energy = 0
	m.RampageCount += 1
	for _, player := range m.Member {
		player.Combo = 0
		player.lastHitTime = time.Unix(0, 0)
	}
}

func (m *Match) leaveRampage() {
	m.initButtons()
}

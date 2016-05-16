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
	isSimulator   bool
}

func NewMatch(s *Srv, controllerIDs []string, matchData *MatchData, mode string, teamID string, isSimulator bool) *Match {
	m := Match{}
	m.srv = s
	m.Member = make([]*Player, len(controllerIDs))
	for i, id := range controllerIDs {
		m.Member[i] = NewPlayer(id, isSimulator)
	}
	m.ID = matchData.ID
	m.matchData = matchData
	m.Stage = "before"
	m.opt = GetOptions()
	m.Mode = mode
	m.msgCh = make(chan *InboxMessage, 100)
	m.closeCh = make(chan bool)
	m.TeamID = teamID
	m.MaxEnergy = GetOptions().MaxEnergy
	m.isSimulator = isSimulator
	return &m
}

func (m *Match) Run() {
	dt := 33 * time.Millisecond
	tickChan := time.Tick(dt)
	if m.Mode == "g" {
		m.TotalTime = m.opt.Mode1TotalTime
	} else {
		m.Gold = m.opt.Mode2InitGold[len(m.Member)-1]
	}
	if m.isSimulator {
		for _, member := range m.Member {
			member.Pos = m.opt.RealPosition(m.opt.ArenaEntrance)
		}
	}
	m.WarmupTime = m.opt.Warmup
	m.setStage("warmup-1")
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
		m.TotalTime = math.Max(m.TotalTime-sec, 0)
	}
	if m.isWarmup() {
		m.WarmupTime = math.Max(m.WarmupTime-sec, 0)
	} else if m.isOngoing() {
		m.RampageTime = math.Max(m.RampageTime-sec, 0)
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
	m.updateStage()
}

func (m *Match) setStage(s string) {
	if m.Stage == s {
		return
	}
	switch s {
	case "warmup-1":
		m.srv.setAllWearableStatus("02")
		m.srv.ledControl(1, "3")
		m.srv.ledControl(2, "23")
	case "warmup-2":
		m.srv.ledControl(3, "4", "1", "2", "3")
	case "ongoing-low":
		if m.Stage == "ongoing-rampage" {
			m.initButtons()
		} else if m.isWarmup() {
			if m.Mode == "s" {
				m.goldDropTime = m.opt.Mode2GoldDropInterval
			}
			m.initLasers()
			m.initButtons()
		}
		m.srv.ledControl(3, "5")
	case "ongoing-high":
		m.srv.ledControl(3, "9")
	case "ongoing-full":
		if m.Mode == "g" {
			m.srv.ledControl(3, "19")
		} else {
			m.srv.ledControl(3, "20")
		}
		m.srv.setAllWearableStatus("03")
	case "ongoing-rampage":
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
		m.srv.ledRampageEffect()
		m.srv.setAllWearableStatus("04")
	case "ongoing-countdown":
		m.srv.ledControl(1, "47")
		m.srv.ledControl(2, "46")
	case "after":
		m.srv.ledFlowEffect()
	}
	m.Stage = s
}

func (m *Match) updateStage() {
	if m.RampageTime > 0 {
		m.setStage("ongoing-rampage")
		return
	}
	if m.WarmupTime > 0 {
		if m.WarmupTime <= m.opt.Warmup-m.opt.WarmupFirstStage {
			m.setStage("warmup-2")
		} else {
			m.setStage("warmup-1")
		}
		return
	}
	s := m.Stage
	r := m.Energy / m.opt.MaxEnergy
	if r < 0.8 {
		s = "ongoing-low"
	} else if r < 1 {
		s = "ongoing-high"
	} else {
		if len(m.Member) == 1 {
			s = "ongoing-rampage"
		} else {
			together := true
			if m.isSimulator {
				p, pBool := m.opt.TilePosition(m.Member[0].Pos)
				if pBool {
					for i := 1; i < len(m.Member); i++ {
						pp, ppBool := m.opt.TilePosition(m.Member[i].Pos)
						if !ppBool || pp.X != p.X || pp.Y != p.Y {
							together = false
							break
						}
					}
				} else {
					together = false
				}
			} else {
				tp := m.Member[0].tilePos
				for i := 1; i < len(m.Member); i++ {
					if m.Member[i].tilePos.X != tp.X || m.Member[i].tilePos.Y != tp.Y {
						together = false
					}
				}
			}
			if together {
				s = "ongoing-rampage"
			} else {
				s = "ongoing-full"
			}
		}
	}
	if m.Mode == "g" && m.opt.Mode1TotalTime-m.opt.Mode1CountDown < m.Elasped {
		s = "ongoing-countdown"
	}
	if m.Mode == "g" && m.TotalTime <= 0 || m.Mode == "s" && m.Gold <= 0 {
		s = "after"
	}
	m.setStage(s)
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
	case "wearableLoc":
		if player := m.getPlayer(msg.Address.String()); player != nil {
			loc, _ := strconv.Atoi(msg.GetStr("loc"))
			if loc > 0 {
				player.updateLoc(loc)
			}
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
	player.InvincibleTime = math.Max(player.InvincibleTime-sec, 0)
	if m.isSimulator {
		moved := player.UpdatePos(sec, m.opt)
		if !m.isOngoing() {
			return
		}
		if moved && player.Button != "" {
			m.consumeButton(player.Button, player)
		}
		if !moved {
			player.Stay(sec, m.opt, m.RampageTime > 0)
		}
	} else {
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
	count := len(m.opt.Buttons)
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	randList := r.Perm(count)
	n := m.opt.InitButtonNum[len(m.Member)-1]
	m.OnButtons = make(map[string]bool)
	m.offButtons = make([]string, count-n)
	m.hiddenButtons = make(map[string]float64)
	for i, j := range randList {
		id := m.opt.Buttons[j].Id
		if i < n {
			m.OnButtons[id] = true
		} else {
			m.offButtons[i-n] = id
		}
	}
	if !m.isSimulator {
		for id, _ := range m.OnButtons {
		}
	}
	if m.isSimulator {
	} else {
		//count := len(m.opt.MainArduino)
		//src := rand.NewSource(time.Now().UnixNano())
		//r := rand.New(src)
		//list := r.Perm(count)
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

func (m *Match) isWarmup() bool {
	return strings.HasPrefix(m.Stage, "warmup")
}

func (m *Match) isOngoing() bool {
	return strings.HasPrefix(m.Stage, "ongoing")
}

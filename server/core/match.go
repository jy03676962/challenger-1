package core

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	MATCH_CAPACITY = 4
)

type Match struct {
	*Hub         `json:"-"`
	Capacity     int             `json:"capacity"`
	Member       []*Player       `json:"member"`
	Stage        string          `json:"stage"`
	TotalTime    float64         `json:"totalTime"`
	Elasped      float64         `json:"elasped"`
	WarmupTime   float64         `json:"warmupTime"`
	RampageTime  float64         `json:"rampageTime"`
	Mode         int             `json:"mode"`
	Gold         float64         `json:"gold"`
	Energy       float64         `json:"energy"`
	OnButtons    map[string]bool `json:"onButtons"`
	RampageCount int             `json:"rampageCount"`
	Lasers       []*Laser        `json:"lasers"`

	offButtons    []string
	hiddenButtons map[string]float64
	goldDropTime  float64
}

func NewMatch(hub *Hub) *Match {
	m := Match{}
	m.Hub = hub
	m.Capacity = MATCH_CAPACITY
	m.Member = make([]*Player, 0)
	m.Stage = "before"
	return &m
}

func (m *Match) Run() {
	dt := 33 * time.Millisecond
	tickChan := time.Tick(dt)
	for {
		select {
		case <-m.ServerQuitCh:
			return
		case <-tickChan:
			isRunning := m.isRunning()
			needUpdate := m.handleInputs()
			if isRunning {
				needUpdate = true
				m.tick(dt)
			}
			if needUpdate {
				go m.sync()
			}
		}
	}
}

func (m *Match) tick(dt time.Duration) {
	if !m.isRunning() {
		return
	}
	sec := dt.Seconds()
	m.Elasped += sec
	if m.Mode == 1 {
		m.TotalTime -= sec
	}
	if m.Stage == "warmup" {
		m.WarmupTime -= sec
		if m.WarmupTime <= 0 {
			m.WarmupTime = 0
			m.enterOngoing()
		}
	} else if m.Stage == "ongoing" {
		if m.Mode == 2 && m.goldDropTime > 0 {
			m.goldDropTime -= sec
			if m.goldDropTime <= 0 {
				m.Gold -= m.Options.Mode2GoldDropRate[len(m.Member)-1]
				m.goldDropTime = m.Options.mode2GoldDropInterval
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
	if m.Mode == 1 && m.TotalTime <= 0 || m.Mode == 2 && m.Gold <= 0 {
		m.enterAfter()
	}
}

func (m *Match) enterOngoing() {
	m.Stage = "ongoing"
	if m.Mode == 2 {
		m.goldDropTime = m.Options.mode2GoldDropInterval
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
	} else if m.Energy >= m.Options.MaxEnergy {
		if len(m.Member) == 1 {
			m.enterRampage()
		} else {
			together := true
			p, pBool := m.Options.TilePosition(m.Member[0].Pos)
			if pBool {
				for i := 1; i < len(m.Member); i++ {
					pp, ppBool := m.Options.TilePosition(m.Member[i].Pos)
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
	data := NewHubMap()
	data.SetCmd("sync")
	data.Set("match", m)
	select {
	case m.MatchOutputCh <- data:
	case <-m.ServerQuitCh:
	}
}

func (m *Match) handleInputs() bool {
	hasInputs := false
	for {
		select {
		case input := <-m.MatchInputCh:
			hasInputs = true
			m.handleInput(input)
		default:
			return hasInputs
		}
	}
}

func (m *Match) handleInput(input *HubMap) {
	cmd := input.GetCmd()
	switch cmd {
	case "login":
		if m.isFull() || m.Stage != "before" {
			return
		}
		name := input.GetStr("name")
		id := input.Get("cid").(int)
		player := NewPlayer(name, id)
		m.Member = append(m.Member, player)
	case "startMatch":
		m.Mode = int(input.Get("mode").(float64))
		m.Stage = "warmup"
		m.WarmupTime = m.Options.Warmup
		if m.Mode == 1 {
			m.TotalTime = m.Options.Mode1TotalTime
		} else {
			m.Gold = m.Options.Mode2InitGold[len(m.Member)-1]
		}
		for _, member := range m.Member {
			member.Pos = m.Options.RealPosition(m.Options.ArenaEntrance)
		}
	case "playerMove":
		name := input.GetStr("name")
		dir := input.GetStr("dir")
		if player := m.getPlayer(name); player != nil {
			player.moving = true
			player.Direction = dir
		}
	case "playerStop":
		name := input.GetStr("name")
		if player := m.getPlayer(name); player != nil {
			player.moving = false
		}
	case "disconnect":
		cid := input.Get("cid").(int)
		m.removePlayer(cid)
	}
}

func (m *Match) getPlayer(name string) *Player {
	for _, player := range m.Member {
		if player.Name == name {
			return player
		}
	}
	return nil
}

func (m *Match) removePlayer(cid int) {
	if m.Stage == "after" {
		return
	}
	destIdx := -1
	destName := ""
	for idx, player := range m.Member {
		if player.clientID == cid {
			destIdx = idx
			destName = player.Name
		}
	}
	if destIdx >= 0 {
		m.Member = append(m.Member[:destIdx], m.Member[destIdx+1:]...)
		idx := -1
		for i, laser := range m.Lasers {
			if laser.IsFollow(destName) {
				idx = i
			}
		}
		if idx >= 0 {
			m.Lasers = append(m.Lasers[:idx], m.Lasers[idx+1:]...)
		}
	}
}

func (m *Match) playerTick(player *Player, sec float64) {
	if player.InvincibleTime > 0 {
		player.InvincibleTime = math.Max(player.InvincibleTime-sec, 0)
	}
	moved := player.UpdatePos(sec, m.Options)
	if m.Stage != "ongoing" {
		return
	}
	if moved && player.Button != "" {
		m.consumeButton(player.Button, player)
	}
	if !moved {
		player.Stay(sec, m.Options, m.RampageTime > 0)
	}
}

func (m *Match) isFull() bool {
	return len(m.Member) == m.Capacity
}

func (m *Match) isRunning() bool {
	return m.Stage == "ongoing" || m.Stage == "warmup"
}

func (m *Match) initLasers() {
	m.Lasers = make([]*Laser, len(m.Member))
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	l := r.Perm(m.Options.ArenaWidth * m.Options.ArenaHeight)
	for i, player := range m.Member {
		loc := l[i]
		p := P{loc % m.Options.ArenaWidth, loc / m.Options.ArenaWidth}
		m.Lasers[i] = NewLaser(p, player, m)
		m.Lasers[i].Pause(m.Options.laserAppearTime)
	}
}

func (m *Match) initButtons() {
	for _, player := range m.Member {
		player.Button = ""
		player.lastButton = ""
		player.ButtonLevel = 0
		player.ButtonTime = 0
	}
	l := len(m.Options.Buttons)
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	list := r.Perm(l)
	n := m.Options.initButtonNum[len(m.Member)-1]
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
		m.Gold += m.Options.GoldBonus[m.Mode-1]
		player.Gold += 1
		if m.RampageTime <= 0 {
			sec := time.Since(player.lastHitTime).Seconds()
			player.lastHitTime = time.Now()
			var max float64
			if player.Combo == 0 {
				max = m.Options.firstComboInterval[len(m.Member)-1]
			} else {
				max = m.Options.firstComboInterval[len(m.Member)-1]
			}
			if sec <= max {
				player.Combo += 1
			} else {
				player.Combo = 0
			}
			extra := 0.0
			if player.Combo == 1 {
				extra = m.Options.firstComboExtra
			} else if player.Combo > 1 {
				extra = m.Options.comboExtra
			}
			delta := m.Options.energyBonus[player.ButtonLevel][len(m.Member)-1] + extra
			m.Energy = math.Min(m.Options.MaxEnergy, m.Energy+delta)
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
	t := m.Options.buttonHideTime[m.Mode-1]
	m.hiddenButtons[key] = t
}

func (m *Match) enterRampage() {
	m.RampageTime = m.Options.rampageTime[m.Mode-1]
	for i := 0; i < len(m.Options.Buttons); i++ {
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

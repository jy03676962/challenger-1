package core

import (
  // "fmt"
  "math"
  "math/rand"
  "strconv"
  "time"
)

const (
  MATCH_CAPACITY = 4
)

/*
 * A Room object just hold a preparing state of a match,
 * Hoster: decide which client can fire off the match for web side
 * Member: all paticipants
 * Stage is before, warmup, ongoing, after
 */
type Match struct {
  Capacity      int             `json:"capacity"`
  Hoster        string          `json:"hoster"`
  Member        []*Player       `json:"member"`
  Stage         string          `json:"stage"`
  StartAt       time.Time       `json:"startAt"`
  TimeElapsed   float64         `json:"elasped"`
  RampageTime   time.Time       `json:"rampageTime"`
  RampageRemain float64         `json:"rampageRemain"`
  Rampage       bool            `json:"rampage"`
  Mode          int             `json:"mode"`
  Gold          float64         `json:"gold"`
  Energy        float64         `json:"energy"`
  LiveButtons   map[string]bool `json:"liveButtons"`
  RampageCount  int             `json:"rampageCount"`
  // private
  messageCh     chan string
  matchCh       chan string
  options       *MatchOptions
  backupButtons []string
  buttonDoneCh  chan struct{}
  buttonCh      chan string
  closeCh       <-chan struct{}
}

func NewMatch(matchCh chan string) *Match {
  m := Match{}
  m.Capacity = MATCH_CAPACITY
  m.Hoster = ""
  m.Member = make([]*Player, 0)
  m.Stage = "before"
  m.StartAt = time.Now()
  m.TimeElapsed = 0
  m.Rampage = false
  m.Mode = 0
  m.messageCh = make(chan string)
  m.matchCh = matchCh
  m.buttonDoneCh = make(chan struct{})
  m.buttonCh = make(chan string, 20)
  m.options = DefaultMatchOptions()
  m.RampageCount = 0
  return &m
}

func (m *Match) GetMessageCh() chan string {
  return m.messageCh
}

func (m *Match) GetOptions() *MatchOptions {
  return m.options
}

func (m *Match) AddMember(name string) bool {
  if m.IsFull() {
    return false
  }
  player := NewPlayer(name)
  m.Member = append(m.Member, player)
  return true
}

func (m *Match) PlayerMove(name string, dir string) {
  for _, player := range m.Member {
    if player.Name == name {
      player.moving = true
      player.Direction = dir
    }
  }
}

func (m *Match) PlayerStop(name string) {
  for _, player := range m.Member {
    if player.Name == name {
      player.moving = false
    }
  }
}

func (m *Match) RemoveMember(name string) bool {
  if len(name) == 0 {
    return false
  }
  if m.Hoster == name {
    m.Hoster = ""
  }
  for p, v := range m.Member {
    if v.Name == name {
      m.Member = append(m.Member[:p], m.Member[p+1:]...)
      return true
    }
  }
  return false
}

func (m *Match) IsFull() bool {
  return len(m.Member) == m.Capacity
}

func (m *Match) IsRunning() bool {
  return m.Stage == "ongoing" || m.Stage == "warmup"
}

func (m *Match) Start(mode int, closeCh <-chan struct{}) {
  m.Mode = mode
  m.Stage = "warmup"
  m.closeCh = closeCh
  m.StartAt = time.Now()
  for _, member := range m.Member {
    member.Pos = m.options.RealPosition(m.options.ArenaEntrance)
  }
  if m.Mode == 2 {
    m.Gold = m.options.Mode2InitGold[len(m.Member)-1]
  }
  m.resetButtons()
  go m.tickWarmup(m.options.Warmup)
  go m.gameLoop()
}

func (m *Match) useButton(btn string) {
  delete(m.LiveButtons, btn)
  if m.Rampage {
    return
  }
  i := rand.Intn(len(m.backupButtons))
  key := m.backupButtons[i]
  m.backupButtons[i] = btn
  t := m.options.buttonHideTime[m.Mode-1]
  go func() {
    c := time.After(time.Second * time.Duration(t))
    select {
    case <-c:
      m.buttonCh <- key
    case <-m.buttonDoneCh:
    case <-m.closeCh:
    }
  }()
}

func (m *Match) resetButtons() {
  for _, player := range m.Member {
    player.Button = ""
    player.lastButton = ""
    player.ButtonLevel = 0
    player.ButtonTime = 0
  }
  l := len(m.options.Buttons)
  list := rand.Perm(l)
  n := m.options.initButtonNum[len(m.Member)-1]
  m.LiveButtons = make(map[string]bool)
  m.backupButtons = make([]string, l-n)
  for i, v := range list {
    id := strconv.Itoa(v)
    if i < n {
      m.LiveButtons[id] = true
    } else {
      m.backupButtons[i-n] = id
    }
  }
}

func (m *Match) enterRampage() {
  close(m.buttonDoneCh)
  m.buttonCh = make(chan string, 20)
  m.buttonDoneCh = make(chan struct{})
  for i := 0; i < len(m.options.Buttons); i++ {
    k := strconv.Itoa(i)
    m.LiveButtons[k] = true
  }
  m.backupButtons = nil
  m.Rampage = true
  m.RampageTime = time.Now()
  m.Energy = 0
  m.RampageCount += 1
  t := m.options.rampageTime[m.Mode-1]
  go func() {
    c := time.After(time.Duration(t) * time.Second)
    select {
    case <-c:
      m.messageCh <- "RampageEnd"
    case <-m.closeCh:
    }
  }()
}

func (m *Match) leaveRampage() {
  m.resetButtons()
  m.Rampage = false
}

func (m *Match) endMatch() {
  m.Stage = "after"
  m.matchCh <- "matchEnd"
}

func (m *Match) tick() {
  m.TimeElapsed = time.Since(m.StartAt).Seconds()
  if m.Rampage {
    elapsed := time.Since(m.RampageTime).Seconds()
    m.RampageRemain = m.options.rampageTime[m.Mode-1] - elapsed
  }
  m.matchCh <- "tick"
}

func (m *Match) gameLoop() {
  tickChan := time.Tick(33 * time.Millisecond)
  for {
    select {
    case <-m.closeCh:
      return
    case msg := <-m.messageCh:
      switch msg {
      case "WarmupEnd":
        m.Stage = "ongoing"
        if m.Mode == 1 {
          go func() {
            c := time.After(time.Duration(m.options.Mode1TotalTime) * time.Second)
            select {
            case <-c:
              m.messageCh <- "MatchEnd"
            case <-m.closeCh:
            }
          }()
        }
      case "RampageEnd":
        m.leaveRampage()
      case "MatchEnd":
        m.endMatch()
        return
      }
      m.tick()
    case btn := <-m.buttonCh:
      m.LiveButtons[btn] = true
      m.tick()
    case <-tickChan:
      for _, player := range m.Member {
        if player.moving {
          delta := 1.0 / 30 * m.options.playerSpeed
          var dx, dy float64
          switch player.Direction {
          case "up":
            dx = 0
            dy = -delta
          case "right":
            dx = delta
            dy = 0
          case "down":
            dx = 0
            dy = delta
          case "left":
            dx = -delta
            dy = 0
          }
          minXY := (float64(m.options.ArenaBorder) + m.options.PlayerSize) / 2
          maxX := float64((m.options.ArenaBorder+m.options.ArenaCellSize)*m.options.ArenaWidth) - minXY
          maxY := float64((m.options.ArenaBorder+m.options.ArenaCellSize)*m.options.ArenaHeight) - minXY
          x := MinMaxfloat64(player.Pos.X+dx, minXY, maxX)
          y := MinMaxfloat64(player.Pos.Y+dy, minXY, maxY)
          rect := Rect{
            float64(x) - float64(m.options.PlayerSize)/2,
            float64(y) - float64(m.options.PlayerSize)/2,
            float64(m.options.PlayerSize),
            float64(m.options.PlayerSize),
          }
          if !m.options.CollideWall(&rect) {
            player.Pos = RP{x, y}
          }
          if player.Button != "" {
            if player.ButtonLevel > 0 {
              m.Gold += m.options.GoldBonus[m.Mode-1]
              player.Gold += 1
            }
            if !m.Rampage {
              delta := m.options.energyBonus[player.ButtonLevel][len(m.Member)-1]
              m.Energy = math.Min(m.options.MaxEnergy, m.Energy+delta)
              player.LevelData[player.ButtonLevel] += 1
              player.Energy += delta
            }
            m.useButton(player.Button)
            player.lastButton = player.Button
            player.Button = ""
            player.ButtonTime = 0
            player.ButtonLevel = 0
          }
        } else if m.Stage == "ongoing" {
          if player.Button != "" {
            player.ButtonTime += 1.0 / 30
            t := player.ButtonTime
            level := 0
            if m.Rampage {
              if t > m.options.TRampage {
                level = 1
              }
            } else {
              if t < m.options.T1 {
                level = 0
              } else if t < m.options.T2 {
                level = 1
              } else if t < m.options.T3 {
                level = 2
              } else {
                level = 3
              }
            }
            player.ButtonLevel = level
          } else {
            rect := Rect{
              float64(player.Pos.X) - float64(m.options.PlayerSize)/2,
              float64(player.Pos.Y) - float64(m.options.PlayerSize)/2,
              float64(m.options.PlayerSize),
              float64(m.options.PlayerSize),
            }
            buttons := m.options.PressingButtons(&rect)
            if buttons != nil {
              var id string
              if len(buttons[1]) > 0 {
                if buttons[0] == player.lastButton {
                  id = buttons[1]
                } else {
                  id = buttons[0]
                }
              } else {
                id = buttons[0]
              }
              player.Button = id
            }
          }
        }
      }
      if !m.Rampage && m.Energy >= m.options.MaxEnergy {
        if len(m.Member) == 1 {
          m.enterRampage()
        } else {
          together := true
          p, pBool := m.options.TilePosition(m.Member[0].Pos)
          if pBool {
            for i := 1; i < len(m.Member); i++ {
              pp, ppBool := m.options.TilePosition(m.Member[i].Pos)
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
      m.tick()
    }
  }
}

func (m *Match) tickWarmup(timeout int) {
  c := time.After(time.Duration(timeout) * time.Second)
  select {
  case <-c:
    m.messageCh <- "WarmupEnd"
  case <-m.closeCh:
  }
}

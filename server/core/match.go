package core

import (
  "fmt"
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
  Capacity    int       `json:"capacity"`
  Hoster      string    `json:"hoster"`
  Member      []*Player `json:"member"`
  Stage       string    `json:"stage"`
  StartAt     time.Time `json:"startAt"`
  TimeElapsed float64   `json:"elasped"`
  Rampage     bool      `json:"rampage"`
  Mode        int       `json:"mode"`
  Gold        float64   `json:"gold"`
  Energy      float64   `json:"energy"`
  // private
  messageCh chan string
  matchCh   chan string
  options   *MatchOptions
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
  m.options = DefaultMatchOptions()
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

func (m *Match) Start(mode int) {
  m.Mode = mode
  m.Stage = "warmup"
  m.StartAt = time.Now()
  for _, member := range m.Member {
    member.Pos = m.options.RealPosition(m.options.ArenaEntrance)
  }
  go m.tickWarmup(m.options.Warmup)
  go m.gameLoop()
}

func (m *Match) gameLoop() {
  tickChan := time.Tick(33 * time.Millisecond)
  for {
    select {
    case message := <-m.messageCh:
      switch message {
      case "WarmupEnd":
        m.Stage = "ongoing"
      }
      m.matchCh <- "tick"
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
          player.lastButton = player.Button
          player.Button = ""
          player.ButtonTime = 0
        } else {
          if player.Button != "" {
            player.ButtonTime += 1.0 / 30
            fmt.Println(player.ButtonTime)
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
      m.TimeElapsed = time.Since(m.StartAt).Seconds()
      m.matchCh <- "tick"
    }
  }
}

func (m *Match) tickWarmup(timeout int) {
  time.Sleep(time.Duration(timeout) * time.Second)
  m.messageCh <- "WarmupEnd"
}

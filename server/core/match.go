package core

import (
  "time"
)

var COLOR_ARRAY [4]string

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
  // private
  messageCh chan string
  matchCh   chan string
  options   *MatchOptions
}

func NewMatch(matchCh chan string) *Match {
  COLOR_ARRAY = [...]string{"red", "yellow", "green", "blue"}
  options := DefaultMatchOptions()
  messageCh := make(chan string)
  return &Match{
    MATCH_CAPACITY,
    "",
    make([]*Player, 0),
    "before",
    time.Now(),
    0,
    messageCh,
    matchCh,
    options,
  }
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
  for _, color := range COLOR_ARRAY {
    used := false
    for _, member := range m.Member {
      if member.Color == color {
        used = true
      }
    }
    if !used {
      player.Color = color
      break
    }
  }
  m.Member = append(m.Member, player)
  return true
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

func (m *Match) Start() {
  // TODO: init match and start
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
      m.TimeElapsed = time.Since(m.StartAt).Seconds()
      m.matchCh <- "tick"
    }
  }
}

func (m *Match) tickWarmup(timeout int) {
  time.Sleep(time.Duration(timeout) * time.Second)
  m.messageCh <- "WarmupEnd"
}

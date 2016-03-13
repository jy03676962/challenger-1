package core

import (
  "time"
)

const MATCH_CAPACITY = 4

/*
 * A Room object just hold a preparing state of a match,
 * Hoster: decide which client can fire off the match for web side
 * Member: all paticipants
 * Stage is before, warmup, ongoing, after
 */
type Match struct {
  Capacity    int       `json:"capacity"`
  Hoster      string    `json:"hoster"`
  Member      StrSlice  `json:"member"`
  Stage       string    `json:"stage"`
  StartAt     time.Time `json:"startAt"`
  TimeElapsed float64   `json:"elasped"`
  // private
  messageCh chan string
  matchCh   chan string
  options   *MatchOptions
}

func NewMatch(matchCh chan string) *Match {
  options := DefaultMatchOptions()
  messageCh := make(chan string)
  return &Match{
    MATCH_CAPACITY,
    "",
    make(StrSlice, 0),
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
  m.Member = append(m.Member, name)
  return true
}

func (m *Match) IsFull() bool {
  return len(m.Member) == m.Capacity
}

func (m *Match) Start() {
  // TODO: init match and start
  m.Stage = "warmup"
  m.StartAt = time.Now()
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

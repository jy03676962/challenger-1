package core

const MATCH_CAPACITY = 4

/*
 * A Room object just hold a preparing state of a match,
 * Hoster: decide which client can fire off the match for web side
 * Member: all paticipants
 * Stage is before, ongoing, after
 */
type Match struct {
  Capacity int           `json:"capacity"`
  Hoster   string        `json:"hoster"`
  Member   []string      `json:"member"`
  Stage    string        `json:"stage"`
  Options  *MatchOptions `json:"mapOptions"`
}

func NewMatch() *Match {
  options := DefaultMatchOptions()
  return &Match{
    MATCH_CAPACITY,
    "",
    make([]string, 0),
    "before",
    options,
  }
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
  m.Stage = "ongoing"
}

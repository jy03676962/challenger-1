package core

type Player struct {
  Name        string  `json:"name"`
  Pos         RP      `json:"pos"`
  Direction   string  `json:"dir"` // values:up,right,down,left
  Button      string  `json:"button"`
  ButtonTime  float64 `json:"buttonTime"`
  ButtonLevel int     `json:"buttonLevel"`
  moving      bool
  lastButton  string
}

func NewPlayer(name string) *Player {
  return &Player{
    name,
    RP{0, 0},
    "up",
    "",
    0,
    0,
    false,
    "",
  }
}

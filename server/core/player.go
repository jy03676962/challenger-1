package core

type Player struct {
  Name      string `json:"name"`
  Color     string `json:"color"`
  Pos       RP     `json:"pos"`
  Direction string `json:"dir"` // values:up,right,down,left
  moving    bool
}

func NewPlayer(name string) *Player {
  return &Player{
    name,
    "",
    RP{0, 0},
    "up",
    false,
  }
}

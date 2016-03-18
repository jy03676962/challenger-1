package core

type Player struct {
	Name           string  `json:"name"`
	Pos            RP      `json:"pos"`
	Direction      string  `json:"dir"` // values:up,right,down,left
	Button         string  `json:"button"`
	ButtonTime     float64 `json:"buttonTime"`
	ButtonLevel    int     `json:"buttonLevel"`
	Gold           float64 `json:"gold"`
	Energy         float64 `json:"energy"`
	LevelData      [4]int  `json:"levelData"`
	HitCount       int     `json:"hitCount"`
	IsInvincible   bool    `json:"isInvincible"`
	moving         bool
	lastButton     string
	invincibleTime float64
}

func NewPlayer(name string) *Player {
	p := Player{}
	p.Name = name
	p.Pos = RP{0, 0}
	p.Direction = "up"
	p.Button = ""
	p.ButtonTime = 0
	p.ButtonLevel = 0
	p.Gold = 0
	p.IsInvincible = false
	p.invincibleTime = 0
	p.Energy = 0
	p.LevelData = [4]int{0, 0, 0, 0}
	p.HitCount = 0
	p.moving = false
	p.lastButton = ""
	return &p
}

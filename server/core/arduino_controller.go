package core

import (
	"strconv"
	"strings"
)

type ArduinoMode int

const ArduinoModeUnknown ArduinoMode = -1

const (
	ArduinoModeOff = iota
	ArduinoModeOn
	ArduinoModeScan
	ArduinoModeFree
)

type ArduinoController struct {
	Address      InboxAddress `json:"address"`
	ID           string       `json:"id"`
	Mode         ArduinoMode  `json:"mode"`
	Online       bool         `json:"online"`
	ScoreUpdated bool         `json:"scoreUpdated"`
	//private
	dir int
	x   int
	y   int
	t   string
	num int
}

func NewArduinoController(addr InboxAddress) *ArduinoController {
	a := ArduinoController{}
	a.Address = addr
	a.ID = addr.String()
	a.Mode = ArduinoModeUnknown
	a.Online = false
	a.ScoreUpdated = false
	id := a.Address.ID
	li := strings.Split(id, "-")
	a.x, _ = strconv.Atoi(li[1])
	a.y, _ = strconv.Atoi(li[2])
	a.dir, _ = strconv.Atoi(li[3])
	a.t = li[4]
	a.num, _ = strconv.Atoi(li[5])
	return &a
}

func (c *ArduinoController) NeedUpdateScore() bool {
	return c.Address.Type == InboxAddressTypeMainArduinoDevice || c.ScoreUpdated == false
}

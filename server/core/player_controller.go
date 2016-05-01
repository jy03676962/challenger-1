package core

type PCStatus int

const (
	PCStatusOffline = iota
	PCStatusIdle
	PCStatusUsing
)

type PlayerController struct {
	Address InboxAddress `json:"address"`
	Status  PCStatus     `json:"status"`
}

func NewPlayerController(addr InboxAddress, st PCStatus) *PlayerController {
	c := PlayerController{}
	c.Address = addr
	c.Status = st
	return &c
}

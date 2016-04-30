package core

type PlayerController struct {
	Address InboxAddress `json:"address"`
}

func NewPlayerController(addr InboxAddress) *PlayerController {
	c := PlayerController{}
	c.Address = addr
	return &c
}

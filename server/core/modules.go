package core

type ArenaPosition struct {
	X int
	Y int
}

type P ArenaPosition

type RealPosition struct {
	X float64
	Y float64
}

type RP RealPosition

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

type ArenaWall struct {
	P1 P
	P2 P
}

type W ArenaWall

type SocketEventType int

const (
	S_Add SocketEventType = iota
	S_Del SocketEventType = 1 << iota
	S_Msg
	S_Err
)

type SocketOutput struct {
	Client        *Client
	ID            int
	Type          SocketEventType
	Error         error
	SocketMessage *HubMap
}

type SocketInput struct {
	Broadcast     bool
	DestID        int
	SocketMessage *HubMap
}

type Button struct {
	Id string `json:"id"`
	R  Rect   `json:"r"`
}

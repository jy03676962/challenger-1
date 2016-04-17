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

type SocketGroupType int

const (
	SG_Game   SocketGroupType = iota
	SG_Api    SocketGroupType = 1 << iota
	SG_Client SocketGroupType = 2 << iota
)

type SocketOutput struct {
	Client        *Client
	ID            int
	Group         SocketGroupType
	Type          SocketEventType
	Error         error
	SocketMessage *HubMap
}

type SocketInput struct {
	Broadcast     bool
	DestID        int
	Group         SocketGroupType
	SocketMessage *HubMap
}

type TCPOutput struct {
	Client  *TCPClient
	ID      string
	Addr    string
	Type    SocketEventType
	Message *HubMap
	Error   error
}

type TCPInput struct {
	Message *HubMap
}

type Button struct {
	Id string `json:"id"`
	R  Rect   `json:"r"`
}

type PlayerData struct {
	Name   string  `json:"name"`
	Gold   int     `json:"gold"`
	Energy float64 `json:"energy"`
	Combo  int     `json:"combo"`
	Grade  string  `json:"grade"`
	Level  int     `json:"level"`
}

type MatchData struct {
	Mode    int           `json:"mode"`
	Elasped float64       `json:"elasped"`
	Gold    float64       `json:"gold"`
	Member  []*PlayerData `json:"member"`
}

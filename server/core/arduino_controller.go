package core

type ArduinoMode int

const ArduinoModeUnknown ArduinoMode = -1

const (
	ArduinoModeOff = iota
	ArduinoModeOn
	ArduinoModeScan
	ArduinoModeFree
)

type ArduinoStatus int

const ArduinoStatusUnknown ArduinoStatus = -1

const (
	ArduinoStatusIdle = iota
	ArduinoStatusNormal
)

type ArduinoController struct {
	Address InboxAddress  `json:"address"`
	ID      string        `json:"id"`
	Mode    ArduinoMode   `json:"mode"`
	Status  ArduinoStatus `json:"status"`
}

func NewArduinoController(addr InboxAddress) *ArduinoController {
	a := ArduinoController{}
	a.Address = addr
	a.ID = addr.String()
	a.Mode = ArduinoModeUnknown
	a.Status = ArduinoStatusUnknown
	return &a
}

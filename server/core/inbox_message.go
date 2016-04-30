package core

import (
	"encoding/json"
)

type InboxMessage struct {
	Data map[string]interface{}
}

func NewInboxMessage() *InboxMessage {
	return &InboxMessage{make(map[string]interface{})}
}

func (message *InboxMessage) Get(key string) interface{} {
	if v, ok := message.Data[key]; ok {
		return v
	}
	return nil
}

func (message *InboxMessage) Set(key string, value interface{}) {
	message.Data[key] = value
}

func (message *InboxMessage) GetStr(key string) string {
	if v, ok := message.Data[key]; ok {
		return v.(string)
	}
	return ""
}

func (message *InboxMessage) GetCmd() string {
	return message.GetStr("cmd")
}

func (message *InboxMessage) SetCmd(v string) {
	message.Set("cmd", v)
}

func (message *InboxMessage) Marshal() (b []byte, e error) {
	b, e = json.Marshal(message.Data)
	return
}

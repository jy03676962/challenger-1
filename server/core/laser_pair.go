package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var _ = log.Printf

type ReceiverInfo struct {
	ID    string `json:"id"`
	Idx   string `json:"idx"`
	Valid int    `json:"valid"`
}

type LaserPair struct {
	m map[string]ReceiverInfo
}

var laserPair = LoadLaserPair()

func GetLaserPair() *LaserPair {
	return laserPair
}

func LoadLaserPair() *LaserPair {
	m := make(map[string]ReceiverInfo)
	b, e := ioutil.ReadFile("./laser.json")
	if os.IsNotExist(e) {
		return &LaserPair{m}
	}
	if e != nil {
		log.Printf("parse laser pair error:%v\n", e.Error())
		os.Exit(1)
	}
	e = json.Unmarshal(b, &m)
	if e != nil {
		log.Printf("parse laser pair error:%v\n", e.Error())
		os.Exit(1)
	}
	return &LaserPair{m}
}

func (l *LaserPair) Save() {
	b, _ := json.Marshal(l.m)
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	ioutil.WriteFile("./laser.json", out.Bytes(), 0640)
}

func (l *LaserPair) Record(key string, receiverID string, receiverIdx string) {
	info := ReceiverInfo{}
	info.ID = receiverID
	info.Idx = receiverIdx
	info.Valid = 1
	l.m[key] = info
}

package core

import (
	"strconv"
)

type ReceiverStatus int

const (
	ReceiverStatusUnknown           = 0
	ReceiverStatusBroken            = 1
	ReceiverStatusBrokenButReceived = 2
	ReceiverStatusNotReceived       = 3
	ReceiverStatusNormal            = 4
)

type QuickChecker struct {
	srv       *Srv
	receivers map[string]bool
	statusMap map[string]ReceiverStatus
}

func NewQuickChecker(srv *Srv) *QuickChecker {
	qc := QuickChecker{}
	qc.srv = srv
	qc.receivers = GetLaserPair().GetValidReceivers()
	qc.statusMap = make(map[string]ReceiverStatus)
	for k, _ := range qc.receivers {
		qc.statusMap[k] = ReceiverStatusUnknown
	}
	qc.openAllLasers()
	return &qc
}

func (qc *QuickChecker) OnArduinoHeartBeat(hb *InboxMessage) {
	ur := hb.GetStr("UR")
	id := hb.Address.ID
	for i, r := range ur {
		c := string(r)
		key := id + ":" + strconv.Itoa(i)
		if c == "1" {
			if qc.receivers[key] {
				qc.statusMap[key] = ReceiverStatusNormal
			} else {
				qc.statusMap[key] = ReceiverStatusBrokenButReceived
			}
		} else {
			if qc.receivers[key] {
				qc.statusMap[key] = ReceiverStatusNotReceived
			} else {
				qc.statusMap[key] = ReceiverStatusBroken
			}
		}
	}
}

func (qc *QuickChecker) Query() {
	msg := NewInboxMessage()
	msg.SetCmd("QuickCheck")
	data := make(map[string]ReceiverStatus)
	for k, v := range qc.statusMap {
		data[k] = v
	}
	msg.Set("data", data)
	qc.srv.sends(msg, InboxAddressTypeAdminDevice)
}

func (qc *QuickChecker) Record() {
	ret := make([]string, 0)
	for k, v := range qc.statusMap {
		if v == ReceiverStatusNotReceived {
			ret = append(ret, k)
		}
	}
	GetLaserPair().RecordBrokens(ret)
}

func (qc *QuickChecker) openAllLasers() {
	senders := GetLaserPair().GetValidSenders()
	for id, li := range senders {
		info := arduinoInfoFromID(id)
		msg := NewInboxMessage()
		msg.SetCmd("laser_ctrl")
		laserList := make([]map[string]string, len(li))
		for i, v := range li {
			v += 1
			if info.LaserNum == 5 {
				v += 5
			}
			laser := make(map[string]string)
			laser["laser_s"] = "1"
			laser["laser_n"] = strconv.Itoa(v)
			laserList[i] = laser
		}
		msg.Set("laser", laserList)
		qc.srv.sendToOne(msg, InboxAddress{InboxAddressTypeMainArduinoDevice, id})
	}
}

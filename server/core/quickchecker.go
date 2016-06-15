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
	hbCh    chan *InboxMessage
	qCh     chan struct{}
	closeCh chan chan struct{}
	closed  bool
	srv     *Srv
}

func NewQuickChecker(srv *Srv) *QuickChecker {
	qc := QuickChecker{}
	qc.srv = srv
	qc.hbCh = make(chan *InboxMessage)
	qc.qCh = make(chan struct{})
	qc.closeCh = make(chan chan struct{}, 1)
	qc.closed = false
	go qc.run()
	return &qc
}

func (qc *QuickChecker) Close(ret chan struct{}) {
	if qc.closed {
		return
	}
	qc.closeCh <- ret
	qc.closed = true
}

func (qc *QuickChecker) OnArduinoHeartBeat(msg *InboxMessage) {
	qc.hbCh <- msg
}

func (qc *QuickChecker) Query() {
	qc.qCh <- struct{}{}
}

func (qc *QuickChecker) run() {
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
		qc.srv.sendToOne(msg, InboxAddress{InboxAddressTypeMainArduinoDevice, id})
	}
	receivers := GetLaserPair().GetValidReceivers()
	statusMap := make(map[string]ReceiverStatus)
	for k, _ := range receivers {
		statusMap[k] = ReceiverStatusUnknown
	}
	for {
		select {
		case ch := <-qc.closeCh:
			ret := make([]string, 0)
			for k, v := range statusMap {
				if v == ReceiverStatusNotReceived {
					ret = append(ret, k)
				}
			}
			GetLaserPair().RecordBrokens(ret)
			ch <- struct{}{}
		case hb := <-qc.hbCh:
			ur := hb.GetStr("UR")
			id := hb.Address.ID
			for i, r := range ur {
				c := string(r)
				key := id + ":" + strconv.Itoa(i)
				if c == "1" {
					if receivers[key] {
						statusMap[key] = ReceiverStatusNormal
					} else {
						statusMap[key] = ReceiverStatusBrokenButReceived
					}
				} else {
					if receivers[key] {
						statusMap[key] = ReceiverStatusNotReceived
					} else {
						statusMap[key] = ReceiverStatusBroken
					}
				}
			}
		case <-qc.qCh:
			msg := NewInboxMessage()
			msg.SetCmd("QuickCheck")
			msg.Set("data", statusMap)
			qc.srv.sends(msg, InboxAddressTypeAdminDevice)
		}
	}
}

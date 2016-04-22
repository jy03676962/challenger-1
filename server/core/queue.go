package core

import (
	"container/list"
	"log"
	"sync"
)

const initCursor = 2000

var _ = log.Printf

var q = newQueue()

type Team struct {
	Count int `json:"count"`
	Num   int `json:"num"`
}

type Queue struct {
	li   *list.List
	cur  int
	lock *sync.RWMutex
}

func newQueue() *Queue {
	q := Queue{}
	q.li = list.New()
	q.cur = initCursor
	q.lock = new(sync.RWMutex)
	return &q
}

func AddTeam(count int) (*Team, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.cur += 1
	t := Team{Count: count, Num: q.cur}
	q.li.PushBack(&t)
	return &t, nil
}

func ResetQueue() error {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.li.Init()
	q.cur = initCursor
	return nil
}

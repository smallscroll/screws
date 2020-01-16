package screws

import (
	"log"
	"time"
)

//ILoopBomb ...
type ILoopBomb interface {
	Defuse()
}

//NewLoopBomb ...
func NewLoopBomb(function func(), clock [3]int, expiration time.Duration) ILoopBomb {
	lb := &loopBomb{
		Powder:     function,
		Clock:      clock,
		Defused:    make(chan bool),
		Expiration: expiration,
	}
	go lb.ignite()
	return lb
}

type loopBomb struct {
	Powder     func()
	Clock      [3]int
	Defused    chan bool
	Expiration time.Duration
}

func (lb *loopBomb) ignite() {
	lb.Powder()
	for {
		today := time.Now()
		tomorrow := today.Add(time.Hour * 24)
		tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), lb.Clock[0], lb.Clock[1], lb.Clock[2], 0, tomorrow.Location())
		select {
		case <-lb.Defused:
			return
		case <-time.After(lb.Expiration):
			return
		case <-time.After(tomorrow.Sub(today)):
			lb.Powder()
		}
	}
}

func (lb *loopBomb) Defuse() {
	lb.Defused <- true
	log.Println("LoopBomb  Defused")
}

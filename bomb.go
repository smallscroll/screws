package screws

import (
	"log"
	"time"
)

//IDailyBomb ...
type IDailyBomb interface {
	Defuse()
}

//NewDailyBomb ...
func NewDailyBomb(function func(), clock [3]int, expiration time.Duration) IDailyBomb {
	lb := &dailyBomb{
		Powder:     function,
		Clock:      clock,
		Defused:    make(chan bool),
		Expiration: expiration,
	}
	go lb.ignite()
	return lb
}

//dailyBomb ...
type dailyBomb struct {
	Powder     func()
	Clock      [3]int
	Defused    chan bool
	Expiration time.Duration
}

//ignite ...
func (d *dailyBomb) ignite() {
	d.Powder()
	for {
		today := time.Now()
		tomorrow := today.Add(time.Hour * 24)
		tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), d.Clock[0], d.Clock[1], d.Clock[2], 0, tomorrow.Location())
		select {
		case <-d.Defused:
			return
		case <-time.After(d.Expiration):
			return
		case <-time.After(tomorrow.Sub(today)):
			d.Powder()
		}
	}
}

//Defuse ...
func (d *dailyBomb) Defuse() {
	d.Defused <- true
	log.Println("LoopBomb  Defused")
}

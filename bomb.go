package screws

import (
	"log"
	"time"
)

//IDailyBomb ...
type IDailyBomb interface {
	Reset(clock [3]int)
	Defuse()
}

//NewDailyBomb ...
func NewDailyBomb(function func(d IDailyBomb), clock [3]int) IDailyBomb {
	lb := &dailyBomb{
		Powder:  function,
		Clock:   clock,
		Defused: make(chan bool),
	}
	go lb.ignite()
	return lb
}

//dailyBomb ...
type dailyBomb struct {
	Powder  func(d IDailyBomb)
	Clock   [3]int
	Defused chan bool
}

//ignite ...
func (d *dailyBomb) ignite() {
	d.Powder(d)
	for {
		today := time.Now()
		tomorrow := today.Add(time.Hour * 24)
		tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), d.Clock[0], d.Clock[1], d.Clock[2], 0, tomorrow.Location())
		select {
		case <-d.Defused:
			return
		case <-time.After(tomorrow.Sub(today)):
			d.Powder(d)
		}
	}
}

//Defuse ...
func (d *dailyBomb) Reset(clock [3]int) {
	d.Clock = clock
}

//Defuse ...
func (d *dailyBomb) Defuse() {
	d.Defused <- true
	log.Println("LoopBomb  Defused")
}

package screws

import (
	"fmt"
	"testing"
	"time"
)

func TestLoopBomb(t *testing.T) {
	b := NewDailyBomb(hi, [3]int{19, 50, 00})
	time.Sleep(time.Second * 600)
	b.Defuse()
	select {}
}

func hi(d IDailyBomb) {
	fmt.Println("hello~")
}

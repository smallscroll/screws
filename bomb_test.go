package screws

import (
	"fmt"
	"testing"
	"time"
)

func TestLoopBomb(t *testing.T) {
	b := NewLoopBomb(hi, [3]int{19, 50, 00}, time.Hour*1000000)
	time.Sleep(time.Second * 600)
	b.Defuse()
	select {}
}

func hi() {
	fmt.Println("hello~")
}

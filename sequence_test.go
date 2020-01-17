package screws

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSequence(t *testing.T) {
	testSequence := NewSequence()
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			fmt.Println(strconv.FormatInt(testSequence.Get(), 36))
			fmt.Println(testSequence.Get())
			wg.Done()
		}()
	}
	wg.Wait()

	for {
		for i := 0; i < 5; i++ {
			go func() {
				fmt.Println(testSequence.Get())
			}()
		}
		time.Sleep(time.Second * 60)
	}
}

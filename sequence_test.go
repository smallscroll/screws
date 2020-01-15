package screws

import (
	"fmt"
	"sync"
	"testing"
)

func TestSequence(t *testing.T) {
	testSequence := StartAndInitializeASequenceService(1)
	wg := sync.WaitGroup{}
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			//fmt.Println(strconv.FormatInt(s.Get(), 36))
			fmt.Println(testSequence.Get())
			wg.Done()
		}()
	}
	wg.Wait()

}

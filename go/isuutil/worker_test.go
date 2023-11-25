package isuutil

import (
	"sync"
	"testing"
	"time"
)

func TestWorker_Run(t *testing.T) {
	exampleWorker := NewWorker[int](1 * time.Second)
	wg := sync.WaitGroup{}
	go exampleWorker.Run(func(items []int) {
		t.Log(items)
		for i := 0; i < len(items); i++ {
			wg.Done()
		}
	})

	for i := 0; i < 10; i++ {
		wg.Add(1)
		exampleWorker.Send(i)
		time.Sleep(500 * time.Millisecond)
	}
	wg.Wait()
}

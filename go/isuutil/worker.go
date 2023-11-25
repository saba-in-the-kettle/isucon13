package isuutil

import (
	"time"
)

// Worker はgoroutineで動く非同期workerです。
// 処理を一定間隔で非同期に実行したいときに使う。
type Worker[T any] struct {
	ch       chan T
	interval time.Duration
}

func NewWorker[T any](interval time.Duration) *Worker[T] {
	return &Worker[T]{
		// sizeをめちゃくちゃでかくしといて、channelへの送信がブロックされないようにする
		// ISUCONなら良いが、業務ではあまりやらないほうが良い
		ch:       make(chan T, 100000),
		interval: interval,
	}
}

// Send はworkerにitemを送信します。
func (w *Worker[T]) Send(item T) {
	w.ch <- item
}

// Run はworkerを起動します。
// Run 関数はgoroutineで動くことが想定されています。
// main関数で一度実行すると良いでしょう。
func (w *Worker[T]) Run(fun func([]T)) {
	var items []T

	// この時間ごとに処理をする
	timer := time.Tick(w.interval)
	for {
		select {
		case <-timer:
			if len(items) == 0 {
				break
			}

			go func(items []T) {
				// ここで定期的に何かの処理をする
				// channelからの受信をブロックしても良いなら、goroutineで実行せずにそのまま実行してもよい
				fun(items)
			}(items)

			// 処理が終わったら、itemsを空にする
			items = []T{}
		case item := <-w.ch:
			items = append(items, item)
		}
	}
}

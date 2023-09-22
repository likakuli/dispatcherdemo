package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	dispatcher := newDispatcher()
	count := atomic.Int32{}
	dispatcher.addProcessor(func(obj interface{}) {
		count.Add(1)
		fmt.Println("processor1", obj)
	})
	dispatcher.addProcessor(func(obj interface{}) {
		count.Add(1)
		fmt.Println("processor2", obj)
	})

	stopCh := make(chan struct{})

	go func() {
		for i := 0; i < 100; i++ {
			dispatcher.dispatch(i)
		}
	}()

	go dispatcher.run(stopCh)

	time.Sleep(1 * time.Second)
	fmt.Printf("%d event processed\n", count.Load())
	close(stopCh)
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cnt := 0

	ch := make(chan struct{})

	var wg sync.WaitGroup

	multiplyDetector := func(num int, ch chan struct{}) {
		defer wg.Done()
		for i := 0; i <= 500; i++ {
			if i%num == 0 {
				ch <- struct{}{}
			}
		}
	}

	wg.Add(3)
	go multiplyDetector(3, ch)
	go multiplyDetector(5, ch)
	go multiplyDetector(7, ch)

	cntWg := sync.WaitGroup{}
	cntWg.Add(1)
	go func() {
		defer cntWg.Done()
		for range ch {
			time.Sleep(time.Millisecond)
			cnt++
		}

		fmt.Println("channel closed")
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	wg.Wait()
	cntWg.Wait()

	fmt.Println(cnt)
}

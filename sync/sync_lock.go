/*
 * @Author: zengzh
 * @Date: 2023-02-20 14:38:56
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-02-20 16:06:41
 */
package sync

import (
	"fmt"
	"sync"
	"time"
	"sync/atomic"
)

func recv(c chan int) {
	fmt.Println("ready")
	ret := <-c
	fmt.Println("recive success", ret)
}

func send(i int) {
	ch := make(chan int)
	go recv(ch)
	ch <- i
	fmt.Println("success")
}

func waitGroup() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			fmt.Println(num)
		}(i)
	}
	wg.Wait()
	fmt.Println("success")
}

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) GetValue() int {
	// c.mu.Lock()
	// defer c.mu.Unlock()
	return c.value
}

func SyncMutex() {
	// var wg = sync.WaitGroup{}
	// counter := &Counter{}
	counter := 0
	for i := 0; i < 1000; i++ {
		// wg.Add(1)
		go func() {
			counter++
			// defer wg.Done()
			// counter.Increment()
		}()
	}
	time.Sleep(time.Second)
	// wg.Wait()
	fmt.Println("Counter value", counter)
}

var x int64
var wg sync.WaitGroup
var l sync.Mutex

func Add(){
	x++ 
	wg.Done()
}

func mutexAdd(){
	l.Lock()
	x++
	l.Unlock()
	wg.Done()
}

func atomicAdd(){
	atomic.AddInt64(&x,1)
	wg.Done()
}

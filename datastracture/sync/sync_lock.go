/*
 * @Author: zengzh 
 * @Date: 2023-02-20 14:38:56 
 * @Last Modified by:   zengzh 
 * @Last Modified time: 2023-02-20 14:38:56 
 */
package sync

import (
	"fmt"
	"sync"
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

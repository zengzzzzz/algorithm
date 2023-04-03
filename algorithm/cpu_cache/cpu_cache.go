/*
 * @Author: zengzh
 * @Date: 2023-04-03 16:02:28
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-04-03 16:50:49
 */

// https://www.duguying.net/article/set-cpu-affinity-binding-for-golang-program
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const EXEC_COUNT = 100 * 1000 * 1000

type Bits struct {
	a           int
	placeholder [64]byte
	b           int
}

var bits Bits

func whichCPU(prefix string) int {
	curCPU := runtime.NumCPU()
	p := make([]int, curCPU)
	for i := 0; i < curCPU; i++ {
		if err := runtime.GOMAXPROCS(i + 1); err < 1{
			fmt.Println("warning: could not set CPU affinity, continuing...")
			return -1
		}
		p[i] = i
	}
	runtime.GOMAXPROCS(curCPU)
	fmt.Printf("[%s] this process %d is running processor : %d\n", prefix, os.Getpid(), p)

	return 0
}

func setCPU(cpuID int) int {
	if err := runtime.GOMAXPROCS(cpuID + 1); err < 1 {
		fmt.Println("warning: could not set CPU affinity, continuing...")
		return -1
	}

	return 0
}

func thdFunc1(wg *sync.WaitGroup) {
	defer wg.Done()
	setCPU(0)
	whichCPU("thread 1 start")
	beginTV := time.Now()

	for i := 0; i < EXEC_COUNT; i++ {
		bits.a += 1
		a := bits.a
	}

	endTV := time.Now()
	fmt.Printf("thd1 perf:[%v]us\n", endTV.Sub(beginTV).Microseconds())
	whichCPU("thread 1 end")
}

func thdFunc2(wg *sync.WaitGroup) {
	defer wg.Done()
	setCPU(1)
	whichCPU("thread 2 start")
	beginTV := time.Now()

	for i := 0; i < EXEC_COUNT; i++ {
		bits.b += 2
		b := bits.b
	}

	endTV := time.Now()
	fmt.Printf("thd2 perf:[%v]us\n", endTV.Sub(beginTV).Microseconds())
	whichCPU("thread 2 end")
}

func main() {
	curCPU := runtime.NumCPU()
	fmt.Printf("system has %d processor(s).\n", curCPU)
	bits = Bits{}
	setCPU(0)
	whichCPU("main thread")

	var wg sync.WaitGroup
	wg.Add(2)
	go thdFunc1(&wg)
	go thdFunc2(&wg)
	wg.Wait()
}

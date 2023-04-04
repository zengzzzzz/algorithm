/*
 * @Author: zengzh
 * @Date: 2023-04-03 16:02:28
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-04-03 16:50:49
 */

// https://www.duguying.net/article/set-cpu-affinity-binding-for-golang-program
package cpucache

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	execCount = 100 * 1000 * 1000
)

type Bits struct {
	a           int
	placeholder [64]byte
	b           int
}

func whichCPU(prefix string) {
	curCPU := make([]int, runtime.NumCPU())
	runtime.GOMAXPROCS(len(curCPU))

	for i := range curCPU {
		curCPU[i] = i
	}
	fmt.Printf("[%s] this process %d is running processor(s) : %v\n", prefix, os.Getpid(), curCPU)
}

func setCPU(cpuID int) {
	// var newMask unix.CPUSet
	// newMask.Set(cpuID)
	// if err := unix.SchedSetaffinity(0, &newMask); err != nil {
	// 	fmt.Printf("set cpu affinity failed, err:%v, cpuID:%d \r ", err, cpuID)
	// }
}

func thdFunc1(wg *sync.WaitGroup, bits *Bits) {
	defer wg.Done()
	setCPU(0)
	whichCPU("thread 1 start")
	begin := time.Now()

	for i := 0; i < execCount; i++ {
		bits.a += 1
		a := bits.a
		_ = a
	}

	end := time.Now()
	fmt.Printf("thd1 perf:[%d]us\n", end.Sub(begin)/time.Microsecond)
	whichCPU("thread 1 end")
}

func thdFunc2(wg *sync.WaitGroup, bits *Bits) {
	defer wg.Done()
	setCPU(1)
	whichCPU("thread 2 start")
	begin := time.Now()

	for i := 0; i < execCount; i++ {
		bits.b += 2
		b := bits.b
		_ = b
	}

	end := time.Now()
	fmt.Printf("thd2 perf:[%d]us\n", end.Sub(begin)/time.Microsecond)
	whichCPU("thread 2 end")
}

// func main() {
// 	fmt.Printf("system has %d processor(s).\n", runtime.NumCPU())
// 	bits := Bits{}
// 	whichCPU("main thread")

// 	var wg sync.WaitGroup
// 	wg.Add(2)
// 	go thdFunc1(&wg, &bits)
// 	go thdFunc2(&wg, &bits)

// 	wg.Wait()
// }

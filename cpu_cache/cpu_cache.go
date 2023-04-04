/*
 * @Author: zengzh
 * @Date: 2023-04-03 16:02:28
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-04-03 16:50:49
 */

// https://www.duguying.net/article/set-cpu-affinity-binding-for-golang-program
package cpu_cache

import (
	// "fmt"
	"runtime"
	"sync"
	// "time"
)

const (
	execCount = 100 * 1000 * 1000
)

type Bits struct {
	a int
	b int
}

type BitsWithCache struct {
	a int
	placeholder [64]byte
	b int
}

func SetCPUCore(num int) {
	runtime.GOMAXPROCS(num)
}

func thdFunc1(wg *sync.WaitGroup, bits *Bits) {
	defer wg.Done()
	// begin := time.Now()

	for i := 0; i < execCount; i++ {
		bits.a += 1
		a := bits.a
		_ = a
	}

	// end := time.Now()
	// fmt.Printf("thd1 perf:[%d]us\n", end.Sub(begin)/time.Microsecond)
}

func thdFunc2(wg *sync.WaitGroup, bits *BitsWithCache) {
	defer wg.Done()
	// begin := time.Now()

	for i := 0; i < execCount; i++ {
		bits.a += 1
		a := bits.a
		_ = a
	}

	// end := time.Now()
	// fmt.Printf("thd1 perf:[%d]us\n", end.Sub(begin)/time.Microsecond)
}

func thdFunc3(wg *sync.WaitGroup, bits *Bits) {
	defer wg.Done()
	// begin := time.Now()

	for i := 0; i < execCount; i++ {
		bits.b += 2
		b := bits.b
		_ = b
	}

	// end := time.Now()
	// fmt.Printf("thd2 perf:[%d]us\n", end.Sub(begin)/time.Microsecond)
}

func thdFunc4(wg *sync.WaitGroup, bits *BitsWithCache) {
	defer wg.Done()
	// begin := time.Now()

	for i := 0; i < execCount; i++ {
		bits.b += 2
		b := bits.b
		_ = b
	}

	// end := time.Now()
	// fmt.Printf("thd2 perf:[%d]us\n", end.Sub(begin)/time.Microsecond)
}
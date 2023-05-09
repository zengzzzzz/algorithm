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
	"reflect"
	"runtime"
	"sync"
	// "time"
)

const (
	execCount = 100 * 1000 * 1000
)

type Bits struct {
	A int
	B int
}

type BitsWithCache struct {
	A           int
	placeholder [64]byte
	B           int
}

func SetCPUCore(num int) {
	runtime.GOMAXPROCS(num)
}

func thdFunc1(wg *sync.WaitGroup, bits interface{}) {
	defer wg.Done()
	// begin := time.Now()

	v := reflect.ValueOf(bits).Elem().FieldByName("A")
	for i := 0; i < execCount; i++ {
		v.SetInt(v.Int() + 1)
		a := v.Int()
		_ = a
	}

	// end := time.Now()
	// fmt.Printf("thd1 perf:[%d]us\n", end.Sub(begin)/time.Microsecond)
}


func thdFunc2(wg *sync.WaitGroup, bits interface{}) {
	defer wg.Done()
	// begin := time.Now()

	v := reflect.ValueOf(bits).Elem().FieldByName("B")
	for i := 0; i < execCount; i++ {
		v.SetInt(v.Int() + 2)
		b := v.Int()
		_ = b
	}
	// end := time.Now()
	// fmt.Printf("thd2 perf:[%d]us\n", end.Sub(begin)/time.Microsecond)
}


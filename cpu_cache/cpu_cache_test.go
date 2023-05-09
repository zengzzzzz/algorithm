package cpu_cache

import (
	"sync"
	"testing"
)

func BenchmarkTestRunOnTwoCoreWithoutCache(b *testing.B) {
	SetCPUCore(2)
	bits := Bits{}
	var wg sync.WaitGroup
	wg.Add(2)
	go thdFunc1(&wg, &bits)
	go thdFunc2(&wg, &bits)
	wg.Wait()
}

func BenchmarkTestRunOnTwoCoreWithCache(b *testing.B) {
	SetCPUCore(2)
	bits := BitsWithCache{}
	var wg sync.WaitGroup
	wg.Add(2)
	go thdFunc1(&wg, &bits)
	go thdFunc2(&wg, &bits)
	wg.Wait()
}

func BenchmarkTestRunOnOneCoreWithoutCache(b *testing.B) {
	SetCPUCore(1)
	bits := Bits{}
	var wg sync.WaitGroup
	wg.Add(2)
	go thdFunc1(&wg, &bits)
	go thdFunc2(&wg, &bits)
	wg.Wait()
}
func BenchmarkTestRunOnOneCoreWithCache(b *testing.B) {
	SetCPUCore(1)
	bits := BitsWithCache{}
	var wg sync.WaitGroup
	wg.Add(2)
	go thdFunc1(&wg, &bits)
	go thdFunc2(&wg, &bits)
	wg.Wait()
}

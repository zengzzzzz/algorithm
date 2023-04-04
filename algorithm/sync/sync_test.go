package sync

import (
	"testing"
)

func TestSyncChan(t *testing.T) {
	send(100)
}

func TestWaitGroup(t *testing.T) {
	waitGroup()
}

func TestSyncMutex(t *testing.T) {
	SyncMutex()
}

func BenchmarkAtomicAdd(t *testing.B) {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go atomicAdd()
	}
	wg.Wait()
}

func BenchmarkAdd(t *testing.B) {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go Add()
	}
	wg.Wait()
}

func BenchmarkMutexAdd(t *testing.B) {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go mutexAdd()
	}
	wg.Wait()
}

package cpu_cache

import (
    "fmt"
    "golang.org/x/sys/unix"
    "testing"
    "time"
)


func TestSchedSetaffinity(t *testing.T) {

    var newMask unix.CPUSet
    newMask.Set(0)

    err := unix.SchedSetaffinity(0, &newMask)
    if err != nil {
            fmt.Printf("SchedSetaffinity: %v", err)
    }

    for i := 0; i < 50; i++ {
            fmt.Println("Hello world")
            time.Sleep(2 * time.Second)
    }
}
package sync

import (
    "testing"
)

func TestSyncChan(t *testing.T){
    send(100)
}

func TestWaitGroup(t *testing.T){
    waitGroup()
}

func TestSyncMutex(t *testing.T){
    SyncMutex()
}
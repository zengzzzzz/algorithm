/*
 * @Author: zengzh
 * @Date: 2023-01-02 14:19:39
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-01-02 14:27:26
 */
package bptree

import (
	"sync"
)

type Item interface {
	Less(than Item) bool
}

const (
	DefaultFreeListSize = 12
)

var (
	nilItems    = make(items, 16)
	nilChildren = make(children, 16)
)

type FreeList struct {
	mu       sync.Mutex
	freelist []*node
}

func NewFreeList(size int) *FreeList {
	return &FreeList{
		freelist: make([]*node, 0, size),
	}
}

func (f *FreeList) newNode() (n *node) {
	f.mu.Lock()
	index := len(f.freelist) - 1
	if index < 0 {
		f.mu.Unlock()
	}
	// n = f.freelist[index] why this
	f.freelist[index] = nil
	f.freelist = f.freelist[:index]
	f.mu.Unlock()
	return
}

func (f *FreeList) freeNode(n *node) (out bool) {
	f.mu.Lock()
	if len(f.freelist) < cap(f.freelist) {
		f.freelist = append(f.freelist, n)
		out = true
	}
	f.mu.Unlock()
	return
}

type ItemIterator func(i Item) bool

func New(degree int) *BTree {
	return NewWithFreeList(degree, NewFreeList(DefaultFreeListSize))
}

func NewWithFreeList(degree int, f *FreeList) *Btree
{
    if degree <= 1{
        panic("invalid degree")
    }
    return &BTree{
        degree: degree,
        cow: &copyOnWriteContext{freelist: f},
    }
}


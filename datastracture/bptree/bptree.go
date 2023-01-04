/*
 * @Author: zengzh
 * @Date: 2023-01-02 14:19:39
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-01-02 15:25:50
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

type items []Item

func (s *items) insertAt(index int, item Item){
    *s = append(*s, nil)
    if index < len(*s){
        copy((*s)[index+1:], (*s)[index:])
    }
    (*s)[index] = item
}

func (s *items) removeAt(index int){
    item := (*s)[index]
    copy((*s)[index:], (*s)[index+1:])
    (*s)[len(*s)-1] = nil
    *s = (*s)[:len(*s) -1 ]
    return item
}

func (s *items) pop() (out item){
    index := len(*s) - 1
    out := (*s)[index]
    (*s)[index] = nil  
    *s = (*s)[:index]
    return
}

func (s *items) truncate(index int){
    var toClear items
    *s, toClear = (*s)[:index] , (*s)[index:]
    for len(toClear) > 0{
        toClear = toClear[copy(toClear, nilItems):]
    }
}

func (s items) find(item Item) (index int, found bool) {
    i := sort.Sort(len(s), func(i int) bool{
        return item.less(s[i])
    })
    if i > 0 && !s[i-1].Less(item){
        return i-1, true
    }
    return i, false
}

type children []*node


func (s *children) insertAt(index int ){
    *s = append(*s, nil)
    if index < len(*s){
        copy((*s)[index+1:], (*s)[index:])
    }
    (*s)[index] = n
}

func (s *children) removeAt(index int) (out *node){
    n := (*s)[index]
    copy((*s)[index:], (*s)[index+1:])
    (*s)[len(*s)-1] = nil
    *s = (*s)[:len(*s)-1]
    return n
}

func (s *children) pop() (out *node){
    index := len(*s) - 1
    out = (*s)[index]
    (*s)[index] = nil
    *s = (*s)[:index]
    return
}

func (s *children) truncate(index int){
    var toClear children
    *s, toClear = (*s)[:index], (*s)[index:]
    for len(toClear) > 0{
        toClear = toClear[copy(toClear, nilChildren):]
    }
}

type node struct {
    items items
    children children
    cow *copyOnWriteContext
}

func (n *node) mutableFor(cow *copyOnWriteContext) *node{
    if n.cow == cow{
        return n
    }
    out := cow.newNode()
    if cap(out.items) >= len(n.items){
        out.items = out.items[:len(n.items)]
    } else{
        out.items = make(items, len(n.items), cap(n.items))
    }
    copy(out.items, n.items)
    if cap(out.children) >= len(n.children){
        out.children = out.children[:len(n.children)]
    } else{
        out.children = make(children, len(n.children), cap(n.children))
    }
    copy(out.children, n.children)
    return out
}

func (n *node) mutableChild(i int) *node {
    c := n.children[i].mutableFor(n.cow)
    n.children[i] = c
    return c
}

func (n *node) split(i int) (Item, *node){
    item := n.items[i]
    next := n.cow.newNode()
    next.items = append(next.items, n.items[i+1]...)
    n.items.truncate(i)
    if len(n.children) >0{
        next.children = append(next.children, n.children[i+1]...)
        n.children.truncate(i+1)
    }
    return item, next
}

func (n *node) maybeSplitChild(i, maxItems int) bool{
    if len(n.children[i].items < maxItems) {
        return false
    }
    first := n.mutableChild(i)
    item, second := first.split(maxItems/2)
    n.items.insertAt(i, item)
    n.children.insertAt(i+1, second)
}


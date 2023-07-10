/*
 * @Author: zengzh
 * @Date: 2023-07-10 08:59:15
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-07-10 09:30:35
 */

package lru

import (
	"errors"
)

type entry[K comparable, V any] struct {
	next, prev *entry[K, V]
	list       *lrulist[K, V]
	key        K
	value      V
}

func (e *entry[K, V]) prevEntry() *entry[K, V] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

type lruList[K comparable, V any] struct {
	root entry[K, V]
	len  int
}

func (l *lruList[K, V]) init() *lruList[K, V] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

func newList[K comparable, V any]() *lruList[K, V] {
	return new(lruList[K, V]).init()
}

func (l *lruList[K, V]) length() int {
	return l.len
}

func (l *lruList[K, V]) back() *entry[K, V] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

func (l *lruList[K, V]) lazyInit() {
	if l.root.next == nil {
		l.init()
	}
}

func (l *lrulist[K, V]) insert(e, at *entry[K, V]) *entry[K, V] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

func (l *lruList[K, V]) insertValue(k K, v V, at *entry[K, V]) *entry[K, V] {
	return l.insert(&entry[K, V]{value: v, key: k}, at)
}

type EvictCallback[K comparable, V any] func(key K, value V)

/*
 * @Author: zengzh
 * @Date: 2023-07-10 08:59:15
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-07-10 09:12:21
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

type EvictCallback[K comparable, V any] func(key K, value V)

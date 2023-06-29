/*
 * @Author: zengzh
 * @Date: 2023-06-26 09:14:55
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-06-28 09:23:19
 */
package radix_tree

import (
	"sort"
	"strings"
)

type WalkFn func(key string, value interface{}) bool

type leafNode struct {
	key string
	val interface{}
}

type edge struct {
	label byte
	node  *node
}

type node struct {
	leaf   *leafNode
	prefix string
	edges  edges
}

func (n *node) isLeaf() bool {
	return n.leaf != nil
}

func (n *node) addEdge(e edge) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= e.label
	})
	n.edges = append(n.edges, edge{})
	copy(n.edges[idx+1:], n.edges[idx:])
	n.edges[idx] = e
}

func (n *node) updateEdge(label byte, node *node) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		n.edges[idx].node = node
		return
	}
	panic("updateEdge: edge not found")
}

func (n *node) getEdge(label byte) *node {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		return n.edges[idx].node
	}
	return nil
}

func (n *node) delEdge(label byte) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		copy(n.edges[idx:], n.edges[idx+1:])
		n.edges[len(n.edges)-1] = edge{}
		n.edges = n.edges[:len(n.edges)-1]
	}
}

type edges []edge

func (e edges) Len() int {
	return len(e)
}

func (e edges) Less(i, j int) bool {
	return e[i].label < e[j].label
}

func (e edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e edges) Sort() {
	sort.Sort(e)
}

type Tree struct {
	root *node
	size int
}

func New() *Tree {
	return NewFromMap(nil)
}

func NewFromMap(m map[string]interface{}) *Tree {
	t := &Tree{root: &node{}}
	for k, v := range m {
		t.Insert(k, v)
	}
	return t
}

func (t *Tree) Len() int {
	return t.size
}

func longestPrefix(k1, k2 string) int {
	max := len(k1)
	if l := len(k2); l < max {
		max = l
	}
	var i int
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

func (t *Tree) Insert(s string, v interface{}) (interface{}, bool) {
	var parent *node
	n := t.root
	search := s
	for {
		if len(search) == 0 {
			if n.isLeaf() {
				old := n.leaf.val
				n.leaf.val = v
				return old, true
			}
			n.leaf = &leafNode{search, v}
			t.size++
			return nil, false
		}
		parent = n
		n = n.getEdge(search[0])
		if n == nil {
			e := edge{
				label: search[0],
				node: &node{
					leaf:   &leafNode{s, v},
					prefix: search,
				},
			}
			parent.addEdge(e)
			t.size++
			return nil, false
		}
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			search = search[commonPrefix:]
			continue
		}
		t.size++
		child := &node{
			prefix: search[:commonPrefix],
		}
		parent.updateEdge(search[0], child)
		child.addEdge(edge{
			label: n.prefix[commonPrefix],
			node:  n,
		})
		n.prefix = n.prefix[commonPrefix:]
		leaf := &leafNode{
			key: s,
			val: v,
		}
		search = search[commonPrefix:]
		if len(search) == 0 {
			child.leaf = leaf
			return nil, false
		}
		child.addEdge(edge{
			label: search[0],
			node: &node{
				leaf:   leaf,
				prefix: search,
			},
		})
		return nil, false
	}
}

func (t *Tree) Delete(s string) (interface{}, bool) {
	var parent *node
	var label byte
	n := t.root
	search := s
	for {
		if len(search) == 0 {
			if !n.isLeaf() {
				break
			}
			goto DELETE
		}
		parent = n
		label = search[0]
		n = n.getEdge(label)
		if n == nil {
			break
		}
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	return nil, false

DELETE:
	leaf := n.leaf
	n.leaf = nil
	t.size--

	if parent != nil && len(n.edges) == 0 {
		parent.delEdge(label)
	}

	if n != t.root && len(n.edges) == 1 {
		n.mergeChild()
	}

	if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
		parent.mergeChild()
	}
	return leaf.val, true
}

func (t *Tree) DeletePrefix(s string) int {
	return t.deletePrefix(nil, t.root, s)
}

func (t *Tree) deletePrefix(parent, n *node, prefix string) int {
	if len(prefix) == 0 {
		subTreeSize := 0
		recursiveWalk(n, func(s string, v interface{}) bool {
			subTreeSize++
			return false
		})
		if n.isLeaf() {
			n.leaf = nil
		}
		n.edges = nil
		if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
			parent.mergeChild()
		}
		t.size -= subTreeSize
		return subTreeSize
	}
	label := prefix[0]
	child := n.getEdge(label)
	if child == nil || (!strings.HasPrefix(child.prefix, prefix) && !strings.HasPrefix(prefix, child.prefix)) {
		return 0
	}
	if len(child.prefix) > len(prefix) {
		prefix = prefix[len(prefix):]
	} else {
		prefix = prefix[len(child.prefix):]
	}
	return t.deletePrefix(n, child, prefix)
}

func (n *node) mergeChild() {
	e := n.edges[0]
	child := e.node
	n.prefix += child.prefix
	n.leaf = child.leaf
	n.edges = child.edges
}

func (t *Tree) Get(s string) (interface{}, bool) {
	n := t.root
	search := s
	for {
		if len(search) == 0 {
			if n.isLeaf() {
				return n.leaf.val, true
			}
			break
		}
		n = n.getEdge(search[0])
		if n == nil {
			break
		}
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	return nil, false
}

// not leaf node how to match
func (t *Tree) longestPrefix(s string) (string, interface{}, bool) {
	var last *leafNode
	n := t.root
	search := s
	for {
		if n.isLeaf() {
			last = n.leaf
		}
		if len(search) == 0 {
			break
		}
		n = n.getEdge(search[0])
		if n == nil {
			break
		}
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	if last != nil {
		return last.key, last.val, true
	}
	return "", nil, false
}

func (t *Tree) Minimum() (string, interface{}, bool) {
	n := t.root
	for {
		if n.isLeaf() {
			return n.leaf.key, n.leaf.val, true
		}
		if len(n.edges) == 0 {
			break
		}
		n = n.edges[0].node
	}
	return "", nil, false
}

func (t *Tree) Maximum() (string, interface{}, bool) {
	n := t.root
	for {
		if num := len(n.edges); num > 0 {
			n = n.edges[num-1].node
			continue
		}
		if n.isLeaf() {
			return n.leaf.key, n.leaf.val, true
		}
		break
	}
	return "", nil, false
}

func (t *Tree) Walk(fn WalkFn) {
	recursiveWalk(t.root, fn)
}

func (t *Tree) WalkPrefix(prefix string, fn WalkFn) {
	n := t.root
	search := prefix
	for {
		if len(search) == 0 {
			recursiveWalk(n, fn)
			return
		}
		n = n.getEdge(search[0])
		if n == nil {
			return
		}
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
			continue
		}
		if strings.HasPrefix(n.prefix, search) {
			recursiveWalk(n, fn)
		}
		return
	}
}

func (t *Tree) WalkPath(path string, fn WalkFn) {
	n := t.root
	search := path
	for {
		if n.leaf != nil && fn(n.leaf.key, n.leaf.val) {
			return
		}
		if len(search) == 0 {
			return
		}
		n = n.getEdge(search[0])
		if n == nil {
			return
		}
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
}
func recursiveWalk(n *node, fn WalkFn) bool {
	if n.leaf != nil && fn(n.leaf.key, n.leaf.val) {
		return true
	}
	i := 0
	k := len(n.edges)
	for i < k {
		e := n.edges[i]
		if recursiveWalk(e.node, fn) {
			return true
		}
		if len(n.edges) == 0 {
			return recursiveWalk(n, fn)
		}
		if len(n.edges) >= k {
			i++
		}
		k = len(n.edges)
	}
	return false
}

func (t *Tree) ToMap() map[string]interface{} {
	out := make(map[string]interface{}, t.size)
	t.Walk(func(k string, v interface{}) bool {
		out[k] = v
		return false
	})
	return out
}

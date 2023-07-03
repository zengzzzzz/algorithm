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

// WalkFn is the type of the function used visiting each item visited by Walk.
// Takes a key and value and returns a boolean if iteration should be terminated.
type WalkFn func(key string, value interface{}) bool

// leafNode is a leaf node in the tree
type leafNode struct {
	key string
	val interface{}
}

// edge is used to respresent a edge leading to a child node
type edge struct {
	// label is the first byte of the edge
	label byte
	node  *node
}

// node is a node in the tree
type node struct {
	// leaf is used to store possible leaf
	leaf *leafNode
	// prefix is the common prefix we ignore for edges
	prefix string
	// Edges should be stored in-order for iteration
	// to avoid a fully materialized slice to save memory,
	// since in most cases we expect to be sparse.
	edges edges
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

// Tree is a radix tree. This can be treated as a map[string]interface{}.
// The main advantage of this over a map is prefix-based lookups and oredered
// iteration.
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

// longestPrefix returns the length of the longest prefix shared by k1 and k2.
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

// Insert is used to insert or update a value in the tree. returns true if an
// existing value was updated.
func (t *Tree) Insert(s string, v interface{}) (interface{}, bool) {
	var parent *node
	n := t.root
	search := s
	for {
		// handle key exhaustion
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
		// look for the edge
		parent = n
		n = n.getEdge(search[0])

		// no edge found, create one
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
		// determine longest prefix of the search key on match
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			search = search[commonPrefix:]
			continue
		}
		// split the node
		t.size++
		child := &node{
			prefix: search[:commonPrefix],
		}
		parent.updateEdge(search[0], child)

		// restore the existing node
		child.addEdge(edge{
			label: n.prefix[commonPrefix],
			node:  n,
		})
		n.prefix = n.prefix[commonPrefix:]

		// create a new leaf node
		leaf := &leafNode{
			key: s,
			val: v,
		}

		// if the new key is a subset, add this to node
		search = search[commonPrefix:]
		if len(search) == 0 {
			child.leaf = leaf
			return nil, false
		}

		// create a new edge for the node
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

// Delete is used to delete a key, returning the previous value and if it was deleted.
func (t *Tree) Delete(s string) (interface{}, bool) {
	var parent *node
	var label byte
	n := t.root
	search := s
	for {
		// check for key exhaustion
		if len(search) == 0 {
			if !n.isLeaf() {
				break
			}
			goto DELETE
		}
		// look for an edge
		parent = n
		label = search[0]
		n = n.getEdge(label)
		if n == nil {
			break
		}
		// consume the search prefix
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	return nil, false

DELETE:
    // delete the leaf node
	leaf := n.leaf
	n.leaf = nil
	t.size--

	// check if we should delete this node from the parent
	if parent != nil && len(n.edges) == 0 {
		parent.delEdge(label)
	}

	// check if we should merge this node
	if n != t.root && len(n.edges) == 1 {
		n.mergeChild()
	}
	// check if we should merge the parent`s other child
	if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
		parent.mergeChild()
	}
	return leaf.val, true
}

// DeletePrefix is used to delete the subtree under a prefix
// returns how many nodes were deleted
// use this to delete large subtrees efficiently
func (t *Tree) DeletePrefix(s string) int {
	return t.deletePrefix(nil, t.root, s)
}


// delete does a recursive deletion
func (t *Tree) deletePrefix(parent, n *node, prefix string) int {
	// check for key exhaustion
	if len(prefix) == 0 {
		// remove the leaf node
		subTreeSize := 0
		// recursively walk from all edges of the node to be deleted
		recursiveWalk(n, func(s string, v interface{}) bool {
			subTreeSize++
			return false
		})
		if n.isLeaf() {
			n.leaf = nil
		}
		// delete the entire subtree
		n.edges = nil
		// check if we should merge the parent's other child
		if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
			parent.mergeChild()
		}
		t.size -= subTreeSize
		return subTreeSize
	}
	// look for an edge
	label := prefix[0]
	child := n.getEdge(label)
	if child == nil || (!strings.HasPrefix(child.prefix, prefix) && !strings.HasPrefix(prefix, child.prefix)) {
		return 0
	}
	// consume the search prefix
	if len(child.prefix) > len(prefix) {
		prefix = prefix[len(prefix):]
	} else {
		prefix = prefix[len(child.prefix):]
	}
	return t.deletePrefix(n, child, prefix)
}

// merge child node to current node
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

// LongesetPrefix is used to find the longest prefix of a key
func (t *Tree) LongestPrefix(s string) (string, interface{}, bool) {
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

// walkprefix is used to walk the tree under a prefix
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

// walkpath is used to walk the tree under a path, but only visiting nodes
// from the root down to a given leaf
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

// recursiveWalk is used to walk the tree recursively.
// returns true if the walk should be aborted.
func recursiveWalk(n *node, fn WalkFn) bool {
	if n.leaf != nil && fn(n.leaf.key, n.leaf.val) {
		return true
	}
	// Recurse on the children
	i := 0
	// keeps track of number of edges in previous iteration
	k := len(n.edges)
	for i < k {
		e := n.edges[i]
		if recursiveWalk(e.node, fn) {
			return true
		}
		// when we are iterating on the node, the node is possibility modified.
		// If there are no more edges, mergeChild happened, so the last edge 
		// became the current node n, on which we'll iterate one last time
		if len(n.edges) == 0 {
			return recursiveWalk(n, fn)
		}
		// if n.edges < k , it means that mergeChild happened, so we need to
		// iterate on the current node n
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

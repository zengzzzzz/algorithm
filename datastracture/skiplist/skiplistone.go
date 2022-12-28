/*
 * @Author: zengzh
 * @Date: 2022-12-28 15:38:33
 * @Last Modified by:   zengzh
 * @Last Modified time: 2022-12-28 15:38:33
 */
package skiplist

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	DefaultMaxlevel    int     = 10
	DefaultProbability float64 = 1 / math.E
)

type elementNode struct {
	next []*Element
}

type Element struct {
	elementNode
	key   float64
	value interface{}
}

func (e *Element) Key() float64 {
	return e.key
}

func (e *Element) Value() interface{} {
	return e.value
}

type SkipList struct {
	elementNode
	maxLevel      int
	length        int
	randSource    rand.Source
	probability   float64
	probTabale    []float64
	mutex         sync.RWMutex
	prevNodeCache []*elementNode
}

func NewSkipList() *SkipList {
	return NewWithMaxLevel(DefaultMaxlevel)
}

func ProbabilityTable(probability float64, maxlevel int) (table []float64) {
	for i := 1; i <= maxlevel; i++ {
		prob := math.Pow(probability, float64(i-1))
		table = append(table, prob)
	}
	return table
}

func NewWithMaxLevel(maxLevel int) *SkipList {
	if maxLevel < 1 || maxLevel > DefaultMaxlevel {
		panic("invalid maxlevel")
	}

	return &SkipList{
		elementNode:   elementNode{next: make([]*Element, maxLevel)},
		prevNodeCache: make([]*elementNode, maxLevel),
		maxLevel:      maxLevel,
		randSource:    rand.New(rand.NewSource((time.Now().UnixNano()))),
		probability:   DefaultProbability,
		probTabale:    ProbabilityTable(DefaultProbability, maxLevel),
	}
}

func (list *SkipList) randLevel() (level int) {
	r := float64(list.randSource.Int63()) / (1 << 63)
	level = 1
	for level < list.maxLevel && r < list.probTabale[level] {
		level++
	}
	return level
}

func (list *SkipList) SetProbability(newProbability float64) {
	list.probability = newProbability
	list.probTabale = ProbabilityTable(newProbability, list.maxLevel)
}

func (list *SkipList) Set(key float64, value interface{}) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var element *Element
	prevs := list.getPrevElementNodes(key)
	if element = prevs[0].next[0]; element != nil && key == element.key {
		element.value = value
		return element
	}

	element = &Element{
		elementNode: elementNode{next: make([]*Element, list.randLevel())},
		key:         key,
		value:       value,
	}
	list.length++

	for i := range element.next {
		element.next[i] = prevs[i].next[i]
		prevs[i].next[i] = element
	}
	return element
}

func (list *SkipList) getPrevElementNodes(key float64) []*elementNode {
	var prev *elementNode = &list.elementNode
	var next *Element
	prevs := list.prevNodeCache
	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && key > next.key {
			prev = &next.elementNode
			next = prev.next[i]
		}
		prevs[i] = prev
	}
	return prevs
}

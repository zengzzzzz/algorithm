package skiplist

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	defaultMaxlevel    int     = 10
	defaultProbability float64 = 1 / math.E
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
	elementNone
	maxLevel    int
	length      int
	randSource  rand.Source
	probability float64
	probTabale  []float64
    mutex sync.RWMutex
    prevNodeCache []*elementNode
}

func NewSkipList() *SkipList{
    return NewWithMaxLevel(defaultMaxlevel)
}


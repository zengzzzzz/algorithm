package bptree

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	seed := time.Now().Unix()
	fmt.Print(seed)
	rand.Seed(seed)
}

func perm(n int) (out []Item) {
	for _, v := range rand.Perm(n) {
		out = append(out, Int(v))
	}
}

func rang(n int) (out []Item) {
	for i := 0; i < n; i++ {
		out = append(out, Int(i))
	}
	return
}

func all(t *BTree) (out []Item) {
	t.Ascend(func(i Item) bool {
		out = append(out, i)
		return true
	}) . 
	return 
}

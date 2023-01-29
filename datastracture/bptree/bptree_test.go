/*
 * @Author: zengzh 
 * @Date: 2023-01-29 10:10:56 
 * @Last Modified by:   zengzh 
 * @Last Modified time: 2023-01-29 10:10:56 
 */
package bptree

import (
	"flag"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
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
	return
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
	})
	return
}

func rangrev(n int) (out []Item) {
	for i := n - 1; i >= 0; i-- {
		out = append(out, Int(i))
	}
	return
}

func allrev(t *BTree) (out []Item) {
	t.Descend(func(i Item) bool {
		out = append(out, i)
		return true
	})
	return
}

var btreeDegree = flag.Int("degree", 32, "B-Tree degree")

func TestBTree(t *testing.T) {
	tr := New(*btreeDegree)
	const treeSize = 1000
	for i := 0; i < 10; i++ {
		if min := tr.Min(); min != nil {
			t.Fatalf("min: got %v, want nil", min)
		}
		if max := tr.Max(); max != nil {
			t.Fatalf("max: got %v, want nil", max)
		}
		for _, item := range perm(treeSize) {
			if x := tr.ReplaceOrInsert(item); x != nil {
				t.Fatal("insert found item", x)
			}
		}
		for _, item := range perm(treeSize) {
			if !tr.Has(item) {
				t.Fatal("has didn't find item", item)
			}
		}
		for _, item := range perm(treeSize) {
			if x := tr.ReplaceOrInsert(item); x == nil {
				t.Fatal("insert didn't find item", item)
			}
		}
		if min, want := tr.Min(), Int(0); min != want {
			t.Fatalf("min: got %v, want %v", min, want)
		}
		if max, want := tr.Max(), Int(treeSize-1); max != want {
			t.Fatalf("max: got %v, want %v", max, want)
		}
		got := all(tr)
		want := rang(treeSize)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("all: got %v, want %v", got, want)
		}
		gotrev := allrev(tr)
		wantrev := rangrev(treeSize)
		if !reflect.DeepEqual(gotrev, wantrev) {
			t.Fatalf("mismatch: \n got: %v, \nwant: %v", got, want)
		}
		for _, item := range perm(treeSize) {
			if x := tr.Delete(item); x == nil {
				t.Fatalf("delete didn't find %v", item)
			}
		}
		if got = all(tr); len(got) != 0 {
			t.Fatalf("all: got %v, want empty", got)
		}
	}
	return
}

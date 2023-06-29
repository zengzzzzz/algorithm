/*
 * @Author: zengzh
 * @Date: 2023-06-28 09:18:13
 * @Last Modified by:   zengzh
 * @Last Modified time: 2023-06-28 09:18:13
 */

package radix_tree

import (
	crand "crypto/rand"
	"fmt"
	"testing"
)

func TestRadix(t *testing.T) {
	var min, max string
	inp := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		gen := generateUUID()
		inp[gen] = i
        if gen < min || i == 0 {
            min = gen
        }
        if gen > max || i == 0 {
            max = gen
        }
	}
	r := NewFromMap(inp)
	if r.Len() != len(inp) {
		t.Fatalf("expected %d, got %d", len(inp), r.Len())
	}
    r.Walk(func(k string, v interface{}) bool {
        println(k)
        return false
    })
    for k, v :=  range inp {
        out, ok := r.Get(k)
        if !ok {
            t.Fatalf("missing key %v", k)
        }
        if out != v {
            t.Fatalf("expected %v, got %v", v, out)
        }
    }
    outMint, _, _ := r.Minimum()
    if outMint != min {
        t.Fatalf("expected %v, got %v", min, outMint)
    }
    outMaxt, _, _ := r.Maximum()
    if outMaxt != max {
        t.Fatalf("expected %v, got %v", max, outMaxt)
    }
    for k, v := range inp {
        out, ok := r.Delete(k)
        if !ok {
            t.Fatalf("missing key %v", k)
        }
        if out != v {
            t.Fatalf("expected %v, got %v", v, out)
        }
    }
    if r.Len() != 0 {
        t.Fatalf("expected %d, got %d", 0, r.Len())
    }
}

func generateUUID() string {
	buf := make([]byte, 16)
	if _, err := crand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}

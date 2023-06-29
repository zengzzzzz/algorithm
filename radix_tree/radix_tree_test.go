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
    "reflect"
)

func TestRadix(t *testing.T) {
	var min, max string
	inp := make(map[string]interface{})
	for i := 0; i < 10; i++ {
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
	for k, v := range inp {
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

func TestRoot(t *testing.T) {
    r := New()
    _, ok := r.Delete("")
    if ok {
        t.Fatalf("bad")
    }
    _, ok = r.Insert("", true)
    if ok {
        t.Fatalf("bad")
    }
    val, ok := r.Get("")
    if !ok || val != true {
        t.Fatalf("bad %v", val)
    }
    val, ok = r.Delete("")
    if !ok || val != true {
        t.Fatalf("bad %v", val)
    }
}

func TestDeletePrefix(t *testing.T){
    type exp struct {
        inp []string
        prefix string
        out []string
        numDeleted int
    }
    cases := []exp{
        {[]string{"", "A", "AB", "ABC", "R", "S"}, "A", []string{"", "R", "S"}, 3},
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "ABC", []string{"", "A", "AB", "R", "S"}, 1},
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "", []string{}, 6},
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "S", []string{"", "A", "AB", "ABC", "R"}, 1},
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "SS", []string{"", "A", "AB", "ABC", "R", "S"}, 0},
    }
    for _, test := range cases {
        r := New()
        for _, ss := range test.inp {
            r.Insert(ss, true)
        }
        deleted := r.DeletePrefix(test.prefix)
        if deleted != test.numDeleted {
            t.Fatalf("Bad delete, excepted %v, got %v", test.numDeleted, deleted)
        }
        out :=  []string{}
        fn := func(s string, v interface{}) bool {
            out = append(out, s)
            return false
        }
        r.Walk(fn)
        if !reflect.DeepEqual(out, test.out) {
            t.Fatalf("Bad delete, excepted %v, got %v", test.out, out)
        }
    }
}

func TestLongestPrefix(t *testing.T){
    r := New()
    keys  :=  []string {
        "",
		"foo",
		"foobar",
		"foobarbaz",
		"foobarbazzip",
		"foozip",
    }
    for _, k := range keys {
        r.Insert(k, nil)
    }
    if r.Len() != len(keys) {
        t.Fatal("bad len: %v %v", r.Len(), len(keys))
    }
    type exp struct {
        inp string
        out string
    }
    cases := []exp{
        {"a", ""},
        {"abc", ""},
		{"fo", ""},
		{"foo", "foo"},
		{"foob", "foo"},
		{"foobar", "foobar"},
		{"foobarba", "foobar"},
		{"foobarbaz", "foobarbaz"},
		{"foobarbazzi", "foobarbaz"},
		{"foobarbazzip", "foobarbazzip"},
		{"foozi", "foo"},
		{"foozip", "foozip"},
		{"foozipzap", "foozip"},
    }
    for _, test := range cases {
        m, _, ok := r.LongestPrefix(test.inp)
        if !ok {
            t.Fatalf("not match: %v", test)
        }
        if m != test.out {
            t.Fatalf("mis match: %v", test)
        }
    }
}

func generateUUID() string {
	buf := make([]byte, 16)
	if _, err := crand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}
	va := fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
	fmt.Println(va[0])
	return va
}

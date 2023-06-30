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
	"reflect"
	"sort"
	"testing"
    "strconv"
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

func TestDeletePrefix(t *testing.T) {
	type exp struct {
		inp        []string
		prefix     string
		out        []string
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
		out := []string{}
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

func TestLongestPrefix(t *testing.T) {
	r := New()
	keys := []string{
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
		t.Fatalf("bad len: %v %v", r.Len(), len(keys))
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

func TestWalkPrefix(t *testing.T) {
	r := New()
	keys := []string{
		"foobar",
		"foo/bar/baz",
		"foo/baz/bar",
		"foo/zip/zap",
		"zipzap",
	}
	for _, k := range keys {
		r.Insert(k, nil)
	}
	if r.Len() != len(keys) {
		t.Fatalf("bad len: %v %v", r.Len(), len(keys))
	}
	type exp struct {
		inp string
		out []string
	}
	cases := []exp{
		{
			"f",
			[]string{"foobar", "foo/bar/baz", "foo/baz/bar", "foo/zip/zap"},
		},
		{
			"foo",
			[]string{"foobar", "foo/bar/baz", "foo/baz/bar", "foo/zip/zap"},
		},
		{
			"foob",
			[]string{"foobar"},
		},
		{
			"foo/",
			[]string{"foo/bar/baz", "foo/baz/bar", "foo/zip/zap"},
		},
		{
			"foo/b",
			[]string{"foo/bar/baz", "foo/baz/bar"},
		},
		{
			"foo/ba",
			[]string{"foo/bar/baz", "foo/baz/bar"},
		},
		{
			"foo/bar",
			[]string{"foo/bar/baz"},
		},
		{
			"foo/bar/baz",
			[]string{"foo/bar/baz"},
		},
		{
			"foo/bar/bazoo",
			[]string{},
		},
		{
			"z",
			[]string{"zipzap"},
		},
	}
	for _, test := range cases {
		out := []string{}
		fn := func(s string, v interface{}) bool {
			out = append(out, s)
			return false
		}
		r.WalkPrefix(test.inp, fn)
		sort.Strings(out)
		sort.Strings(test.out)
		if !reflect.DeepEqual(out, test.out) {
			t.Fatalf("mis match: %v %v", out, test.out)
		}
	}
}

func TestWalkPath(t *testing.T) {
	r := New()
	keys := []string{
		"foo",
		"foo/bar",
		"foo/bar/baz",
		"foo/baz/bar",
		"foo/zip/zap",
		"zipzap",
	}
	for _, k := range keys {
		r.Insert(k, nil)
	}
	if r.Len() != len(keys) {
		t.Fatalf("bad len: %v %v", r.Len(), len(keys))
	}
	type exp struct {
		inp string
		out []string
	}
	cases := []exp{
		{
			"f",
			[]string{},
		},
		{
			"foo",
			[]string{"foo"},
		},
		{
			"foo/",
			[]string{"foo"},
		},
		{
			"foo/ba",
			[]string{"foo"},
		},
		{
			"foo/bar",
			[]string{"foo", "foo/bar"},
		},
		{
			"foo/bar/baz",
			[]string{"foo", "foo/bar", "foo/bar/baz"},
		},
		{
			"foo/bar/bazoo",
			[]string{"foo", "foo/bar", "foo/bar/baz"},
		},
		{
			"z",
			[]string{},
		},
	}
	for _, test := range cases {
		out := []string{}
		fn := func(s string, v interface{}) bool {
			out = append(out, s)
			return false
		}
		r.WalkPath(test.inp, fn)
		sort.Strings(out)
		sort.Strings(test.out)
		if !reflect.DeepEqual(out, test.out) {
			t.Fatalf("mis match: %v %v", out, test.out)
		}
	}
}

func TestWalkDelete(t *testing.T) {
	r := New()
	r.Insert("init0/0", nil)
	r.Insert("init0/1", nil)
	r.Insert("init0/2", nil)
	r.Insert("init0/3", nil)
	r.Insert("init1/0", nil)
	r.Insert("init1/1", nil)
	r.Insert("init1/2", nil)
	r.Insert("init1/3", nil)
	r.Insert("init2", nil)
	deleteFn := func(s string, v interface{}) bool {
		r.Delete(s)
		return false
	}
	r.WalkPrefix("init1", deleteFn)
	for _, s := range []string{"init0/0", "init0/1", "init0/2", "init0/3", "init2"} {
		if _, ok := r.Get(s); !ok {
			t.Fatalf("missing key: %v", s)
		}
	}
	if n := r.Len(); n != 5 {
		t.Fatalf("bad len: %v %v", n, r.ToMap())
	}
	r.Walk(deleteFn)
	if n := r.Len(); n != 0 {
		t.Fatalf("bad len: %v %v", n, r.ToMap())
	}
}

func BenchmarkInsert(b *testing.B) {
	r := New()
	for i := 0; i < 10000; i++ {
		r.Insert(fmt.Sprintf("init%d", i), true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, updated := r.Insert(strconv.Itoa(n), true)
		if updated {
			b.Fatal("updated")
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

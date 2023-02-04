/*
 * @Author: zengzh
 * @Date: 2022-12-30 07:53:34
 * @Last Modified by: zengzh
 * @Last Modified time: 2022-12-30 07:53:55
 */
package skiplist

import (
	"testing"
)

func TestSkipList(t *testing.T) {
	lis := NewSkipList()
	lis.Set(1, 1)
	lis.Set(2, 2)
	lis.Set(2, 3)
	b := lis.Get(1)
	if b.Value() != 1 {
		t.Error("error")
	}
}


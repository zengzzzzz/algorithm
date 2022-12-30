/*
 * @Author: zengzh
 * @Date: 2022-12-29 13:12:27
 * @Last Modified by: zengzh
 * @Last Modified time: 2022-12-29 14:27:44
 */
package main

import (
	"datastracture/skiplist"
	"fmt"
)

func main() {
	lis := skiplist.NewSkipList()
	lis.Set(1, 1)
	lis.Set(2, 2)
	lis.Set(2, 3)
	b := lis.Get(1)
	fmt.Println(b.Value())
}

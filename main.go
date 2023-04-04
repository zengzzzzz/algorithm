/*
 * @Author: zengzh
 * @Date: 2022-12-29 13:12:27
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-01-02 16:14:10
 */
package main

import (
	"algorithm/cpu_cache"
	"fmt"
)

type Item interface {
	Less(than Item) bool
}

type Int int

func (a Int) Less(than Item) bool {
	return a < than.(Int)
}

func main() {
	cpu_cache.GetCPU("main")
	c := Int(1)
	fmt.Print(c.Less(Int(3)))
}

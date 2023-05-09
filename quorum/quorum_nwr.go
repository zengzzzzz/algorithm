/*
 * @Author: zengzh
 * @Date: 2023-05-09 12:31:31
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-05-09 12:44:20
 */
package quorum

import (
	"fmt"
	"sync"
)

type Data struct {
	Value string
	Seq   int
}

type Node struct {
	ID       int
	DataList []*Data
}

func (n *Node) Write(data *Data) {
	n.DataList = append(n.DataList, data)
	fmt.Printf("Node %d write data %s\n", n.ID, data.Value)
}

func (n *Node) Read(seq int) *Data {
	for _, data := range n.DataList {
		if data.Seq == seq {
			fmt.Printf("Node %d read data %s\n", n.ID, data.Value)
			return data
		}
	}
	return nil
}

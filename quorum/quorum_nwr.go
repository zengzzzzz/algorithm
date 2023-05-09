/*
 * @Author: zengzh
 * @Date: 2023-05-09 12:31:31
 * @Last Modified by: zengzh
 * @Last Modified time: 2023-05-09 13:48:15
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

func (n *Node) Read(Value string) *Data {
	for _, data := range n.DataList {
		if data.Value == Value {
			fmt.Printf("Node %d read data %s\n", n.ID, data.Value)
			return data
		}
	}
	return nil
}

func quorumNWWrite(nodes []*Node, data *Data, quorum int) bool {
	var wg sync.WaitGroup
	count := 0
	for _, node := range nodes[:quorum] {
		wg.Add(1)
		go func(node *Node) {
			defer wg.Done()
			node.Write(data)
			count++
		}(node)
	}
	wg.Wait()
	return count >= quorum
}

func quorumNWRead(nodes []*Node, value string, quorum int) *Data {
	var wg sync.WaitGroup
	datalist := make(chan *Data, len(nodes))
	for _, node := range nodes[:quorum] {
		wg.Add(1)
		go func(node *Node) {
			defer wg.Done()
			data := node.Read(value)
			if data != nil {
				datalist <- data
			}
		}(node)
	}
	wg.Wait()
	close(datalist)
	if len(datalist) < quorum {
		return nil
	}
	var lastestData *Data
	for data := range datalist {
		if lastestData == nil || data.Seq > lastestData.Seq {
			lastestData = data
		}
	}
	return lastestData
}

func quorumNWR(node []*Node, data *Data, writeQuorum, readQuorum int) bool {
	ok := quorumNWWrite(node, data, writeQuorum)
	if !ok {
		return false
	}
	readData := quorumNWRead(node, data.Value, readQuorum)
	if readData == nil || readData.Value != data.Value {
		return false
	}
	return true
}

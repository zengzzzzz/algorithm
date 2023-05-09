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

func (n *Node) Read(seq int) *Data {
	for _, data := range n.DataList {
		if data.Seq == seq {
			fmt.Printf("Node %d read data %s\n", n.ID, data.Value)
			return data
		}
	}
	return nil
}

package main

import (
	"fmt"
	"sync"
)

// 定义一个结构体来存储数据
type Data struct {
	Value string // 数据的值
	Seq   int    // 数据的序列号
}

// 定义一个结构体来表示节点
type Node struct {
	ID       int      // 节点的 ID
	DataList []*Data // 节点存储的数据列表
}

// 实现一个写操作
func (n *Node) Write(data *Data) {
	n.DataList = append(n.DataList, data)
	fmt.Printf("Node %d writes data: %v\n", n.ID, data)
}

// 实现一个读操作
func (n *Node) Read(seq int) *Data {
	for _, data := range n.DataList {
		if data.Seq == seq {
			fmt.Printf("Node %d reads data: %v\n", n.ID, data)
			return data
		}
	}
	return nil
}

// 实现 Quorum NW 算法的写操作
func quorumNWWrite(nodes []*Node, data *Data, quorum int) bool {
	var wg sync.WaitGroup
	count := 0

	// 遍历写组中的节点进行写操作
	for _, node := range nodes[:quorum] {
		wg.Add(1)
		go func(node *Node) {
			defer wg.Done()
			node.Write(data)
			count++
		}(node)
	}
	wg.Wait()

	// 如果写操作无法达到写组的阈值，则返回错误
	if count < quorum {
		return false
	}

	// 如果写操作成功，则返回成功
	return true
}

// 实现 Quorum NR 算法的读操作
func quorumNRRead(nodes []*Node, seq int, quorum int) *Data {
	var wg sync.WaitGroup
	dataList := make(chan *Data, len(nodes))

	// 遍历读组中的节点进行读操作
	for _, node := range nodes[:quorum] {
		wg.Add(1)
		go func(node *Node) {
			defer wg.Done()
			data := node.Read(seq)
			if data != nil {
				dataList <- data
			}
		}(node)
	}
	wg.Wait()

	// 如果读操作无法达到读组的阈值，则返回 nil
	if len(dataList) < quorum {
		return nil
	}

	// 如果读操作成功，则返回最新的数据
	var latestData *Data
	for data := range dataList {
		if latestData == nil || data.Seq > latestData.Seq {
			latestData = data
		}
	}
	return latestData
}

// 实现 Quorum NWR 算法
func quorumNWR(nodes []*Node, data *Data, writeQuorum, readQuorum int) bool {
	// 进行写操作
	ok := quorumNWWrite(nodes, data, writeQuorum)
	if !ok {
		return false
	}

	// 进行读操作
	readData := quorumNRRead(nodes, data.Seq, readQuorum)
	if readData == nil || readData.Value != data.Value {
		return false
	}

	// 如果写和读操作均成功，则返回成功
	return true
}

func main() {
	// 创建三个节点
	node1 := &Node{ID: 1}
	node2 := &Node{ID: 2}
	node3 := &Node{ID: 3}
	nodes := []*Node{node1, node2, node3}

	// 定义一个数据
	data := &Data{Value: "Hello, world!", Seq: 1}

	// 调用 Quorum NWR 算法进行写操作和读操作
	ok := quorumNWR(nodes, data, 2, 2)
	fmt.Printf("Quorum NWR result: %v\n", ok)
}
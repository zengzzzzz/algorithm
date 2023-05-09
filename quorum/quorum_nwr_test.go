/*
 * @Author: zengzh
 * @Date: 2023-05-09 12:31:33
 * @Last Modified by:   zengzh
 * @Last Modified time: 2023-05-09 12:31:33
 */
package quorum

import (
	"testing"
)

func TestNodeWrite(t *testing.T) {
	node := &Node{1, nil}
	data := &Data{"data1", 1}
	node.Write(data)
	if len(node.DataList) != 1 || node.DataList[0].Value != "data1" {
		t.Errorf("Node write data failed")
	}
}

func TestNodeRead(t *testing.T) {
	node := &Node{1, []*Data{{Value: "data1", Seq: 1}}}
	data := node.Read("data1")
	if data.Value != "data1" {
		t.Errorf("Node read data failed")
	}
	data = node.Read("deta2")
	if data != nil {
		t.Errorf("Node read data should fail")
	}
}

func TestQuorumNWWrite(t *testing.T) {
	nodes := []*Node{
		{1, nil},
		{2, nil},
		{3, nil},
	}
	data := &Data{"data1", 1}
	// quorum is 2, should succeed
	ok := quorumNWWrite(nodes, data, 2)
	if !ok {
		t.Errorf("Expected quorumNWWrite to succeed")
	}

	// quorum is 3, should fail
	ok = quorumNWWrite(nodes, data, 3)
	if !ok {
		t.Errorf("Expected quorumNWWrite to fail")
	}
}
func TestQuorumNWRead(t *testing.T) {
	nodes := []*Node{
		{1, []*Data{{Value: "data1", Seq: 1}}},
		{2, []*Data{{Value: "data1", Seq: 2}}},
		{3, []*Data{{Value: "data2", Seq: 2}}},
	}
	// quorum is 2, should read data with seq 1
	data := quorumNWRead(nodes, "data1", 2)
	if data.Value != "data1" {
		t.Errorf("Expected to read data1, actual is %s", data.Value)
	}

	// quorum is 3, should fail to read
	data = quorumNWRead(nodes, "data2", 3)
	if data != nil {
		t.Errorf("Expected read to fail")
	}
}

func TestQuorumNWR(t *testing.T) {
	nodes := []*Node{
		{1, nil},
		{2, nil},
		{3, nil},
	}
	data := &Data{"data1", 1}

	// Write quorum is 2, read quorum is 2, should succeed
	ok := quorumNWR(nodes, data, 2, 2)
	if !ok {
		t.Errorf("Expected quorumNWR to succeed")
	}

	// Write quorum is 2, read quorum is 3, should fail
	ok = quorumNWR(nodes, data, 2, 3)
	if ok {
		t.Errorf("Expected quorumNWR to fail")
	}
}

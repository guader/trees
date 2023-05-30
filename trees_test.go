package trees

import (
	"encoding/json"
	"fmt"
	"testing"
)

type node struct {
	Id       int64   `json:",omitempty"`
	ParentId int64   `json:",omitempty"`
	Children []*node `json:",omitempty"`

	Seq int `json:",omitempty"`
}

func (n *node) GetId() int64 {
	return n.Id
}

func (n *node) GetParentId() int64 {
	return n.ParentId
}

func (n *node) GetChildren() []*node {
	return n.Children
}

func (n *node) SetChildren(children []*node) {
	n.Children = children
}

func (n *node) PrintAsNode() {
	m := *n
	m.Children = nil
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func (n *node) PrintAsTree() {
	m := *n
	m.ParentId = 0
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

var (
	_ Walker[int64, *node] = (*node)(nil)
	_ Maker[int64, *node]  = (*node)(nil)
)

func makeTrees() []*node {
	return []*node{
		{Id: 1, Children: []*node{ // 1
			{Id: 11}, // 1-1
			{Id: 12, Children: []*node{ // 1-2
				{Id: 121}, // 1-2-1
				{Id: 122}, // 1-2-2
			}},
		}},
		{Id: 2, Children: []*node{ // 2
			{Id: 21}, // 2-1
			{Id: 22}, // 2-2
		}},
		{Id: 3}, // 3
	}
}

func makeNodes() []*node {
	return []*node{
		{Id: 1, ParentId: 0}, // 1
		{Id: 2, ParentId: 0}, // 2
		{Id: 3, ParentId: 0}, // 3

		{Id: 11, ParentId: 1}, // 1-1
		{Id: 12, ParentId: 1}, // 1-2

		{Id: 21, ParentId: 2}, // 2-1
		{Id: 22, ParentId: 2}, // 2-2

		{Id: 121, ParentId: 12}, // 1-2-1
		{Id: 122, ParentId: 12}, // 1-2-2
	}
}

func TestWalk(t *testing.T) {
	t.Log("walk trees, depth first")
	Walk[int64, *node](makeTrees(), 0, func(i int, tree *node, parentId int64) bool {
		tree.PrintAsNode()
		return true
	})
	/*
		{"Id":1}
		{"Id":11}
		{"Id":12}
		{"Id":121}
		{"Id":122}
		{"Id":2}
		{"Id":21}
		{"Id":22}
		{"Id":3}
	*/

	t.Log("stop walking by returning false from handler func")
	Walk[int64, *node](makeTrees(), 0, func(i int, tree *node, parentId int64) bool {
		tree.PrintAsNode()
		return tree.Id != 121
	})
	/*
		{"Id":1}
		{"Id":11}
		{"Id":12}
		{"Id":121}
	*/

	t.Log("sequence and parentId")
	Walk[int64, *node](makeTrees(), 0, func(i int, tree *node, parentId int64) bool {
		tree.Seq = i + 1
		tree.ParentId = parentId
		tree.PrintAsNode()
		return true
	})
	/*
		{"Id":1,"Seq":1}
		{"Id":11,"ParentId":1,"Seq":1}
		{"Id":12,"ParentId":1,"Seq":2}
		{"Id":121,"ParentId":12,"Seq":1}
		{"Id":122,"ParentId":12,"Seq":2}
		{"Id":2,"Seq":2}
		{"Id":21,"ParentId":2,"Seq":1}
		{"Id":22,"ParentId":2,"Seq":2}
		{"Id":3,"Seq":3}
	*/

	t.Log("cycle detection")
	node1 := &node{Id: 1}
	node2 := &node{Id: 2}
	node1.Children = append(node1.Children, node2)
	node2.Children = append(node2.Children, node1)
	handledIds := make(map[int64]struct{})
	Walk[int64, *node]([]*node{node1, node2}, 0, func(i int, tree *node, parentId int64) bool {
		_, ok := handledIds[tree.Id]
		if ok {
			fmt.Println("cycle detected with node id:", tree.Id, ", parent id:", parentId)
			return false
		}
		handledIds[tree.Id] = struct{}{}
		return true
	})
	/*
		cycle detected with node id: 1 , parent id: 2
	*/
}

func TestMake(t *testing.T) {
	t.Log("make trees from node list, a tree would be set to the top level when its parent not exist")
	for _, tree := range Make[int64, *node](makeNodes()) {
		tree.PrintAsTree()
	}
	/*
		{
		  "Id": 1,
		  "Children": [
		    {
		      "Id": 11,
		      "ParentId": 1
		    },
		    {
		      "Id": 12,
		      "ParentId": 1,
		      "Children": [
		        {
		          "Id": 121,
		          "ParentId": 12
		        },
		        {
		          "Id": 122,
		          "ParentId": 12
		        }
		      ]
		    }
		  ]
		}
		{
		  "Id": 2,
		  "Children": [
		    {
		      "Id": 21,
		      "ParentId": 2
		    },
		    {
		      "Id": 22,
		      "ParentId": 2
		    }
		  ]
		}
		{
		  "Id": 3
		}
	*/
}

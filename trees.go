package trees

type Walker[ID comparable, NODE any] interface {
	GetId() ID
	GetChildren() []NODE
}

type Maker[ID comparable, NODE any] interface {
	Walker[ID, NODE]
	GetParentId() ID
	SetChildren([]NODE)
}

func Walk[ID comparable, T Walker[ID, T]](trees []T, parentId ID, handle func(i int, tree T, parentId ID) bool) bool {
	for i, tree := range trees {
		if !handle(i, tree, parentId) {
			return false
		}
		if !Walk(tree.GetChildren(), tree.GetId(), handle) {
			return false
		}
	}
	return true
}

func Make[ID comparable, T Maker[ID, T]](nodes []T) []T {
	nodeById := make(map[ID]T, len(nodes))
	for _, node := range nodes {
		nodeById[node.GetId()] = node
	}
	var trees []T
	for _, node := range nodes {
		parent, ok := nodeById[node.GetParentId()]
		if ok {
			parent.SetChildren(append(parent.GetChildren(), node))
		} else {
			trees = append(trees, node)
		}
	}
	return trees
}

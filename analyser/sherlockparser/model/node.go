package model

type Node struct {
	tag      Tag
	parent *Node
	children []*Node
}

func (node *Node) Tag() Tag {
	return node.tag
}

func (node *Node) Children() []*Node {
	return node.children
}


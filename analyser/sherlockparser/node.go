package sherlockparser

/*
Node struct.
*/
type Node struct {
	tag      *Tag
	parent   *Node
	children []*Node
}

/*
Tag returns a pointer to the Tag of the Node.
*/
func (node *Node) Tag() *Tag {
	return node.tag
}

/*
Children returns a pointer to a slice of nodes/children of the Node.
*/
func (node *Node) Children() []*Node {
	return node.children
}

/*
Parent returns a pointer to the parent Node.
*/
func (node *Node) Parent() *Node {
	return node.parent
}

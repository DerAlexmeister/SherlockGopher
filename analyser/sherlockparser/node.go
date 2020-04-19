package sherlockparser

type Node struct {
	tag      *Tag
	parent *Node
	children []*Node
}

/*
Returns a pointer to the Tag of the Node.
 */
func (node *Node) Tag() *Tag {
	return node.tag
}

/*
Returns a pointer to a slice of nodes/children of the Node.
 */
func (node *Node) Children() []*Node {
	return node.children
}

/*
Returns a pointer to the parent Node.
*/
func (node *Node) Parent() *Node {
	return node.parent
}

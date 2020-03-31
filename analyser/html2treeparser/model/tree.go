package model

type HTMLTree struct {
	htmlRaw string
	rootNode *Node
}

type Node struct {

}

func NewHTMLTree(html string) *HTMLTree {
	return &HTMLTree{htmlRaw:html}
}

func (tree *HTMLTree) Parse() *Node {

	return tree.rootNode
}

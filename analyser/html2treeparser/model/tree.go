package model

type HTMLTree struct {
	htmlRaw string
	rootNode *Node
}

func NewHTMLTree(html string) *HTMLTree {
	return &HTMLTree{htmlRaw:html}
}

func (tree *HTMLTree) Parse() *Node {
	return tree.rootNode
}

func (tree *HTMLTree) RootNode() *Node {
	return tree.rootNode
}


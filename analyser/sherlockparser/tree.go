package sherlockparser

import (
	"fmt"
	"log"

	"github.com/golang-collections/collections/stack"
)

/*
HTMLTree struct.
*/
type HTMLTree struct {
	htmlRaw  string
	rootNode *Node
}

/*
RootNode returns the RootNode of the HTMLTree.
*/
func (tree *HTMLTree) RootNode() *Node {
	return tree.rootNode
}

/*
NewHTMLTree returns pointer to a new empty HTMLTree with set html string.
*/
func NewHTMLTree(html string) *HTMLTree {
	return &HTMLTree{htmlRaw: html}
}

/*
Parse returns a pointer to the Root-Node of the parsed HTMLTree. Set verbose to true for info about malformed HTML-Documents.
*/
func (tree *HTMLTree) Parse(verbose bool) *Node {
	tokenStream := tree.tokenize()
	stack := stack.New()
	isRoot := true
	tree.rootNode = &Node{
		tag:      &Tag{},
		parent:   nil,
		children: nil,
	}
	currentNode := tree.rootNode
	for i := 0; i < len(tokenStream); i++ {
		switch currentToken := tokenStream[i].(type) {
		case *TagToken:
			switch currentToken.Type() {
			case StartTag:
				if isRoot {
					currentNode.tag = &Tag{
						tagType:       currentToken.TagType(),
						tagAttributes: tree.extractAttributes(currentToken.RawContent()),
						tagContent:    "",
					}
					currentNode.parent = nil
					isRoot = false
				} else {
					parent := currentNode
					currentNode = &Node{
						tag:      &Tag{},
						parent:   parent,
						children: nil,
					}
					parent.children = append(parent.Children(), currentNode)

					currentNode.tag = &Tag{
						tagType:       currentToken.TagType(),
						tagAttributes: tree.extractAttributes(currentToken.RawContent()),
						tagContent:    "",
					}
					currentNode.parent = parent
				}
				stack.Push(currentNode)
			case EndTag:
				if stack.Len() == 0 {
					currentNode = tree.RootNode()
				} else if nextNode, ok := stack.Pop().(*Node); ok {
					switch {
					case nextNode.Tag().TagType() != currentToken.TagType():
						if verbose {
							fmt.Printf("Expected Closing Tag for CurrentNode %s, but got Closing Tag for CurrentToken: %s", currentNode.Tag().TagType(), currentToken.TagType())
						}
					case stack.Len() > 0:
						currentNode = nextNode.Parent()
					default:
						currentNode = tree.RootNode()

					}
				} else {
					log.Fatal("Something went wrong while popping from the stack.")
				}

			case SelfClosingTag:
				newNode := &Node{
					tag: &Tag{
						tagType:       currentToken.TagType(),
						tagAttributes: tree.extractAttributes(currentToken.RawContent()),
						tagContent:    "",
					},
					parent:   currentNode,
					children: nil,
				}
				currentNode.children = append(currentNode.Children(), newNode)
			default:
				log.Fatalf("Type mismatch! Expected Start, End, or SelfClosingTag, but was %s", string(currentToken.Type()))
			}

		case *TextToken:
			currentNode.tag.tagContent = currentNode.Tag().TagContent() + currentToken.RawContent()
		default:
			log.Fatal("Type mismatch! TagToken or TextToken, but was none")
		}
	}
	return tree.rootNode
}

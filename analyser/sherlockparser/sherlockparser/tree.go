package sherlockparser

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
	"log"
	"strings"
)

type HTMLTree struct {
	htmlRaw  string
	rootNode *Node
}

/*
Return the RootNode of the HTMLTree.
*/
func (tree *HTMLTree) RootNode() *Node {
	return tree.rootNode
}

/*
Returns pointer to a new empty HTMLTree with set html string.
*/
func NewHTMLTree(html string) *HTMLTree {
	return &HTMLTree{htmlRaw: html}
}

/*
Returns a pointer to the Root-Node of the parsed HTMLTree.
*/
func (tree *HTMLTree) Parse() *Node {
	tokenStream := tree.tokenize()
	stack := stack.New()
	isRoot := true
	tree.rootNode = &Node{
		tag:      &Tag{},
		parent:   nil,
		children: nil,
	}
	currentNode := tree.rootNode
	for i := 0; i < len(*tokenStream); i++ {
		switch currentToken := (*tokenStream)[i].(type) {
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
				if currentNode, ok := stack.Pop().(*Node); ok {
					if currentNode.Tag().TagType() != currentToken.TagType() {
						matching := false
						for !matching {
							if node, ok := stack.Pop().(*Node); ok {
								fmt.Printf("Malformed HTML-Document! Expected closing tag %s, but was %s \n", currentNode.Tag().TagType(), currentToken.TagType())
								currentNode = node
								matching = node.Tag().TagType() == currentToken.TagType()
							}
						}
					}
				} else {
					log.Fatal("Something went wrong while popping from the stack.")
				}
				if stack.Len() > 0 {
					if nextNode, ok := stack.Peek().(*Node); ok {
						currentNode = nextNode

					} else {
						log.Fatal() // TODO: Print error
					}
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
				log.Fatalf("Type mismatch! Expected Start, End, or SelfClosingTag, but was %s", currentToken.Type())
			}

		case *TextToken:
			currentNode.tag.tagContent = currentNode.Tag().TagContent() + currentToken.RawContent()
		default:
			log.Fatal("Type mismatch! TagToken or TextToken, but was none")
		}
	}
	return tree.rootNode
}

/*
Returns a pointer to a TagAttribute-Slice generated based on the input string.
*/
func (tree *HTMLTree) extractAttributes(tagContent string) []*TagAttribute { //TODO: Handle attributes like style="font-size: 1px" whitespace is a problem
	attributes := make([]*TagAttribute, 0)
	nonQuotedWhiteSpaces := make([]int, 0)

	for i := 0; i < len(tagContent); i++ {
		if tagContent[i] == '"' { //forward quoted string
			for l := i + 1; l < len(tagContent); l++ {
				if tagContent[l] == '"' {
					i = l
					break
				}
			}
		}
		if tagContent[i] == '\'' {
			for l := i + 1; l < len(tagContent); l++ {
				if tagContent[l] == '\'' {
					i = l
					break
				}
			}
		}
		if tagContent[i] == ' ' {
			nonQuotedWhiteSpaces = append(nonQuotedWhiteSpaces, i)
		}
	}

	for i, whitespaceIndex := range nonQuotedWhiteSpaces {
		if i == (len(nonQuotedWhiteSpaces) - 1) {
			attributes = append(attributes, tree.processAttribute(tagContent[whitespaceIndex+1:]))
		} else {
			attributes = append(attributes, tree.processAttribute(tagContent[whitespaceIndex+1:nonQuotedWhiteSpaces[i+1]]))
		}
	}
	return attributes
}

/*
Returns a pointer to a TagAttribute which may contain several attributes from a single tag based on the input string.
*/
func (tree *HTMLTree) processAttribute(attributeRaw string) *TagAttribute {
	attributeSplit := strings.Split(attributeRaw, "=") //TODO: Kann eventuell ein Problem werden wenn = in gequotetem String enthalten sein darf, abhilfe -> nach ersten = suchen
	attribute := &TagAttribute{}
	attribute.attributeType = attributeSplit[0]
	if len(attributeSplit) == 2 {
		attribute.value = attributeSplit[1] //TODO: eventuell " bzw. ' trimmen, mit johnny abklären
	}
	return attribute
}


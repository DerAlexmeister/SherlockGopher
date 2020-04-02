package model

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
	"strings"
)

type HTMLTree struct {
	htmlRaw  string
	rootNode *Node
}

func NewHTMLTree(html string) *HTMLTree {
	return &HTMLTree{htmlRaw: html}
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func (tree *HTMLTree) handleTag(tag string) *TagToken {
	var token *TagToken
	tagRaw := strings.TrimSpace(tag)
	voidElements := []string{"area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr"}
	tagType := strings.Split(tagRaw, " ")[0]

	if _, contained := find(voidElements, tagType); contained {
		token = &TagToken{
			tokenType:  SelfClosingTag,
			tagType:    tagType,
			rawContent: tagRaw,
		}
	} else if tagRaw[0] == '/' {
		token = &TagToken{
			tokenType:  EndTag,
			tagType:    tagType,
			rawContent: tagRaw,
		}
	} else {
		token = &TagToken{
			tokenType:  StartTag,
			tagType:    tagType,
			rawContent: tagRaw,
		}
	}

	return token
}

func (tree *HTMLTree) Tokenize() []HtmlToken {
	element := ""
	classified := []HtmlToken{}

	for i := 0; i < len(tree.htmlRaw); i++ {
		if tree.htmlRaw[i] == '<' {
			if element != "" {
				element = strings.TrimSpace(element)
				if element != "" {
					classified = append(classified, &TextToken{
						tokenType:  PlainText,
						rawContent: element,
					})
				}
			}
			element = ""

			tag := ""
			for l := i + 1; l < len(tree.htmlRaw); l++ {
				if tree.htmlRaw[l] == '>' {
					i = l + 1
					classified = append(classified, tree.handleTag(tag))
					break
				} else {
					tag = tag + string(tree.htmlRaw[l])
				}
			}
		}
		element = element + string(tree.htmlRaw[i])
	}
	return classified
}

func (tree *HTMLTree) Parse() *Node {
	return nil
}

func (tree *HTMLTree) ParseRuneByRune() *Node {
	stack := stack.New()
	var currentNode *Node
	currentTag := ""
	var inTag bool
	var inClosingTag bool

	for i, char := range tree.htmlRaw {
		if char == '<' {
			currentTag = ""
			if tree.htmlRaw[i+1] == '/' {
				inClosingTag = true
				inTag = false
			} else {
				inTag = true
				inClosingTag = false
				if currentNode == nil {
					currentNode = &Node{
						tag:      Tag{},
						parent:   nil,
						children: nil,
					}
					tree.rootNode = currentNode
				} else {
					newNode := &Node{
						tag:      Tag{},
						parent:   currentNode,
						children: nil,
					}
					currentNode.children = append(currentNode.children, newNode)
					currentNode = newNode
				}
			}
		}
		if inTag {
			if char == '>' {
				if tree.htmlRaw[i-1] == '/' {
					// check for selfclosing tag
				} else {
					currentNode.tag = Tag{
						tagType:       currentTag[1:],
						tagAttributes: nil,
						tagContent:    "",
					}
					stack.Push(currentNode)
				}
				inTag = false
			}
			currentTag = currentTag + string(char)
		} else if inClosingTag {
			if char == '>' {
				expected, ok := stack.Peek().(*Node)
				if !ok {
					fmt.Errorf("Element on top of stack %v is not of type *Pointer", expected)
				}
				currentTag = currentTag[2:]
				if expected.tag.tagType == currentTag {
					currentNode = expected.parent
					stack.Pop()
				} else {
					fmt.Errorf("Malformed HTML string, %s does not match expected closing tag %s", expected.tag.tagType, currentTag)
				}
			} else {
				currentTag = currentTag + string(char)
			}
		}

	}

	return tree.rootNode
}

func (tree *HTMLTree) RootNode() *Node {
	return tree.rootNode
}

package model

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
			tagType:    tagType[1:],
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
					i = l
					classified = append(classified, tree.handleTag(tag))
					break
				} else {
					tag = tag + string(tree.htmlRaw[l])
				}
			}
			//element = ""
		} else {
			toAdd := string(tree.htmlRaw[i])
			element = element + toAdd
		}
	}
	return classified
}

func (tree *HTMLTree) Parse() *Node {
	tokenStream := tree.Tokenize()
	stack := stack.New()
	isRoot := true
	tree.rootNode = &Node{
		tag:      Tag{},
		parent:   nil,
		children: nil,
	}
	currentNode := tree.rootNode
	for i := 0; i < len(tokenStream); i++ {
		switch currentToken := tokenStream[i].(type) {
		case *TagToken:
			switch currentToken.tokenType {
			case StartTag:
				if isRoot {
					currentNode.tag = Tag{
						tagType:       currentToken.tagType,
						tagAttributes: tree.ExtractAttributes(currentToken.rawContent),
						tagContent:    "",
					}
					currentNode.parent = nil
					isRoot = false
				} else {
					parent := currentNode
					currentNode = &Node{
						tag:      Tag{},
						parent:   parent,
						children: nil,
					}
					parent.children = append(parent.children, currentNode)

					currentNode.tag = Tag{
						tagType:       currentToken.tagType,
						tagAttributes: tree.ExtractAttributes(currentToken.rawContent),
						tagContent:    "",
					}
					currentNode.parent = parent
				}
				stack.Push(currentNode)
			case EndTag:
				if currentNode, ok := stack.Pop().(*Node); ok {
					if currentNode.tag.tagType != currentToken.tagType {
						log.Fatalf("Malformed HTML-Document! Expected closing tag %s, but was %s", currentNode.tag.tagType, currentToken.tagType)
					} else if !ok{
						log.Fatal() // TODO: Print error
					}
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
					tag: Tag{
						tagType:       currentToken.tagType,
						tagAttributes: tree.ExtractAttributes(currentToken.rawContent),
						tagContent:    "",
					},
					parent:   currentNode,
					children: nil,
				}
				currentNode.children = append(currentNode.children, newNode)
			default:
				log.Fatalf("Type mismatch! Expected Start, End, or SelfClosingTag, but was %s", currentToken.tokenType)
			}

		case *TextToken:
			currentNode.tag.tagContent = currentNode.tag.tagContent + currentToken.rawContent
		default:
			log.Fatal("Type mismatch! TagToken or TextToken, but was none")
		}
	}
	return tree.rootNode
}

func (tree *HTMLTree) ExtractAttributes(tagContent string) []TagAttribute { //TODO: Handle attributes like style="font-size: 1px" whitespace is a problem
	attributesRaw := strings.Split(tagContent, " ")
	attributes := make([]TagAttribute, 0)
	for _, attribute := range attributesRaw {
		if contained := strings.IndexRune(attribute, '='); contained != -1 {
			splitAttribute := strings.Split(attribute, "=")
			attributes = append(attributes, TagAttribute{
				attributeType: splitAttribute[0],
				value:         splitAttribute[1],
			})
		}
	}
	return attributes
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

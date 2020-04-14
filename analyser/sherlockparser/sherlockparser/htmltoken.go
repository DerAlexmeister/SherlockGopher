package sherlockparser

import "strings"

type HtmlToken interface {
	Type() TokenType
	RawContent() string
}

type TokenType int

/*
Used to enumerate TokenTypes of HTML-Elements.
*/
const (
	StartTag       TokenType = 0
	EndTag         TokenType = 1
	SelfClosingTag TokenType = 2
	PlainText      TokenType = 3
)

type TagToken struct {
	tokenType  TokenType
	tagType    string
	rawContent string
}

/*
Returns the Type of a TagToken.
*/
func (tgTk *TagToken) Type() TokenType {
	return tgTk.tokenType
}

/*
Returns the raw content of a TagToken.
*/
func (tgTk *TagToken) RawContent() string {
	return tgTk.rawContent
}


/*
Returns the TagType of a TagToken
 */
func (tgTk *TagToken) TagType() string {
	return tgTk.tagType
}

type TextToken struct {
	tokenType  TokenType
	rawContent string
}

/*
Returns the Type of a TextToken
 */
func (txTk *TextToken) Type() TokenType {
	return txTk.tokenType
}

/*
Returns the raw Content of a TextToken
 */
func (txTk *TextToken) RawContent() string {
	return txTk.rawContent
}

/*
Returns a pointer to a HtmlToken-slice which is generated based on the input html stored in the HTMLTree.
*/
func (tree *HTMLTree) tokenize() *[]HtmlToken {
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
			if tree.htmlRaw[i+1] == '!' {
				for l := i + 1; l < len(tree.htmlRaw); l++ {
					if tree.htmlRaw[l] == '>' {
						i = l
						break
					}
				}
			} else {
				for l := i + 1; l < len(tree.htmlRaw); l++ {
					if tree.htmlRaw[l] == '>' {
						i = l
						classified = append(classified, tree.handleTag(tag))
						break
					} else {
						tag = tag + string(tree.htmlRaw[l])
					}
				}
			}
		} else {
			toAdd := string(tree.htmlRaw[i])
			element = element + toAdd
		}
	}
	return &classified
}

/*
Returns a pointer to a TagToken extracted from the input string.
*/
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

/*
Finds a string in a string slice. Returns index of string (-1 if not found) and bool if string was found.
*/
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}


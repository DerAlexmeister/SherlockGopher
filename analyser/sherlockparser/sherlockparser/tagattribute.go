package sherlockparser

import "strings"

type TagAttribute struct {
	attributeType string
	value    string
}
/*
Returns the AttributeType of the attribute
*/
func (attr *TagAttribute) AttributeType() string {
	return attr.attributeType
}

/*
Returns the value of the attribute
*/
func (attr *TagAttribute) Value() string {
	return attr.value
}

/*
Returns a pointer to a TagAttribute-Slice generated based on the input string.
*/
func (tree *HTMLTree) extractAttributes(tagContent string) []*TagAttribute {
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

package sherlockparser

import "fmt"

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

/*
Adds a string to the rawContent variable.
*/
func (tgTk *TagToken) AddToRawContent(toAdd string) {
	tgTk.rawContent = tgTk.rawContent + toAdd
}

/*
Returns the string representation of the struct. Only used for testing.
*/
func (tgTk *TagToken) ToString() string {
	return fmt.Sprintf("Type:%d TagType:%s RawContent:%s", tgTk.Type(), tgTk.TagType(), tgTk.RawContent())
}

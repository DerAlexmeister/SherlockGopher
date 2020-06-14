package sherlockparser

import "fmt"

/*
TagToken struct.
*/
type TagToken struct {
	tokenType  TokenType
	tagType    string
	rawContent string
}

/*
Type returns the Type of a TagToken.
*/
func (tgTk *TagToken) Type() TokenType {
	return tgTk.tokenType
}

/*
RawContent returns the raw content of a TagToken.
*/
func (tgTk *TagToken) RawContent() string {
	return tgTk.rawContent
}

/*
TagType returns the TagType of a TagToken.
*/
func (tgTk *TagToken) TagType() string {
	return tgTk.tagType
}

/*
AddToRawContent adds a string to the rawContent variable.
*/
func (tgTk *TagToken) AddToRawContent(toAdd string) {
	tgTk.rawContent += toAdd
}

/*
ToString returns the string representation of the struct. Only used for testing.
*/
func (tgTk *TagToken) ToString() string {
	return fmt.Sprintf("Type:%d TagType:%s RawContent:%s", tgTk.Type(), tgTk.TagType(), tgTk.RawContent())
}

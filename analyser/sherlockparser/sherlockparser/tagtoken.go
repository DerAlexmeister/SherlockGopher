package sherlockparser

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

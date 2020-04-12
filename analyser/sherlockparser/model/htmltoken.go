package model

type HtmlToken interface {
	Type() TokenType
	RawContent() string
}

type TokenType int

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

func (tgTk *TagToken) Type() TokenType {
	return tgTk.tokenType
}

func (tgTk *TagToken) RawContent() string {
	return tgTk.rawContent
}

func (tgTk *TagToken) TagType() string {
	return tgTk.tagType
}

type TextToken struct {
	tokenType  TokenType
	rawContent string
}

func (txTk *TextToken) Type() TokenType {
	return txTk.tokenType
}

func (txTk *TextToken) RawContent() string {
	return txTk.rawContent
}

package sherlockparser

import "fmt"

/*
TextToken struct.
*/
type TextToken struct {
	tokenType  TokenType
	rawContent string
}

/*
Type returns the Type of a TextToken.
*/
func (txTk *TextToken) Type() TokenType {
	return txTk.tokenType
}

/*
RawContent returns the raw Content of a TextToken.
*/
func (txTk *TextToken) RawContent() string {
	return txTk.rawContent
}

/*
AddToRawContent adds a string to the rawContent variable.
*/
func (txTk *TextToken) AddToRawContent(toAdd string) {
	txTk.rawContent += toAdd
}

/*
ToString returns the string representation of the struct. Only used for testing.
*/
func (txTk *TextToken) ToString() string {
	return fmt.Sprintf("Type:%d RawContent:%s", txTk.Type(), txTk.RawContent())
}

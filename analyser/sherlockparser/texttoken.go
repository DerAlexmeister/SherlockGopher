package sherlockparser

import "fmt"

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
Adds a string to the rawContent variable.
*/
func (txTk *TextToken) AddToRawContent(toAdd string) {
	txTk.rawContent = txTk.rawContent + toAdd
}

/*
   Returns the string representation of the struct. Only used for testing.
*/
func (txTk *TextToken) ToString() string {
	return fmt.Sprintf("Type:%d RawContent:%s", txTk.Type(), txTk.RawContent())
}

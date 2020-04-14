package sherlockparser

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

func (txTk *TextToken) AddToRawContent(toAdd string){
	txTk.rawContent = txTk.rawContent + toAdd
}

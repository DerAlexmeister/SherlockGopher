package sherlockparser

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	tables := []struct {
		input  *HTMLTree
		wanted []HTMLToken
	}{
		{
			input: NewHTMLTree("<span>test</span>"),
			wanted: []HTMLToken{
				&TagToken{
					tokenType:  StartTag,
					tagType:    "span",
					rawContent: "span",
				},
				&TextToken{
					tokenType:  PlainText,
					rawContent: "test",
				},
				&TagToken{
					tokenType:  EndTag,
					tagType:    "span",
					rawContent: "/span",
				},
			},
		},
		{
			input: NewHTMLTree("<html><!-- Comment!><meta value=test>"),
			wanted: []HTMLToken{
				&TagToken{
					tokenType:  StartTag,
					tagType:    "html",
					rawContent: "html",
				},
				&TagToken{
					tokenType:  SelfClosingTag,
					tagType:    "meta",
					rawContent: "meta value=test",
				},
			},
		},
		{
			input: NewHTMLTree("<html    ><title>Ja test = <<> test< /title >"),
			wanted: []HTMLToken{
				&TagToken{
					tokenType:  0,
					tagType:    "html",
					rawContent: "html",
				},
				&TagToken{
					tokenType:  0,
					tagType:    "title",
					rawContent: "title",
				},
				&TextToken{
					tokenType:  PlainText,
					rawContent: "Ja test = <<> test",
				},
				&TagToken{
					tokenType:  EndTag,
					tagType:    "title",
					rawContent: "/title",
				},
			},
		},
		{
			input: NewHTMLTree("<html    ><title>test< /title >"),
			wanted: []HTMLToken{
				&TagToken{
					tokenType:  0,
					tagType:    "html",
					rawContent: "html",
				},
				&TagToken{
					tokenType:  0,
					tagType:    "title",
					rawContent: "title",
				},
				&TextToken{
					tokenType:  PlainText,
					rawContent: "test",
				},
				&TagToken{
					tokenType:  EndTag,
					tagType:    "title",
					rawContent: "/title",
				},
			},
		},
		{
			input: NewHTMLTree("<html    ><title>< /title >"),
			wanted: []HTMLToken{
				&TagToken{
					tokenType:  0,
					tagType:    "html",
					rawContent: "html",
				},
				&TagToken{
					tokenType:  0,
					tagType:    "title",
					rawContent: "title",
				},
				&TagToken{
					tokenType:  EndTag,
					tagType:    "title",
					rawContent: "/title",
				},
			},
		},
	}

	for _, elem := range tables {
		have := elem.input.tokenize()

		for i, elemHave := range have {
			if elem.wanted[i].ToString() != elemHave.ToString() {
				t.Errorf("Wanted token %v, but was %v", elem.wanted[i], elemHave)
			}
		}
	}
}

func TestHandleToken(t *testing.T) {
	sut := NewHTMLTree("")
	tables := []struct {
		input  string
		wanted HTMLToken
	}{
		{
			input: "html",
			wanted: &TagToken{
				tokenType:  StartTag,
				tagType:    "html",
				rawContent: "html",
			},
		},
		{
			input: "meta property=og:type content=product",
			wanted: &TagToken{
				tokenType:  SelfClosingTag,
				tagType:    "meta",
				rawContent: "meta property=og:type content=product",
			},
		},
		{
			input: "/html",
			wanted: &TagToken{
				tokenType:  EndTag,
				tagType:    "html",
				rawContent: "/html",
			},
		},
	}

	for _, elem := range tables {
		have := sut.handleTag(elem.input)

		if have.ToString() != elem.wanted.ToString() {
			t.Errorf("Wanted %s but was %s", elem.wanted.ToString(), have.ToString())
		}
	}
}

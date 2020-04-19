package sherlockparser

import (
	"testing"
)

func TestProcessAttribute(t *testing.T) {
	sut := NewHTMLTree("")

	tables := []struct {
		input  string
		wanted *TagAttribute
	}{
		{"disabled", &TagAttribute{
			attributeType: "disabled",
			value:         "",
		},
		},
		{
			"value=yes", &TagAttribute{
			attributeType: "value",
			value:         "yes",
		},
		},
		{
			"href=\"/mmix.html\"", &TagAttribute{
			attributeType: "href",
			value:         "/mmix.html",
		},
		},
		{
			"href='/impressum.html'", &TagAttribute{
			attributeType: "href",
			value:         "/impressum.html",
		},
		},
		{
			"class=\"position-relative js-header-wrapper \"", &TagAttribute{
			attributeType: "class",
			value:         "position-relative js-header-wrapper ",
		},
		},
		{
			"value=color=red", &TagAttribute{
			attributeType: "value",
			value:         "color=red",
		},
		},
		{
			"value=\"color=red\"", &TagAttribute{
			attributeType: "value",
			value:         "color=red",
		},
		},
	}

	for _, elem := range tables {
		attr := sut.processAttribute(elem.input)
		if attr.AttributeType() != elem.wanted.AttributeType() {
			t.Errorf("Expected AttributeType to be %v, but was %v", elem.wanted.AttributeType(), attr.AttributeType())
		}
		if attr.Value() != elem.wanted.Value() {
			t.Errorf("Expected Value to be %v, but was %v", elem.wanted.Value(), attr.Value())
		}
	}
}

func TestExtractAttributes(t *testing.T) {
	sut := NewHTMLTree("")

	tables := []struct {
		input  string
		wanted []*TagAttribute
	}{
		{"link disabled", []*TagAttribute{
			{
				attributeType: "disabled",
				value:         "",
			},
		},
		},
		{
			"link disabled value=yes", []*TagAttribute{
			{
				attributeType: "value",
				value:         "yes",
			},
			{
				attributeType: "disabled",
				value:         "",
			},
		},
		},
		{
			"link href=\"/mmix.html\" disabled     value=yes ", []*TagAttribute{
			{
				attributeType: "href",
				value:         "/mmix.html",
			},

			{
				attributeType: "value",
				value:         "yes",
			},
			{
				attributeType: "disabled",
				value:         "",
			},
		},
		},
		{
			"link href='/impressum.html' href=\"/mmix.html\" disabled     value=yes ", []*TagAttribute{
			{
				attributeType: "href",
				value:         "/impressum.html",
			},
			{
				attributeType: "href",
				value:         "/mmix.html",
			},

			{
				attributeType: "value",
				value:         "yes",
			},
			{
				attributeType: "disabled",
				value:         "",
			},
		},
		},
		{
			"link class=\"position-relative js-header-wrapper \" href='/impressum.html' href=\"/mmix.html\" disabled     value=yes ", []*TagAttribute{
			{
				attributeType: "class",
				value:         "position-relative js-header-wrapper ",
			},
			{
				attributeType: "href",
				value:         "/impressum.html",
			},
			{
				attributeType: "href",
				value:         "/mmix.html",
			},

			{
				attributeType: "value",
				value:         "yes",
			},
			{
				attributeType: "disabled",
				value:         "",
			},
		},
		},
		{
			"link value=color=red class=\"position-relative js-header-wrapper \" href='/impressum.html' href=\"/mmix.html\" disabled     value=yes ", []*TagAttribute{
			{
				attributeType: "value",
				value:         "color=red",
			},
			{
				attributeType: "class",
				value:         "position-relative js-header-wrapper ",
			},
			{
				attributeType: "href",
				value:         "/impressum.html",
			},
			{
				attributeType: "href",
				value:         "/mmix.html",
			},

			{
				attributeType: "value",
				value:         "yes",
			},
			{
				attributeType: "disabled",
				value:         "",
			},
		},
		},
		{
			"link value=\"color=red\" value=color=blue class=\"position-relative js-header-wrapper \" href='/impressum.html' href=\"/mmix.html\" disabled     value=yes ", []*TagAttribute{
			{
				attributeType: "value",
				value:         "color=red",
			},
			{
				attributeType: "value",
				value:         "color=blue",
			},
			{
				attributeType: "class",
				value:         "position-relative js-header-wrapper ",
			},
			{
				attributeType: "href",
				value:         "/impressum.html",
			},
			{
				attributeType: "href",
				value:         "/mmix.html",
			},

			{
				attributeType: "value",
				value:         "yes",
			},
			{
				attributeType: "disabled",
				value:         "",
			},
		},
		},
	}

	for _, elem := range tables {
		have := sut.extractAttributes(elem.input)

		for _, wantElem := range elem.wanted {
			contained, occurances := contains(wantElem, have)

			if !contained || occurances > 1 {
				t.Errorf("TagAttribute with Type %s and Value %s was not contained or occured more then once for input string %s", wantElem.AttributeType(), wantElem.Value(), elem.input)
			}
		}
	}
}

func contains(attribute *TagAttribute, tbl []*TagAttribute) (bool, int) {
	contained := false
	occurances := 0

	for _, elem := range tbl {
		if elem.Value() == attribute.Value() && elem.AttributeType() == attribute.AttributeType() {
			occurances++
			contained = true
		}
	}

	return contained, occurances
}

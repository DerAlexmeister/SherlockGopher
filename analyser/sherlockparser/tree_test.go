package sherlockparser

import (
	"reflect"
	"testing"
)

func TestParseTitleContent(t *testing.T) {
	input := NewHTMLTree("<html><head><title>Test</title></head></html>")

	rootNode := &Node{
		tag:      &Tag{
			tagType:       "html",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   nil,
		children: nil,
	}

	 firstChild := &Node{
		tag:      &Tag{
			tagType:       "head",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   rootNode,
		children: nil,
	}
	 rootNode.children = append(rootNode.Children(), firstChild)

	secondChild := &Node{
		tag:      &Tag{
			tagType:       "title",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "Test",
		},
		parent:   firstChild,
		children: nil,
	}
	firstChild.children = append(firstChild.Children(), secondChild)

	if !reflect.DeepEqual(rootNode, input.Parse(false)) {
		t.Fail()
	}
}

func TestParseMoreSelfClosing(t *testing.T) {
	input := NewHTMLTree("<html><head><meta style=\"background:black position:fixed\"><link href=/home/something.html></head></html>").Parse(false)

	rootNode := &Node{
		tag:      &Tag{
			tagType:       "html",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   nil,
		children: nil,
	}

	firstChild := &Node{
		tag:      &Tag{
			tagType:       "head",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   rootNode,
		children: nil,
	}
	rootNode.children = append(rootNode.Children(), firstChild)
	attributes := make([]*TagAttribute,0)
	attributes = append(attributes, &TagAttribute{
		attributeType: "style",
		value:         "background:black position:fixed",
	})
	secondChild := &Node{
		tag:      &Tag{
			tagType:       "meta",
			tagAttributes: attributes,
			tagContent:    "",
		},
		parent:   firstChild,
		children: nil,
	}
	firstChild.children = append(firstChild.Children(), secondChild)

	attributes2 := make([]*TagAttribute,0)
	attributes2 = append(attributes2, &TagAttribute{
		attributeType: "href",
		value:         "/home/something.html",
	})
	thirdChild := &Node{
		tag:      &Tag{
			tagType:       "link",
			tagAttributes: attributes2,
			tagContent:    "",
		},
		parent:   firstChild,
		children: nil,
	}
	firstChild.children = append(firstChild.Children(), thirdChild)

	if !reflect.DeepEqual(rootNode, input) {
		t.Fail()
	}
}


func TestParseAttributesSelfClosing(t *testing.T) {
	input := NewHTMLTree("<html><head><meta style=\"background:black position:fixed\"></head></html>")

	rootNode := &Node{
		tag:      &Tag{
			tagType:       "html",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   nil,
		children: nil,
	}

	firstChild := &Node{
		tag:      &Tag{
			tagType:       "head",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   rootNode,
		children: nil,
	}
	rootNode.children = append(rootNode.Children(), firstChild)
	attributes := make([]*TagAttribute,0)
	attributes = append(attributes, &TagAttribute{
		attributeType: "style",
		value:         "background:black position:fixed",
	})
	secondChild := &Node{
		tag:      &Tag{
			tagType:       "meta",
			tagAttributes: attributes,
			tagContent:    "",
		},
		parent:   firstChild,
		children: nil,
	}
	firstChild.children = append(firstChild.Children(), secondChild)

	if !reflect.DeepEqual(rootNode, input.Parse(false)) {
		t.Fail()
	}
}

func TestParseMissingClosingTag(t *testing.T) {
	input := NewHTMLTree("<html><body><p></body></html>")

	rootNode := &Node{
		tag:      &Tag{
			tagType:       "html",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   nil,
		children: nil,
	}

	firstChild := &Node{
		tag:      &Tag{
			tagType:       "body",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   rootNode,
		children: nil,
	}
	rootNode.children = append(rootNode.Children(), firstChild)

	secondChild := &Node{
		tag:      &Tag{
			tagType:       "p",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   firstChild,
		children: nil,
	}
	firstChild.children = append(firstChild.Children(), secondChild)

	if !reflect.DeepEqual(rootNode, input.Parse(false)) {
		t.Fail()
	}
}

func TestParseMissingClosingTagAfterOpening(t *testing.T) {
	input := NewHTMLTree("<html><body><p><div></div></body></html>").Parse(false)

	rootNode := &Node{
		tag:      &Tag{
			tagType:       "html",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   nil,
		children: nil,
	}

	firstChild := &Node{
		tag:      &Tag{
			tagType:       "body",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   rootNode,
		children: nil,
	}
	rootNode.children = append(rootNode.Children(), firstChild)

	secondChild := &Node{
		tag:      &Tag{
			tagType:       "p",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   firstChild,
		children: nil,
	}
	firstChild.children = append(firstChild.Children(), secondChild)
	thirdChild := &Node{
		tag:      &Tag{
			tagType:       "div",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   secondChild,
		children: nil,
	}
	secondChild.children = append(secondChild.Children(), thirdChild)
	if !reflect.DeepEqual(rootNode, input) {
		t.Fail()
	}
}

func TestParseMoreChildren(t *testing.T){
	input := NewHTMLTree("<html><head><title>Test</title><p>Test2</p></head></html>")

	rootNode := &Node{
		tag:      &Tag{
			tagType:       "html",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   nil,
		children: nil,
	}

	firstChild := &Node{
		tag:      &Tag{
			tagType:       "head",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "",
		},
		parent:   rootNode,
		children: nil,
	}
	rootNode.children = append(rootNode.Children(), firstChild)

	secondChild := &Node{
		tag:      &Tag{
			tagType:       "title",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "Test",
		},
		parent:   firstChild,
		children: nil,
	}
	firstChild.children = append(firstChild.Children(), secondChild)

	secondChildTwo := &Node{
		tag:      &Tag{
			tagType:       "p",
			tagAttributes: make([]*TagAttribute,0),
			tagContent:    "Test2",
		},
		parent:   firstChild,
		children: nil,
	}
	firstChild.children = append(firstChild.Children(), secondChildTwo)

	if !reflect.DeepEqual(rootNode, input.Parse(false)) {
		t.Fail()
	}
}

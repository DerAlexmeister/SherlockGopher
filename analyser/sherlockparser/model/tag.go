package model

type Tag struct {
	tagType string
	tagAttributes []*TagAttribute
	tagContent string
}

func (tag *Tag) TagType() string {
	return tag.tagType
}

func (tag *Tag) Attributes() []*TagAttribute {
	return tag.tagAttributes
}

func (tag *Tag) TagContent() string {
	return tag.tagContent
}


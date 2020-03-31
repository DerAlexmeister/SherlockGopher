package model

type Tag struct {
	tagType string
	tagAttributes []string
	tagContent string
}

func (tag *Tag) TagType() string {
	return tag.tagType
}

func (tag *Tag) Attributes() []string {
	return tag.tagAttributes
}

func (tag *Tag) TagContent() string {
	return tag.tagContent
}


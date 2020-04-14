package sherlockparser

type Tag struct {
	tagType string
	tagAttributes []*TagAttribute
	tagContent string
}

/*
Returns the type of the tag
*/
func (tag *Tag) TagType() string {
	return tag.tagType
}

/*
Returns a pointer to a slice of TagAttributes of the tag
*/
func (tag *Tag) Attributes() []*TagAttribute {
	return tag.tagAttributes
}

/*
Returns the content of the tag
*/
func (tag *Tag) TagContent() string {
	return tag.tagContent
}


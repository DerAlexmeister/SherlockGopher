package sherlockparser

/*
Tag struct.
*/
type Tag struct {
	tagType       string
	tagAttributes []*TagAttribute
	tagContent    string
}

/*
TagType returns the type of the tag.
*/
func (tag *Tag) TagType() string {
	return tag.tagType
}

/*
Attributes returns a pointer to a slice of TagAttributes of the tag.
*/
func (tag *Tag) Attributes() []*TagAttribute {
	return tag.tagAttributes
}

/*
TagContent returns the content of the tag.
*/
func (tag *Tag) TagContent() string {
	return tag.tagContent
}

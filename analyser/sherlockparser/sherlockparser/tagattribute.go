package sherlockparser

type TagAttribute struct {
	attributeType string
	value    string
}
/*
Returns the AttributeType of the attribute
*/
func (attr *TagAttribute) AttributeType() string {
	return attr.attributeType
}

/*
Returns the value of the attribute
*/
func (attr *TagAttribute) Value() string {
	return attr.value
}

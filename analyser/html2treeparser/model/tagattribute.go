package model

type TagAttribute struct {
	attributeType string
	value    string
}

func (attr *TagAttribute) AttributeType() string {
	return attr.attributeType
}

func (attr *TagAttribute) Value() string {
	return attr.value
}

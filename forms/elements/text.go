package elements

type Text struct {
	Element
}

func NewText() *Text {
	b := &Text{}
	b.options = make(map[string]interface{})
	b.attributes = make(map[string]interface{})
	b.labelAttributes = make(map[string]interface{})
	b.labelOptions = make(map[string]interface{})
	b.SetAttribute("type", "text")
	return b
}

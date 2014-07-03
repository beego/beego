package elements

type Button struct {
	Element
}

func NewButton() *Button {
	b := &Button{}
	b.options = make(map[string]interface{})
	b.attributes = make(map[string]interface{})
	b.labelAttributes = make(map[string]interface{})
	b.labelOptions = make(map[string]interface{})
	b.SetAttribute("type", "button")
	return b
}

package elements

type Textarea struct {
	Element
}

func NewTextarea() *Textarea {
	b := &Textarea{}
	b.options = make(map[string]interface{})
	b.attributes = make(map[string]interface{})
	b.labelAttributes = make(map[string]interface{})
	b.labelOptions = make(map[string]interface{})
	b.SetAttribute("type", "textarea")
	return b
}

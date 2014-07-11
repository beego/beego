package elements

type Submit struct {
	Element
}

func NewSubmit() *Submit {
	b := &Submit{}
	b.options = make(map[string]interface{})
	b.attributes = make(map[string]interface{})
	b.labelAttributes = make(map[string]interface{})
	b.labelOptions = make(map[string]interface{})
	b.SetAttribute("type", "submit")
	return b
}

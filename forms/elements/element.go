package elements

type LableInterface interface {
	SetLable(label string)
	GetLable() (label string)
	SetLabelAttributes(labelattrs map[string]interface{})
	GetLabelAttributes() (labelattrs map[string]interface{})
	SetLabelOptions(labelOptions map[string]interface{})
	GetLabelOptions() (labelOptions map[string]interface{})
	ClearLabelOptions()
	RemoveLabelOptions(keys []string)
	SetLabelOption(key string, val interface{})
	GetLabelOption(key string) (val interface{})
	RemoveLabelOption(key string)
	HasLabelOption(key string) bool
}

type ElementInterface interface {
	SetName(name string)
	GetName() (name string)
	SetOptions(options map[string]interface{})
	SetOption(key string, val interface{})
	GetOptions() (options map[string]interface{})
	GetOption(key string) (val interface{})
	SetAttribute(key string, val interface{})
	GetAttribute(key string) (val interface{})
	RemoveAttribute(key string)
	HasAttribute(key string) bool
	SetAttributes(attributes map[string]interface{})
	GetAttributes() (attributes map[string]interface{})
	RemoveAttributes(keys []string)
	ClearAttributes()
	SetValue(val interface{})
	GetValue() (val interface{})
	SetMessages(msg string)
	GetMessages() (msg string)
	SetValidator(v ValidatorInterface)
	Valid(val interface{}) bool
	GetMessage() string
	LableInterface
}

type ValidatorInterface interface {
	IsValid(val interface{}) bool
	GetMessages() (msg string)
}

type Element struct {
	options         map[string]interface{}
	attributes      map[string]interface{}
	labelAttributes map[string]interface{}
	labelOptions    map[string]interface{}
	value           interface{}
	lable           string
	messages        string
	validator       ValidatorInterface
}

func (e *Element) SetName(name string) {
	e.SetAttribute("name", name)
}

func (e *Element) GetName() (name string) {
	return e.GetAttribute("name").(string)
}

func (e *Element) SetOptions(options map[string]interface{}) {
	if val, ok := options["label"]; ok {
		e.SetLable(val.(string))
	}

	if val, ok := options["label_attributes"]; ok {
		e.SetLabelAttributes(val.(map[string]interface{}))
	}

	if val, ok := options["label_options"]; ok {
		e.SetLabelOptions(val.(map[string]interface{}))
	}
	e.options = options
}

func (e *Element) SetOption(key string, val interface{}) {
	e.options[key] = val
}

func (e *Element) GetOptions() (options map[string]interface{}) {
	return e.options
}

func (e *Element) GetOption(key string) (val interface{}) {
	if val, ok := e.options[key]; ok {
		return val
	}
	return nil
}

func (e *Element) SetAttribute(key string, val interface{}) {
	e.attributes[key] = val
}

func (e *Element) GetAttribute(key string) (val interface{}) {
	if val, ok := e.attributes[key]; ok {
		return val
	}
	return nil
}

func (e *Element) RemoveAttribute(key string) {
	delete(e.attributes, key)
}

func (e *Element) HasAttribute(key string) bool {
	if _, ok := e.attributes[key]; ok {
		return true
	}
	return false
}

func (e *Element) SetAttributes(attributes map[string]interface{}) {
	for key, val := range attributes {
		e.SetAttribute(key, val)
	}
}

func (e *Element) GetAttributes() (attributes map[string]interface{}) {
	return e.attributes
}

func (e *Element) RemoveAttributes(keys []string) {
	for _, key := range keys {
		e.RemoveAttribute(key)
	}
}

func (e *Element) ClearAttributes() {
	e.attributes = make(map[string]interface{})
}

func (e *Element) SetValue(val interface{}) {
	e.value = val
}

func (e *Element) GetValue() (val interface{}) {
	return e.value
}

func (e *Element) SetLable(label string) {
	e.lable = label
}

func (e *Element) GetLable() (label string) {
	return e.lable
}

func (e *Element) SetLabelAttributes(labelattrs map[string]interface{}) {
	e.labelAttributes = labelattrs
}

func (e *Element) GetLabelAttributes() (labelattrs map[string]interface{}) {
	return e.labelAttributes
}

func (e *Element) SetLabelOptions(labelOptions map[string]interface{}) {
	for key, val := range labelOptions {
		e.SetLabelOption(key, val)
	}
}

func (e *Element) GetLabelOptions() (labelOptions map[string]interface{}) {
	return e.labelOptions
}

func (e *Element) ClearLabelOptions() {
	e.labelOptions = make(map[string]interface{})
}

func (e *Element) RemoveLabelOptions(keys []string) {
	for _, key := range keys {
		e.RemoveLabelOption(key)
	}
}

func (e *Element) SetLabelOption(key string, val interface{}) {
	e.labelOptions[key] = val
}

func (e *Element) GetLabelOption(key string) (val interface{}) {
	if v, ok := e.labelOptions[key]; ok {
		return v
	}
	return nil
}

func (e *Element) RemoveLabelOption(key string) {
	delete(e.labelOptions, key)
}

func (e *Element) HasLabelOption(key string) bool {
	if _, ok := e.labelOptions[key]; ok {
		return true
	}
	return false
}

func (e *Element) SetMessages(msg string) {
	e.messages = msg
}

func (e *Element) GetMessages() (msg string) {
	return e.messages
}

func (e *Element) SetValidator(v ValidatorInterface) {
	e.validator = v
}

func (e *Element) Valid(val interface{}) bool {
	if e.validator == nil {
		return true
	}
	if e.validator.IsValid(val) {
		return true
	}
	return false
}

func (e *Element) GetMessage() string {
	if e.validator == nil {
		return ""
	}
	return e.validator.GetMessages()
}

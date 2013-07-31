package orm

import ()

type fieldError struct {
	name string
	err  error
}

func (f *fieldError) Name() string {
	return f.name
}

func (f *fieldError) Error() error {
	return f.err
}

func NewFieldError(name string, err error) IFieldError {
	return &fieldError{name, err}
}

// non cleaned field errors
type fieldErrors struct {
	errors    map[string]IFieldError
	errorList []IFieldError
}

func (fe *fieldErrors) Get(name string) IFieldError {
	return fe.errors[name]
}

func (fe *fieldErrors) Set(name string, value IFieldError) {
	fe.errors[name] = value
}

func (fe *fieldErrors) List() []IFieldError {
	return fe.errorList
}

func NewFieldErrors() IFieldErrors {
	return &fieldErrors{errors: make(map[string]IFieldError)}
}

type Manager struct {
	ins    Modeler
	inited bool
}

// func (m *Manager) init(model reflect.Value) {
// 	elm := model.Elem()
// 	for i := 0; i < elm.NumField(); i++ {
// 		field := elm.Field(i)
// 		if _, ok := field.Interface().(Fielder); ok && field.CanSet() {
// 			if field.Elem().Kind() != reflect.Struct {
// 				field.Set(reflect.New(field.Type().Elem()))
// 			}
// 		}
// 	}
// }

func (m *Manager) Init(model Modeler) Modeler {
	if m.inited {
		return m.ins
	}
	m.inited = true
	m.ins = model
	return model
}

func (m *Manager) IsInited() bool {
	return m.inited
}

func (m *Manager) Clean() IFieldErrors {
	return nil
}

func (m *Manager) CleanFields(name string) IFieldErrors {
	return nil
}

func (m *Manager) GetTableName() string {
	return getTableName(m.ins)
}

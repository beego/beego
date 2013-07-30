package orm

import ()

// non cleaned field errors
type FieldErrors map[string]error

func (fe FieldErrors) Get(name string) error {
	return fe[name]
}

func (fe FieldErrors) Set(name string, value error) {
	fe[name] = value
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

func (m *Manager) Clean() FieldErrors {
	return nil
}

func (m *Manager) CleanFields(name string) FieldErrors {
	return nil
}

func (m *Manager) GetTableName() string {
	return getTableName(m.ins)
}

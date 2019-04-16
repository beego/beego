package goyaml2

const (
	N_Map = iota
	N_List
	N_String
)

type Node interface {
	Type() int
}

type MapNode map[string]interface{}

type ListNode []interface{}

type StringNode string

func (m *MapNode) Type() int {
	return N_Map
}

func (l *ListNode) Type() int {
	return N_List
}

func (s *StringNode) Type() int {
	return N_String
}

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Interface struct {
	Id   int
	Name string

	Index1 string
	Index2 string

	Unique1 string
	Unique2 string
}

func (i *Interface) TableIndex() [][]string {
	return [][]string{{"index1"}, {"index2"}}
}

func (i *Interface) TableUnique() [][]string {
	return [][]string{{"unique1"}, {"unique2"}}
}

func (i *Interface) TableName() string {
	return "INTERFACE_"
}

func (i *Interface) TableEngine() string {
	return "innodb"
}

func TestDbBase_GetTables(t *testing.T) {
	registry := DefaultModelRegistry
	registry.Register("", true, &Interface{})
	mi, ok := registry.Get("INTERFACE_")
	assert.True(t, ok)
	assert.NotNil(t, mi)

	engine := GetTableEngine(mi.AddrField)
	assert.Equal(t, "innodb", engine)
	uniques := GetTableUnique(mi.AddrField)
	assert.Equal(t, [][]string{{"unique1"}, {"unique2"}}, uniques)
	indexes := GetTableIndex(mi.AddrField)
	assert.Equal(t, [][]string{{"index1"}, {"index2"}}, indexes)
}

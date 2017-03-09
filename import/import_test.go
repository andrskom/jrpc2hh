package jrpc2hh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImportMap(t *testing.T) {
	a := assert.New(t)
	iMap := NewImportMap()

	val, ok := (*iMap)["github.com/andrskom/jrpc2hh/models"]
	a.True(ok)
	a.Equal("jModel", *val)
	a.Len(*iMap, 1)
}

func TestRegister(t *testing.T) {
	a := assert.New(t)
	iMap := NewImportMap()
	eImport := "import"
	iMap.Register(eImport)
	val, ok := (*iMap)[eImport]
	a.True(ok)
	a.Nil(val)
	val, ok = (*iMap)["github.com/andrskom/jrpc2hh/models"]
	a.True(ok)
	a.Equal("jModel", *val)
	a.Len(*iMap, 2)
}

func TestGenerateAlias(t *testing.T) {
	a := assert.New(t)
	iMap := NewImportMap()
	iMap.Register("import")
	iMap.Register("import/import")
	iMap.Register("blah")
	iMap.GenerateAlias()

	val, ok := (*iMap)["import"]
	a.True(ok)
	a.Nil(val)

	val, ok = (*iMap)["import/import"]
	a.True(ok)
	a.Regexp("^import_[0-9]+$", *val)

	val, ok = (*iMap)["blah"]
	a.True(ok)
	a.Nil(val)
}

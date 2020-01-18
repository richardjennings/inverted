package index

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCompositeIndex(t *testing.T) {
	cidx, err := NewIndex(map[string]map[string]string{"field": {"type": "text"}})
	assert.Nil(t, err)

	_, err = cidx.GetFieldIdx("field")
	assert.Nil(t, err)

	// index dupliate uri errors
	err = cidx.Index("1", map[string]interface{}{"field": "some content"})
	assert.Nil(t, err)
	err = cidx.Index("1", map[string]interface{}{"field": "some content"})
	assert.Equal(t, errors.New("document uri already exists"), err)

	assert.Equal(t, &Stats{DocumentCount: 1, Fields: map[string]IdxStats{"field": {TermCount: 2}}}, cidx.Stats())

	// get field index that does not exist returns error
	_, err = cidx.GetFieldIdx("notexists")
	assert.Equal(t, errors.New("field not found"), err)

	// create a field index that already exists
	_, err = cidx.newFieldIndex("field", NewTextIndex())
	assert.Equal(t, errors.New("field index already exists"), err)

	// try to index content with non existent field
	err = cidx.Index("a", map[string]interface{}{"notexists": "b"})
	assert.Equal(t, errors.New("field not found"), err)
}

func TestNewIndex(t *testing.T) {
	// create keyword
	cidx, err := NewIndex(map[string]map[string]string{"field": {"type": "keyword"}})
	assert.Nil(t, err)
	_, err = cidx.GetFieldIdx("field")
	assert.Nil(t, err)

	// create index missing type
	_, err = NewIndex(map[string]map[string]string{"field": {"ty": "text"}})
	assert.Equal(t, errors.New("missing type"), err)

	// create index invalid type
	_, err = NewIndex(map[string]map[string]string{"field": {"type": "magic"}})
	assert.Equal(t, errors.New("unknown field type"), err)

}

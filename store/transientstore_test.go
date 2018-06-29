package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var k, v = []byte("hello"), []byte("world")

func TestTransientStore(t *testing.T) {
	tstore := newTransientStore()

	assert.Nil(t, tstore.Get(k))

	tstore.Set(k, v)

	assert.Equal(t, v, tstore.Get(k))

	tstore.Commit()

	assert.Nil(t, tstore.Get(k))
}

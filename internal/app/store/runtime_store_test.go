package store_test

import (
	"testing"

	"github.com/BroNaz/queue_broker/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func TestRuntimeStore_Load(t *testing.T) {
	s := store.NewRuntimeStorege()
	s.Store("one", "1")
	s.Store("one", "2")
	s.Store("one", "3")

	vals, ok := s.Load("one")

	assert.True(t, ok)
	assert.Contains(t, vals, "1")
	assert.Contains(t, vals, "2")
	assert.Contains(t, vals, "3")

	vals, ok = s.Load("two")
	assert.False(t, ok)
}

func TestRuntimeStore_Pop(t *testing.T) {
	s := store.NewRuntimeStorege()
	s.Store("one", "1")
	s.Store("one", "2")
	s.Store("one", "3")

	val, ok := s.Pop("one")
	assert.Equal(t, val, "1")
	assert.True(t, ok)

	val, ok = s.Pop("one")
	assert.Equal(t, val, "2")
	assert.True(t, ok)

	val, ok = s.Pop("one")
	assert.Equal(t, val, "3")
	assert.True(t, ok)

	val, ok = s.Pop("one")
	assert.Equal(t, val, "")
	assert.False(t, ok)
}

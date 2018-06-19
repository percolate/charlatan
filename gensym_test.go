package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGensym(t *testing.T) {
	g := gensym()

	assert.Equal(t, g, "G1")
}

func TestSymbolGenerator_Next(t *testing.T) {
	s := symbolGenerator{
		Prefix: "A",
		Suffix: "Z",
	}

	n := s.next()
	n2 := s.next()

	assert.Equal(t, n, "A1Z")
	assert.Equal(t, n2, "A2Z")
}

func TestSymbolGenerator_Reset(t *testing.T) {
	s := symbolGenerator{
		Prefix: "A",
		Suffix: "Z",
	}

	preCount := s.count
	n := s.next()
	onceCount := s.count
	n2 := s.next()
	twiceCount := s.count
	s.reset()

	assert.Equal(t, n, "A1Z")
	assert.Equal(t, n2, "A2Z")
	assert.Equal(t, s.count, uint64(0))
	assert.Equal(t, preCount, uint64(0))
	assert.Equal(t, onceCount, uint64(1))
	assert.Equal(t, twiceCount, uint64(2))
}

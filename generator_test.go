package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadPackageDir(t *testing.T) {
	g, err := LoadPackageDir(".")

	assert.Equal(t, err, nil)
	assert.IsType(t, Generator{}, *g)
}

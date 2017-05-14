package astiloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoader_Add(t *testing.T) {
	var l = New()
	l.Add(3)
	assert.Equal(t, 3, l.t)
}

func TestLoader_Done(t *testing.T) {
	var l = New()
	l.Done(3)
	assert.Equal(t, 3, l.d)
}

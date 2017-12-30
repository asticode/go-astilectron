package astilectron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccelerator(t *testing.T) {
	var tb = []byte("1+2+3")
	var a Accelerator
	err := a.UnmarshalText(tb)
	assert.NoError(t, err)
	assert.Equal(t, Accelerator{"1", "2", "3"}, a)
	b, err := a.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, tb, b)
}

package astilectron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier(t *testing.T) {
	var i = newIdentifier()
	assert.Equal(t, "1", i.new())
	assert.Equal(t, "2", i.new())
	assert.Equal(t, "3", i.new())
}

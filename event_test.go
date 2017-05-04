package astilectron

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventMessage(t *testing.T) {
	// Init
	var em = newEventMessage(false)

	// Test marshal
	var b, err = json.Marshal(em)
	assert.NoError(t, err)
	assert.Equal(t, "false", string(b))

	// Test unmarshal
	err = json.Unmarshal([]byte("true"), em)
	assert.NoError(t, err)
	assert.Equal(t, []byte("true"), em.i)
	var v bool
	err = em.Unmarshal(&v)
	assert.NoError(t, err)
	assert.Equal(t, true, v)
}

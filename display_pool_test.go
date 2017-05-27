package astilectron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplayPool(t *testing.T) {
	// Init
	var dp = newDisplayPool()

	// Test update
	dp.update(&EventDisplays{
		All: []*DisplayOptions{
			{ID: PtrInt(1), Rotation: PtrInt(1)},
			{ID: PtrInt(2)},
		},
		Primary: &DisplayOptions{ID: PtrInt(2)},
	})
	assert.Len(t, dp.all(), 2)
	assert.Equal(t, 2, *dp.primary().o.ID)

	// Test removing one display
	dp.update(&EventDisplays{
		All: []*DisplayOptions{
			{ID: PtrInt(1), Rotation: PtrInt(2)},
		},
		Primary: &DisplayOptions{ID: PtrInt(1)},
	})
	assert.Len(t, dp.all(), 1)
	assert.Equal(t, 2, dp.all()[0].Rotation())
	assert.Equal(t, 1, *dp.primary().o.ID)

	// Test adding a new one
	dp.update(&EventDisplays{
		All: []*DisplayOptions{
			{ID: PtrInt(1)},
			{ID: PtrInt(3)},
		},
		Primary: &DisplayOptions{ID: PtrInt(1)},
	})
	assert.Len(t, dp.all(), 2)
	assert.Equal(t, 1, *dp.primary().o.ID)
}

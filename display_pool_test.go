package astilectron

import (
	"testing"

	"github.com/asticode/go-astikit"
	"github.com/stretchr/testify/assert"
)

func TestDisplayPool(t *testing.T) {
	// Init
	var dp = newDisplayPool()

	// Test update
	dp.update(&EventDisplays{
		All: []*DisplayOptions{
			{ID: astikit.Int64Ptr(1), Rotation: astikit.IntPtr(1)},
			{ID: astikit.Int64Ptr(2)},
		},
		Primary: &DisplayOptions{ID: astikit.Int64Ptr(2)},
	})
	assert.Len(t, dp.all(), 2)
	assert.Equal(t, int64(2), *dp.primary().o.ID)

	// Test removing one display
	dp.update(&EventDisplays{
		All: []*DisplayOptions{
			{ID: astikit.Int64Ptr(1), Rotation: astikit.IntPtr(2)},
		},
		Primary: &DisplayOptions{ID: astikit.Int64Ptr(1)},
	})
	assert.Len(t, dp.all(), 1)
	assert.Equal(t, 2, dp.all()[0].Rotation())
	assert.Equal(t, int64(1), *dp.primary().o.ID)

	// Test adding a new one
	dp.update(&EventDisplays{
		All: []*DisplayOptions{
			{ID: astikit.Int64Ptr(1)},
			{ID: astikit.Int64Ptr(3)},
		},
		Primary: &DisplayOptions{ID: astikit.Int64Ptr(1)},
	})
	assert.Len(t, dp.all(), 2)
	assert.Equal(t, int64(1), *dp.primary().o.ID)
}

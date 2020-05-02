package astilectron

import (
	"testing"

	"github.com/asticode/go-astikit"
	"github.com/stretchr/testify/assert"
)

// TestDisplay tests display
func TestDisplay(t *testing.T) {
	var o = &DisplayOptions{
		Bounds:       &RectangleOptions{PositionOptions: PositionOptions{X: astikit.IntPtr(1), Y: astikit.IntPtr(2)}, SizeOptions: SizeOptions{Height: astikit.IntPtr(3), Width: astikit.IntPtr(4)}},
		ID:           astikit.Int64Ptr(1234),
		Rotation:     astikit.IntPtr(5),
		ScaleFactor:  astikit.Float64Ptr(6),
		Size:         &SizeOptions{Height: astikit.IntPtr(7), Width: astikit.IntPtr(8)},
		TouchSupport: astikit.StrPtr("available"),
		WorkArea:     &RectangleOptions{PositionOptions: PositionOptions{X: astikit.IntPtr(9), Y: astikit.IntPtr(10)}, SizeOptions: SizeOptions{Height: astikit.IntPtr(11), Width: astikit.IntPtr(12)}},
		WorkAreaSize: &SizeOptions{Height: astikit.IntPtr(13), Width: astikit.IntPtr(14)},
	}
	var d = newDisplay(o, true)
	assert.Equal(t, Rectangle{Position: Position{X: 1, Y: 2}, Size: Size{Height: 3, Width: 4}}, d.Bounds())
	assert.Equal(t, int64(1234), d.ID())
	assert.True(t, d.IsPrimary())
	assert.Equal(t, 5, d.Rotation())
	assert.Equal(t, float64(6), d.ScaleFactor())
	assert.Equal(t, Size{Height: 7, Width: 8}, d.Size())
	assert.True(t, d.IsTouchAvailable())
	assert.Equal(t, Rectangle{Position: Position{X: 9, Y: 10}, Size: Size{Height: 11, Width: 12}}, d.WorkArea())
	assert.Equal(t, Size{Height: 13, Width: 14}, d.WorkAreaSize())
	o.TouchSupport = astikit.StrPtr("unavailable")
	d = newDisplay(o, false)
	assert.False(t, d.IsPrimary())
	assert.False(t, d.IsTouchAvailable())
}

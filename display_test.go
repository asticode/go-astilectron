package astilectron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDisplay tests display
func TestDisplay(t *testing.T) {
	var o = &DisplayOptions{
		Bounds:       &RectangleOptions{PositionOptions: PositionOptions{X: PtrInt(1), Y: PtrInt(2)}, SizeOptions: SizeOptions{Height: PtrInt(3), Width: PtrInt(4)}},
		Rotation:     PtrInt(5),
		ScaleFactor:  PtrFloat(6),
		Size:         &SizeOptions{Height: PtrInt(7), Width: PtrInt(8)},
		TouchSupport: PtrStr("available"),
		WorkArea:     &RectangleOptions{PositionOptions: PositionOptions{X: PtrInt(9), Y: PtrInt(10)}, SizeOptions: SizeOptions{Height: PtrInt(11), Width: PtrInt(12)}},
		WorkAreaSize: &SizeOptions{Height: PtrInt(13), Width: PtrInt(14)},
	}
	var d = newDisplay(o, true)
	assert.Equal(t, Rectangle{Position: Position{X: 1, Y: 2}, Size: Size{Height: 3, Width: 4}}, d.Bounds())
	assert.True(t, d.IsPrimary())
	assert.Equal(t, 5, d.Rotation())
	assert.Equal(t, float64(6), d.ScaleFactor())
	assert.Equal(t, Size{Height: 7, Width: 8}, d.Size())
	assert.True(t, d.IsTouchAvailable())
	assert.Equal(t, Rectangle{Position: Position{X: 9, Y: 10}, Size: Size{Height: 11, Width: 12}}, d.WorkArea())
	assert.Equal(t, Size{Height: 13, Width: 14}, d.WorkAreaSize())
	o.TouchSupport = PtrStr("unavailable")
	d = newDisplay(o, false)
	assert.False(t, d.IsPrimary())
	assert.False(t, d.IsTouchAvailable())
}

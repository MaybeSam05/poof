// Package helicopter renders a chopper flying across with a spinning rotor.
package helicopter

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w = 92
	h = 28
)

var (
	red   = color.RGBA{216, 54, 54, 255}
	rdark = color.RGBA{150, 32, 34, 255}
	white = color.RGBA{236, 240, 246, 255}
	glass = color.RGBA{112, 182, 240, 255}
	rotor = color.RGBA{150, 154, 168, 255}
	dark  = color.RGBA{50, 52, 62, 255}
	cloud = color.RGBA{60, 64, 82, 255}
)

var clouds = [][3]int{{10, 22, 5}, {58, 24, 6}, {78, 20, 4}}

func lerp(a, b, i, n int) int {
	if n <= 0 {
		return b
	}
	return a + (b-a)*i/n
}

// drawHeli renders a right-facing helicopter with its top reference at (ox, oy).
// phase animates the main and tail rotors.
func drawHeli(c *renderer.Canvas, ox, oy, phase int) {
	// Main rotor: a grey blur line with a bright blade tip sweeping across.
	c.Rect(ox+14, oy, 56, 1, rotor)
	tip := ox + 14 + (phase%7)*8
	c.Rect(tip, oy, 3, 1, white)
	c.Rect(ox+40, oy+1, 3, 4, dark) // mast

	// Fuselage: red body with a white belly stripe, rounded nose to the right.
	c.Rect(ox+30, oy+5, 28, 8, red)
	c.Disc(ox+58, oy+9, 4, red) // nose
	c.Rect(ox+30, oy+11, 30, 2, white)
	c.Rect(ox+30, oy+12, 30, 1, rdark)
	// Cockpit glass.
	c.Rect(ox+52, oy+6, 7, 4, glass)
	c.Set(ox+58, oy+7, glass)

	// Tail boom + fin + tail rotor.
	c.Rect(ox+8, oy+7, 24, 3, red)
	c.Rect(ox+2, oy+2, 4, 9, red) // vertical fin
	c.Rect(ox+1, oy+2, 2, 2, rdark)
	if phase%2 == 0 {
		c.Rect(ox+2, oy+1, 2, 7, rotor) // tail rotor (vertical blur)
	} else {
		c.Rect(ox+1, oy+4, 4, 2, rotor) // tail rotor (horizontal blur)
	}

	// Skids.
	c.Rect(ox+30, oy+15, 26, 1, dark)
	c.Rect(ox+34, oy+13, 1, 2, dark)
	c.Rect(ox+52, oy+13, 1, 2, dark)
}

// Build assembles the flyby: the chopper crosses left→right with a gentle bob
// and drifting clouds, rotor spinning throughout.
func Build() animation.Scene {
	c := renderer.NewCanvas(w, h)
	var frames []animation.Frame
	add := func(ms int) {
		frames = append(frames, animation.Frame{Grid: c.Snapshot(), DurationMs: ms})
	}
	bob := []int{0, -1, -1, 0, 1, 1}
	steps := 58
	for i := 0; i <= steps; i++ {
		c.Clear()
		for _, cl := range clouds {
			cx := (cl[0] - i) % (w + 12)
			if cx < -6 {
				cx += w + 12
			}
			c.Disc(cx, cl[1], cl[2], cloud)
			c.Disc(cx+cl[2], cl[1]+1, cl[2]-1, cloud)
		}
		x := lerp(-58, w, i, steps)
		drawHeli(c, x, 6+bob[i%len(bob)], i)
		add(30)
	}
	return animation.Scene{Name: "helicopter", Frames: frames}
}

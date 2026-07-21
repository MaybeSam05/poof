// Package alien renders a sleek UFO cruising the lower terminal with a beam.
package alien

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w = 108
	h = 26
)

var (
	hullHi  = color.RGBA{216, 221, 233, 255}
	hull    = color.RGBA{150, 156, 172, 255}
	hullDk  = color.RGBA{88, 94, 114, 255}
	domeHi  = color.RGBA{210, 240, 255, 255}
	dome    = color.RGBA{104, 200, 252, 255}
	domeDk  = color.RGBA{38, 118, 190, 255}
	lightY  = color.RGBA{255, 214, 90, 255}
	lightG  = color.RGBA{150, 255, 150, 255}
	lightR  = color.RGBA{255, 120, 120, 255}
	beamCor = color.RGBA{190, 255, 200, 255}
	beamMid = color.RGBA{110, 224, 140, 255}
	beamEdg = color.RGBA{52, 150, 92, 255}
	star    = color.RGBA{228, 233, 255, 255}
	starDim = color.RGBA{120, 128, 150, 255}
)

var stars = [][3]int{
	{6, 2, 1}, {18, 6, 0}, {33, 1, 1}, {52, 4, 0}, {66, 2, 1},
	{80, 6, 0}, {94, 1, 1}, {101, 8, 0}, {12, 11, 0}, {88, 10, 1}, {44, 9, 0},
}

// drawUFO renders a shaded metallic saucer with a glass dome and rim lights.
func drawUFO(c *renderer.Canvas, cx, cy, blink int) {
	// Glass dome (drawn first; the hull overlaps its base).
	c.Disc(cx, cy-1, 4, domeDk)
	c.Disc(cx, cy-1, 3, dome)
	c.Set(cx-1, cy-3, domeHi)
	c.Set(cx-2, cy-2, domeHi)

	// Hull: a metallic lens shaded top-highlight → mid → dark underside.
	c.Rect(cx-9, cy+1, 18, 1, hullHi)
	c.Rect(cx-13, cy+2, 26, 1, hull)
	c.Rect(cx-15, cy+3, 30, 1, hull)
	c.Rect(cx-13, cy+4, 26, 1, hullDk)
	c.Rect(cx-8, cy+5, 16, 1, hullDk)

	// Rim lights along the widest row, blinking in sequence. 2×2 so they survive
	// the quadrant downsample and read as distinct glowing bulbs.
	cols := []color.RGBA{lightY, lightG, lightR}
	for i, dx := range []int{-13, -7, -1, 5, 11} {
		c.Rect(cx+dx, cy+3, 2, 2, cols[(i+blink)%3])
	}
}

// drawBeam paints a widening tractor beam with a bright core, dim edges and a
// light stipple so it reads as translucent light.
func drawBeam(c *renderer.Canvas, cx, cy, length, phase int) {
	for dy := 0; dy < length; dy++ {
		halfw := 3 + dy*6/maxInt(length, 1)
		for x := -halfw; x <= halfw; x++ {
			// stipple ~2/3 coverage for a see-through glow
			if (x+dy+phase)%3 == 0 {
				continue
			}
			col := beamMid
			if abs(x)*3 < halfw {
				col = beamCor
			} else if abs(x) >= halfw-1 {
				col = beamEdg
			}
			c.Set(cx+x, cy+6+dy, col)
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func lerp(a, b, i, n int) int {
	if n <= 0 {
		return b
	}
	return a + (b-a)*i/n
}

// Build assembles the UFO flyby: drift in, extend and pulse a beam, retract, zip
// away — over a quiet starfield in the lower terminal.
func Build() animation.Scene {
	c := renderer.NewCanvas(w, h)
	var frames []animation.Frame
	add := func(ms int) {
		frames = append(frames, animation.Frame{Grid: c.Snapshot(), DurationMs: ms})
	}
	sky := func() {
		c.Clear()
		for _, s := range stars {
			col := star
			if s[2] == 0 {
				col = starDim
			}
			c.Set(s[0], s[1], col)
		}
	}
	const hoverX, hoverY = 54, 5
	const beamMax = 13

	// Phase 1: drift in from the left.
	for i := 0; i <= 10; i++ {
		sky()
		drawUFO(c, lerp(-16, hoverX, i, 10), hoverY, i)
		add(32)
	}
	// Phase 2: beam extends.
	for i := 1; i <= 6; i++ {
		sky()
		drawBeam(c, hoverX, hoverY, lerp(0, beamMax, i, 6), i)
		drawUFO(c, hoverX, hoverY, i)
		add(40)
	}
	// Phase 3: hover, beam shimmers, lights blink.
	for i := 0; i <= 12; i++ {
		sky()
		drawBeam(c, hoverX, hoverY, beamMax, i)
		drawUFO(c, hoverX, hoverY, i)
		add(58)
	}
	// Phase 4: beam retracts.
	for i := 0; i <= 5; i++ {
		sky()
		drawBeam(c, hoverX, hoverY, lerp(beamMax, 0, i, 5), i)
		drawUFO(c, hoverX, hoverY, i)
		add(40)
	}
	// Phase 5: zip away to the upper right.
	for i := 0; i <= 9; i++ {
		sky()
		drawUFO(c, lerp(hoverX, 124, i, 9), lerp(hoverY, 0, i, 9), i)
		add(26)
	}

	return animation.Scene{Name: "alien", Frames: frames}
}

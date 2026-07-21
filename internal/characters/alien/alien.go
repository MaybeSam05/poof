// Package alien renders a compact metallic UFO cruising with a tractor beam.
package alien

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w = 88
	h = 22
)

var (
	hullHi = color.RGBA{216, 221, 233, 255}
	hull   = color.RGBA{150, 156, 172, 255}
	hullDk = color.RGBA{88, 94, 114, 255}
	domeHi = color.RGBA{210, 240, 255, 255}
	dome   = color.RGBA{104, 200, 252, 255}
	domeDk = color.RGBA{38, 118, 190, 255}
	lightY = color.RGBA{255, 214, 90, 255}
	lightG = color.RGBA{150, 255, 150, 255}
	lightR = color.RGBA{255, 120, 120, 255}
	beamC  = color.RGBA{190, 255, 200, 255}
	beamM  = color.RGBA{110, 224, 140, 255}
	beamE  = color.RGBA{52, 150, 92, 255}
	star   = color.RGBA{228, 233, 255, 255}
	starD  = color.RGBA{120, 128, 150, 255}
)

var stars = [][3]int{
	{5, 2, 1}, {16, 6, 0}, {28, 1, 1}, {60, 3, 0}, {70, 6, 1},
	{82, 2, 0}, {78, 9, 1}, {10, 10, 0}, {36, 9, 0},
}

// drawUFO renders a shaded metallic saucer with a glass dome and rim lights.
func drawUFO(c *renderer.Canvas, cx, cy, blink int) {
	c.Disc(cx, cy-1, 3, domeDk)
	c.Disc(cx, cy-1, 2, dome)
	c.Set(cx-1, cy-2, domeHi)

	c.Rect(cx-6, cy+1, 12, 1, hullHi)
	c.Rect(cx-9, cy+2, 18, 1, hull)
	c.Rect(cx-11, cy+3, 22, 1, hull)
	c.Rect(cx-9, cy+4, 18, 1, hullDk)
	c.Rect(cx-5, cy+5, 10, 1, hullDk)

	cols := []color.RGBA{lightY, lightG, lightR}
	for i, dx := range []int{-9, -3, 3, 8} {
		c.Rect(cx+dx, cy+3, 2, 2, cols[(i+blink)%3])
	}
}

// drawBeam paints a widening, stippled tractor beam with a bright core.
func drawBeam(c *renderer.Canvas, cx, cy, length, phase int) {
	for dy := 0; dy < length; dy++ {
		halfw := 2 + dy*4/maxInt(length, 1)
		for x := -halfw; x <= halfw; x++ {
			if (x+dy+phase)%3 == 0 {
				continue
			}
			col := beamM
			if abs(x)*3 < halfw {
				col = beamC
			} else if abs(x) >= halfw-1 {
				col = beamE
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

// Build assembles the UFO flyby.
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
				col = starD
			}
			c.Set(s[0], s[1], col)
		}
	}
	const hoverX, hoverY = 44, 5
	const beamMax = 10

	for i := 0; i <= 10; i++ {
		sky()
		drawUFO(c, lerp(-14, hoverX, i, 10), hoverY, i)
		add(32)
	}
	for i := 1; i <= 6; i++ {
		sky()
		drawBeam(c, hoverX, hoverY, lerp(0, beamMax, i, 6), i)
		drawUFO(c, hoverX, hoverY, i)
		add(40)
	}
	for i := 0; i <= 12; i++ {
		sky()
		drawBeam(c, hoverX, hoverY, beamMax, i)
		drawUFO(c, hoverX, hoverY, i)
		add(58)
	}
	for i := 0; i <= 5; i++ {
		sky()
		drawBeam(c, hoverX, hoverY, lerp(beamMax, 0, i, 5), i)
		drawUFO(c, hoverX, hoverY, i)
		add(40)
	}
	for i := 0; i <= 9; i++ {
		sky()
		drawUFO(c, lerp(hoverX, 102, i, 9), lerp(hoverY, 0, i, 9), i)
		add(26)
	}

	return animation.Scene{Name: "alien", Frames: frames}
}

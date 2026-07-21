// Package surf renders a surfer carving a curling wave along the lower terminal.
package surf

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w     = 112
	h     = 30
	baseY = 24 // resting ocean surface
)

var (
	deep  = color.RGBA{10, 34, 84, 255}
	mid   = color.RGBA{20, 84, 156, 255}
	face  = color.RGBA{46, 146, 200, 255}
	crest = color.RGBA{118, 212, 232, 255}
	foam  = color.RGBA{236, 246, 250, 255}
)

var surferPalette = []color.RGBA{
	{},                   // 0 transparent
	{58, 40, 24, 255},    // 1 hair
	{238, 186, 144, 255}, // 2 skin
	{26, 34, 54, 255},    // 3 wetsuit
	{60, 78, 120, 255},   // 4 wetsuit highlight
	{234, 66, 58, 255},   // 5 board red
	{240, 242, 248, 255}, // 6 board white
}

// A surfer in a low carving crouch, facing left (the direction of travel):
// front arm reaching forward, back arm out for balance, knees bent on the board.
var surfer = renderer.ParseSprite(surferPalette, []string{
	"       11    ",
	"      1221   ",
	"      2222   ",
	"       22    ",
	"    44322    ",
	"  2243332    ",
	"    433332   ",
	"    433332   ",
	"    33333    ",
	"    33 33    ",
	"   33   33   ",
	"  33     3   ",
	" 2       3   ",
	"5666666655   ",
})

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func lerp(a, b, i, n int) int {
	if n <= 0 {
		return b
	}
	return a + (b-a)*i/n
}

// surfaceY returns the top of the water at column x for a swell peaking at peakX:
// a steep front (left) face and a gentle back.
func surfaceY(x, peakX, amp int) int {
	var a int
	if x <= peakX {
		a = amp - (peakX-x)*3/2
	} else {
		a = amp - (x-peakX)/2
	}
	if a < 0 {
		a = 0
	}
	return baseY - a
}

// drawWave fills the ocean with the swell, shades it crest→deep, lays foam along
// the crest and hooks a curling lip that pitches forward over the face.
func drawWave(c *renderer.Canvas, peakX, amp int) {
	for x := 0; x < w; x++ {
		top := surfaceY(x, peakX, amp)
		for y := top; y < h; y++ {
			switch d := y - top; {
			case d == 0:
				c.Set(x, y, crest)
			case d <= 2:
				c.Set(x, y, face)
			case d <= 6:
				c.Set(x, y, mid)
			default:
				c.Set(x, y, deep)
			}
		}
	}
	if amp < 6 {
		return
	}
	peakTop := surfaceY(peakX, peakX, amp)
	// Chunky whitewater cap breaking over the crest, with a little spray.
	c.Rect(peakX-4, peakTop-1, 11, 2, foam)
	c.Rect(peakX-2, peakTop-3, 6, 2, foam)
	c.Set(peakX, peakTop-4, foam)
	c.Set(peakX+3, peakTop-4, foam)
	// Foam spilling a short way down the front face.
	c.Rect(peakX-6, peakTop+1, 3, 1, foam)
	c.Rect(peakX-8, peakTop+2, 2, 1, foam)
}

// Build assembles the surfing scene: a swell travels right→left with a surfer
// carving its face, then the crest breaks into foam that spreads and settles.
func Build() animation.Scene {
	c := renderer.NewCanvas(w, h)
	var frames []animation.Frame
	add := func(ms int) {
		frames = append(frames, animation.Frame{Grid: c.Snapshot(), DurationMs: ms})
	}
	rideSurfer := func(peakX, amp int) {
		sx := peakX - 15 // on the open face, ahead of the curl
		boardY := surfaceY(sx+6, peakX, amp)
		c.Blit(surfer, sx, boardY-13)
	}

	const amp = 15

	// Phase 1: swell builds in from the right.
	for i := 0; i <= 6; i++ {
		c.Clear()
		drawWave(c, 96, lerp(2, amp, i, 6))
		add(38)
	}
	// Phase 2: the swell rolls left, surfer carving the face.
	steps := 42
	for i := 0; i <= steps; i++ {
		c.Clear()
		px := lerp(96, 26, i, steps)
		drawWave(c, px, amp)
		rideSurfer(px, amp)
		add(30)
	}
	// Phase 3: the wave breaks — foam spreads across the surface.
	for i := 0; i <= 6; i++ {
		c.Clear()
		drawWave(c, 20, lerp(amp, 0, i, 6))
		c.Rect(0, baseY-2, lerp(0, w, i, 6), 3, foam)
		add(40)
	}
	// Phase 4: foam settles back into calm water.
	for i := 0; i <= 4; i++ {
		c.Clear()
		for x := 0; x < w; x++ {
			for y := baseY; y < h; y++ {
				if y-baseY <= 1 {
					c.Set(x, y, mid)
				} else {
					c.Set(x, y, deep)
				}
			}
		}
		start := lerp(0, w, i, 4)
		if start < w {
			c.Rect(start, baseY-1, w-start, 2, foam)
		}
		add(55)
	}

	return animation.Scene{Name: "surf", Frames: frames}
}

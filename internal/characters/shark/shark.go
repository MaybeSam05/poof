// Package shark renders a shark fin gliding, then a breach and splash.
package shark

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w   = 98
	h   = 30
	sea = 20 // water surface
)

var (
	grey  = color.RGBA{104, 112, 128, 255}
	dgrey = color.RGBA{62, 70, 86, 255}
	belly = color.RGBA{224, 230, 242, 255}
	mouth = color.RGBA{66, 26, 30, 255}
	teeth = color.RGBA{238, 242, 248, 255}
	eye   = color.RGBA{18, 20, 26, 255}
	deep  = color.RGBA{10, 42, 96, 255}
	mid   = color.RGBA{22, 82, 150, 255}
	foam  = color.RGBA{234, 244, 250, 255}
)

// Side-profile shark facing left (snout left, crescent tail right).
var shark = renderer.ParseSprite([]color.RGBA{
	{}, grey, dgrey, belly, mouth, teeth, eye,
}, []string{
	"              222         ",
	"             22222        ",
	"         11111111111      ",
	"      1111111111111111112  ",
	"    1111111111111111111122 ",
	"  11111111111111111111112222",
	" 1e5111111111111111111112222",
	"1451111111111111111111112222",
	" 33311111111111111111111 22 ",
	"  333331111111222221        ",
	"        222222              ",
})

func lerp(a, b, i, n int) int {
	if n <= 0 {
		return b
	}
	return a + (b-a)*i/n
}

func drawSea(c *renderer.Canvas) {
	for y := sea; y < h; y++ {
		col := mid
		if y-sea > 2 {
			col = deep
		}
		c.HSpan(0, y, w, col)
	}
	c.HSpan(0, sea, w, foam)
}

// drawFin draws a dorsal fin cutting the surface at column fx, with a wake.
func drawFin(c *renderer.Canvas, fx, t int) {
	for dy := 0; dy < 5; dy++ {
		c.Rect(fx-dy, sea-1-dy, dy+1, 1, grey)
	}
	c.Set(fx-4, sea-5, dgrey)
	// V-wake trailing to the right (behind a left-moving fin).
	for i := 1; i < 10; i++ {
		if (i+t)%2 == 0 {
			c.Set(fx+i*2, sea-1, foam)
			c.Set(fx+i*2, sea+1, foam)
		}
	}
}

func splash(c *renderer.Canvas, x, size int) {
	for i := 0; i < size; i++ {
		c.Set(x-i, sea-i/2, foam)
		c.Set(x+i, sea-i/3, foam)
	}
	c.Disc(x, sea-1, size/3, foam)
}

// Build assembles the scene: the fin glides in, the shark breaches in an arc,
// then splashes back with spreading foam.
func Build() animation.Scene {
	c := renderer.NewCanvas(w, h)
	var frames []animation.Frame
	add := func(ms int) {
		frames = append(frames, animation.Frame{Grid: c.Snapshot(), DurationMs: ms})
	}

	// Phase 1: fin glides right→left toward the launch point.
	for i := 0; i <= 14; i++ {
		c.Clear()
		drawSea(c)
		drawFin(c, lerp(w-6, 62, i, 14), i)
		add(45)
	}
	// Phase 2 & 3: breach arc — leap up-left out of the water and back down.
	steps := 22
	for i := 0; i <= steps; i++ {
		c.Clear()
		drawSea(c)
		sx := lerp(60, 16, i, steps)
		// parabola: top of the sprite peaks near the middle of the arc.
		d := 2*i - steps
		peak := 14
		top := sea - 4 - (peak - peak*d*d/(steps*steps))
		c.Blit(shark, sx, top)
		// launch splash early, entry splash late
		if i < 6 {
			splash(c, 66, 6-i)
		}
		if i > steps-6 {
			splash(c, 18, i-(steps-6)+2)
		}
		add(34)
	}
	// Phase 4: foam settles where it re-entered.
	for i := 0; i <= 5; i++ {
		c.Clear()
		drawSea(c)
		r := 5 - i
		if r > 0 {
			for _, o := range []int{-4, 0, 5} {
				c.Disc(18+o, sea-1, r-(abs(o)/3), foam)
			}
		}
		add(55)
	}

	return animation.Scene{Name: "shark", Frames: frames}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

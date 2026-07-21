// Package dino renders a running T-rex hopping a cactus on a scrolling desert.
package dino

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w      = 80
	h      = 24
	ground = 20
)

var (
	green  = color.RGBA{98, 188, 90, 255}
	dgreen = color.RGBA{58, 138, 62, 255}
	belly  = color.RGBA{162, 218, 146, 255}
	eyeW   = color.RGBA{246, 249, 242, 255}
	eyeP   = color.RGBA{22, 26, 22, 255}
	sand   = color.RGBA{198, 168, 108, 255}
	sandD  = color.RGBA{150, 122, 74, 255}
	cactus = color.RGBA{64, 156, 96, 255}
	cactD  = color.RGBA{42, 112, 70, 255}
	sky    = color.RGBA{88, 96, 118, 255}
)

// dinoBody is a right-facing T-rex from the hips up (legs are drawn separately
// so they can cycle). Tail at the left, head at the right.
var dinoBody = renderer.ParseSprite([]color.RGBA{
	{}, green, dgreen, belly, eyeW, eyeP,
}, []string{
	"              11111 ",
	"             1111111",
	"             1145111",
	"             111111 ",
	"  1          11111  ",
	" 11        1111111  ",
	"111      111111111  ",
	" 1111111111111111   ",
	"  1113311111111     ",
	"   1133111111       ",
	"     11111111       ",
})

func lerp(a, b, i, n int) int {
	if n <= 0 {
		return b
	}
	return a + (b-a)*i/n
}

// drawLegs draws the two legs from the hip line at (ox,topY), cycling by phase.
// airborne tucks both legs up.
func drawLegs(c *renderer.Canvas, ox, topY, phase int, airborne bool) {
	legY := topY + 11
	front, back := ox+11, ox+6
	if airborne {
		c.Rect(back, legY, 2, 2, dgreen)
		c.Rect(front, legY, 2, 2, dgreen)
		return
	}
	if phase == 0 {
		c.Rect(back, legY, 2, 3, dgreen)
		c.Set(back-1, legY+3, dgreen)
		c.Rect(front, legY, 2, 2, dgreen)
		c.Set(front+1, legY+2, dgreen)
	} else {
		c.Rect(back, legY, 2, 2, dgreen)
		c.Set(back+1, legY+2, dgreen)
		c.Rect(front, legY, 2, 3, dgreen)
		c.Set(front+1, legY+3, dgreen)
	}
}

func drawCactus(c *renderer.Canvas, x, baseY int) {
	c.Rect(x+2, baseY-8, 2, 9, cactus) // trunk
	c.Rect(x, baseY-5, 2, 3, cactus)   // left arm
	c.Set(x, baseY-6, cactus)
	c.Rect(x+4, baseY-6, 2, 3, cactus) // right arm
	c.Set(x+5, baseY-7, cactus)
	c.Rect(x+2, baseY-8, 1, 9, cactD) // shade
}

func drawGround(c *renderer.Canvas, scroll int) {
	c.HSpan(0, ground, w, sandD)
	c.Rect(0, ground+1, w, h-ground-1, sand)
	// scrolling speckle for a sense of motion
	for x := 0; x < w; x++ {
		if (x+scroll)%9 == 0 {
			c.Set(x, ground+2, sandD)
		}
		if (x+scroll)%13 == 0 {
			c.Set(x, ground+3, sandD)
		}
	}
}

// Build assembles the run: the dino runs in place while the desert scrolls, a
// cactus approaches, the dino hops it, then it runs on.
func Build() animation.Scene {
	c := renderer.NewCanvas(w, h)
	var frames []animation.Frame
	add := func(ms int) {
		frames = append(frames, animation.Frame{Grid: c.Snapshot(), DurationMs: ms})
	}
	const dinoX = 12
	bodyTop := ground - 14 // feet rest on the ground

	// hopArc returns the upward offset (0 = grounded) for jump progress 0..n.
	hopArc := func(i, n, height int) int {
		// parabola peaking at the middle
		d := 2*i - n
		return height - height*d*d/(n*n)
	}

	// A cactus travels leftward across the whole run; the dino hops as it nears.
	total := 46
	cactusStart := w + 4
	for i := 0; i <= total; i++ {
		c.Clear()
		scroll := i * 3
		drawGround(c, scroll)
		cx := cactusStart - i*3
		if cx > -8 && cx < w {
			drawCactus(c, cx, ground)
		}
		// Jump when the cactus is in front of the dino.
		lift := 0
		airborne := false
		jumpStart := (cactusStart - (dinoX + 20)) / 3
		if i >= jumpStart && i <= jumpStart+10 {
			lift = hopArc(i-jumpStart, 10, 6)
			airborne = lift > 1
		}
		top := bodyTop - lift
		drawLegs(c, dinoX, top, (i/2)%2, airborne)
		c.Blit(dinoBody, dinoX, top)
		add(60)
	}

	return animation.Scene{Name: "dino", Frames: frames}
}

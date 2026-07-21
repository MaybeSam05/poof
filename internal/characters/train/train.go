// Package train renders a steam locomotive chugging across the lower terminal.
package train

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w     = 98
	h     = 26
	track = 22
)

var (
	green  = color.RGBA{36, 116, 78, 255}
	gdark  = color.RGBA{22, 82, 56, 255}
	black  = color.RGBA{40, 42, 52, 255}
	red    = color.RGBA{208, 46, 46, 255}
	brass  = color.RGBA{208, 168, 82, 255}
	window = color.RGBA{255, 216, 112, 255}
	steel  = color.RGBA{150, 154, 168, 255}
	rail   = color.RGBA{70, 72, 84, 255}
	smk1   = color.RGBA{158, 160, 170, 255}
	smk2   = color.RGBA{104, 106, 118, 255}
)

func lerp(a, b, i, n int) int {
	if n <= 0 {
		return b
	}
	return a + (b-a)*i/n
}

// drawWheel: dark wheel with a red hub and rotating spokes.
func drawWheel(c *renderer.Canvas, cx, cy, r, spin int) {
	c.Disc(cx, cy, r, black)
	c.Disc(cx, cy, r-1, gdark)
	if spin%2 == 0 {
		for i := -(r - 1); i <= r-1; i++ {
			c.Set(cx+i, cy, steel)
			c.Set(cx, cy+i, steel)
		}
	} else {
		for i := -(r - 1); i <= r-1; i++ {
			if 2*i*i <= (r-1)*(r-1)+r {
				c.Set(cx+i, cy+i, steel)
				c.Set(cx+i, cy-i, steel)
			}
		}
	}
	c.Disc(cx, cy, 1, red)
}

// drawTrain renders a right-facing steam engine anchored at relative x.
func drawTrain(c *renderer.Canvas, x, spin int) {
	// Cab (rear/left).
	c.Rect(x+2, 6, 18, 12, green)
	c.Rect(x+1, 4, 20, 2, black) // roof
	c.Rect(x+4, 8, 12, 5, window)
	c.Rect(x+4, 8, 12, 1, black)

	// Boiler.
	c.Rect(x+18, 9, 48, 9, green)
	c.Rect(x+18, 9, 48, 1, gdark)
	for _, bx := range []int{x + 26, x + 40, x + 54} {
		c.Rect(bx, 9, 1, 9, brass) // boiler bands
	}
	// Steam + sand domes.
	c.Rect(x+34, 6, 6, 3, brass)
	c.Rect(x+46, 7, 4, 2, brass)

	// Smokebox (front) + headlight.
	c.Rect(x+62, 8, 8, 10, black)
	c.Disc(x+66, 13, 3, gdark)
	c.Rect(x+64, 11, 2, 2, window) // headlight

	// Chimney.
	c.Rect(x+54, 3, 6, 6, black)
	c.Rect(x+52, 2, 10, 2, black) // flared cap

	// Cowcatcher (pilot) + running board.
	c.Rect(x+18, 18, 52, 1, black)
	for i := 0; i < 5; i++ {
		c.Rect(x+70+i, 14+i, 2, track-1-(14+i), red)
	}

	// Wheels: three big drivers + a small pilot wheel.
	drawWheel(c, x+28, track-3, 4, spin)
	drawWheel(c, x+42, track-3, 4, spin)
	drawWheel(c, x+56, track-3, 4, spin)
	drawWheel(c, x+70, track-2, 2, spin)
	// Coupling rod linking the drivers.
	c.Rect(x+28, track-3+3, 30, 1, steel)
}

// drawSmoke draws a plume streaming back over the engine from the chimney at
// (chX, chTop) — billowing puffs that grow and drift as they trail behind.
func drawSmoke(c *renderer.Canvas, chX, chTop, t int) {
	for i := 0; i < 9; i++ {
		px := chX - 1 - i*4 + t%3
		py := chTop - i/3
		if py < 0 {
			py = 0
		}
		r := 1 + i/2
		col := smk1
		if (i+t)%2 == 0 {
			col = smk2
		}
		c.Disc(px, py, r, col)
	}
}

// Build assembles the train chugging left→right with a trailing smoke plume.
func Build() animation.Scene {
	c := renderer.NewCanvas(w, h)
	var frames []animation.Frame
	add := func(ms int) {
		frames = append(frames, animation.Frame{Grid: c.Snapshot(), DurationMs: ms})
	}
	steps := 60
	for i := 0; i <= steps; i++ {
		c.Clear()
		// Rails.
		c.HSpan(0, track, w, rail)
		c.HSpan(0, track+2, w, rail)
		for x := 0; x < w; x++ {
			if (x+i*2)%6 == 0 {
				c.Rect(x, track+1, 1, 1, smk2) // sleepers
			}
		}
		x := lerp(-78, w, i, steps)
		drawSmoke(c, x+57, 3, i)
		drawTrain(c, x, i)
		add(30)
	}
	return animation.Scene{Name: "train", Frames: frames}
}

// Package f1 renders a compact red-and-white Formula 1 car across the terminal.
package f1

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w      = 84
	h      = 22
	ground = 19
	axleY  = 14
	wheelR = 5
	floorY = 16
)

var (
	red    = color.RGBA{216, 34, 42, 255}
	rdark  = color.RGBA{148, 20, 28, 255}
	white  = color.RGBA{238, 240, 246, 255}
	tire   = color.RGBA{30, 30, 36, 255}
	floor  = color.RGBA{22, 22, 28, 255}
	cockpt = color.RGBA{16, 16, 22, 255}
	grey   = color.RGBA{120, 125, 140, 255}
	track  = color.RGBA{56, 58, 68, 255}
	speed  = color.RGBA{120, 124, 140, 255}
)

func lerp(a, b, i, n int) int {
	if n <= 0 {
		return b
	}
	return a + (b-a)*i/n
}

// topEdge is the top of the bodywork at relative x (nose at the left).
func topEdge(rx int) int {
	switch {
	case rx < 10:
		return 13
	case rx < 22:
		return 13 - (rx-10)*3/12
	case rx < 28:
		return 10 - (rx-22)*2/6
	case rx < 36:
		return 8 - (rx-28)*4/8
	case rx < 42:
		return 4
	case rx < 64:
		return 4 + (rx-42)*6/22
	default:
		return 10
	}
}

// drawWheel renders a black tire with a thin white rim ring, spokes and a hub.
func drawWheel(c *renderer.Canvas, cx, cy, spin int) {
	c.Disc(cx, cy, wheelR, tire)
	c.Disc(cx, cy, wheelR-2, white)
	c.Disc(cx, cy, wheelR-3, tire)
	if spin%2 == 0 {
		c.Set(cx-2, cy, white)
		c.Set(cx+2, cy, white)
		c.Set(cx, cy-2, white)
		c.Set(cx, cy+2, white)
	} else {
		c.Set(cx-1, cy-1, white)
		c.Set(cx+1, cy+1, white)
		c.Set(cx-1, cy+1, white)
		c.Set(cx+1, cy-1, white)
	}
	c.Set(cx, cy, grey)
}

// drawCar renders a left-facing red/white F1 car at relative x.
func drawCar(c *renderer.Canvas, x, spin int) {
	// Red bodywork shell.
	for rx := 8; rx <= 72; rx++ {
		t := topEdge(rx)
		if t < floorY {
			c.Rect(x+rx, t, 1, floorY-t, red)
		}
	}
	c.HSpan(x+10, floorY-1, 60, rdark)
	c.Rect(x+14, floorY, 54, 2, floor)

	// White top stripe along the engine cover.
	for rx := 28; rx <= 64; rx++ {
		c.Set(x+rx, topEdge(rx), white)
		c.Set(x+rx, topEdge(rx)+1, white)
	}

	// Sidepod: white band + dark inlet.
	c.Rect(x+36, 12, 20, 2, white)
	c.Rect(x+36, 10, 6, 3, cockpt)

	// Airbox / roll hoop.
	c.Rect(x+38, 3, 5, 4, red)
	c.Rect(x+39, 4, 3, 2, cockpt)

	// Cockpit + halo + driver.
	c.Rect(x+28, 7, 9, 2, cockpt)
	c.Rect(x+30, 4, 4, 3, cockpt)
	c.Set(x+31, 5, red)
	c.HSpan(x+27, 6, 11, cockpt)
	c.Rect(x+27, 6, 1, 3, cockpt)

	// White nose cone tapering to the tip.
	for rx := 2; rx < 10; rx++ {
		top := 13 - (rx-2)/3
		c.Rect(x+rx, top, 1, floorY-1-top, white)
	}

	// Front wing (far left, low).
	c.Rect(x+0, floorY-1, 10, 1, red)
	c.Rect(x+0, floorY, 10, 1, white)
	c.Rect(x+0, floorY-3, 1, 4, red)

	// Rear wing (far right, tall).
	c.Rect(x+72, 3, 3, 13, floor)
	c.Rect(x+68, 3, 9, 2, floor)
	c.HSpan(x+68, 3, 9, white)
	c.Rect(x+67, 13, 10, 2, floor)

	drawWheel(c, x+18, axleY, spin)
	drawWheel(c, x+58, axleY, spin)
}

// Build assembles the flyby: the car enters from the right and sweeps left.
func Build() animation.Scene {
	c := renderer.NewCanvas(w, h)
	var frames []animation.Frame
	add := func(ms int) {
		frames = append(frames, animation.Frame{Grid: c.Snapshot(), DurationMs: ms})
	}
	scene := func(carX, spin, streak int) {
		c.Clear()
		c.HSpan(0, ground, w, track)
		if streak > 0 {
			for _, y := range []int{axleY - 5, axleY, floorY} {
				for sx := carX + 76; sx < carX+76+streak; sx += 5 {
					if sx < w {
						c.Rect(sx, y, 3, 1, speed)
					}
				}
			}
		}
		drawCar(c, carX, spin)
	}

	steps := 56
	for i := 0; i <= steps; i++ {
		x := lerp(w, -80, i, steps)
		streak := 14
		if i < 8 {
			streak = i * 2
		}
		scene(x, i, streak)
		add(26)
	}
	for i := 0; i < 5; i++ {
		c.Clear()
		c.HSpan(0, ground, w, track)
		for _, y := range []int{axleY - 5, axleY, floorY} {
			for sx := w - 1; sx > i*18; sx -= 6 {
				c.Rect(sx, y, 2, 1, speed)
			}
		}
		add(45)
	}

	return animation.Scene{Name: "f1", Frames: frames}
}

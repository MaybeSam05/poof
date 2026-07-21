// Package f1 renders a modern papaya/blue Formula 1 car across the lower terminal.
package f1

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w      = 108
	h      = 28
	ground = 25
	axleY  = 19
	wheelR = 6
	floorY = 21
)

var (
	orange = color.RGBA{247, 130, 28, 255}
	odark  = color.RGBA{201, 94, 14, 255}
	blue   = color.RGBA{32, 112, 236, 255}
	tire   = color.RGBA{30, 30, 36, 255}
	floor  = color.RGBA{22, 22, 28, 255}
	cockpt = color.RGBA{16, 16, 22, 255}
	white  = color.RGBA{238, 241, 247, 255}
	grey   = color.RGBA{118, 123, 138, 255}
	track  = color.RGBA{56, 58, 68, 255}
	speed  = color.RGBA{120, 124, 140, 255}
)

func lerp(a, b, i, n int) int {
	if n <= 0 {
		return b
	}
	return a + (b-a)*i/n
}

// topEdge is the top of the bodywork at relative x (nose at the left): a low
// nose rising to the cockpit, a tall airbox peak, then the engine cover sloping
// back to the rear wing.
func topEdge(rx int) int {
	switch {
	case rx < 12:
		return 17
	case rx < 26:
		return 17 - (rx-12)*3/14
	case rx < 34:
		return 14 - (rx-26)*4/8
	case rx < 44:
		return 10 - (rx-34)*4/10
	case rx < 52:
		return 6
	case rx < 78:
		return 6 + (rx-52)*7/26
	default:
		return 13
	}
}

// drawWheel renders a black tire with a thin white rim ring, faint spokes and a
// hub — mostly black, like the reference.
func drawWheel(c *renderer.Canvas, cx, cy, spin int) {
	c.Disc(cx, cy, wheelR, tire)    // black tire
	c.Disc(cx, cy, wheelR-2, white) // white rim ring...
	c.Disc(cx, cy, wheelR-3, tire)  // ...over a black rim face (1px white ring)
	if spin%2 == 0 {                // short rotating spokes
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

// drawCar renders a left-facing papaya/blue F1 car at relative x.
func drawCar(c *renderer.Canvas, x, spin int) {
	// Orange bodywork shell, from the top edge down to the tub line.
	for rx := 10; rx <= 88; rx++ {
		t := topEdge(rx)
		if t < floorY {
			c.Rect(x+rx, t, 1, floorY-t, orange)
		}
	}
	// Lower-body / underside shadow and black floor.
	c.HSpan(x+12, floorY-1, 74, odark)
	c.Rect(x+16, floorY, 66, 2, floor)

	// Blue top stripe along the engine cover + a blue nose tip.
	for rx := 34; rx <= 78; rx++ {
		c.Set(x+rx, topEdge(rx), blue)
		c.Set(x+rx, topEdge(rx)+1, blue)
	}
	c.Rect(x+10, 15, 4, 3, blue) // nose tip

	// Sidepod: orange top, blue lower band, dark inlet at its leading edge.
	c.Rect(x+42, 15, 26, 4, blue)
	c.Rect(x+42, 13, 8, 3, cockpt) // sidepod inlet
	c.Rect(x+52, 16, 6, 2, white)  // livery swoosh

	// Airbox / roll hoop behind the cockpit (blue with a dark intake).
	c.Rect(x+46, 4, 6, 4, blue)
	c.Rect(x+47, 5, 3, 2, cockpt)

	// Cockpit + halo + driver.
	c.Rect(x+34, 9, 10, 2, cockpt) // opening
	c.Rect(x+36, 6, 4, 3, cockpt)  // helmet/head
	c.Set(x+37, 7, blue)           // helmet flash
	c.HSpan(x+33, 8, 13, cockpt)   // halo bar
	c.Rect(x+33, 8, 1, 3, cockpt)  // rear halo post

	// Nose cone (orange) tapering to the tip at the left.
	for rx := 2; rx < 12; rx++ {
		top := 17 - (rx-2)/3
		c.Rect(x+rx, top, 1, floorY-1-top, orange)
	}

	// Front wing (far left, low): blue plane + orange flap + endplate.
	c.Rect(x+0, floorY-1, 12, 1, blue)
	c.Rect(x+0, floorY, 12, 1, orange)
	c.Rect(x+0, floorY-3, 1, 4, blue)

	// Rear wing (far right, tall) with endplate and beam.
	c.Rect(x+87, 4, 3, 16, floor)  // main endplate
	c.Rect(x+82, 4, 10, 3, floor)  // top plane
	c.HSpan(x+82, 4, 10, blue)     // wing flash
	c.Rect(x+81, 16, 11, 2, floor) // beam wing

	// Wheels on top.
	drawWheel(c, x+22, axleY, spin)
	drawWheel(c, x+72, axleY, spin)
}

// Build assembles the flyby: the car enters from the right and sweeps left with
// faint speed lines, then a brief settle.
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
			for _, y := range []int{axleY - 6, axleY, floorY} {
				for sx := carX + 92; sx < carX+92+streak; sx += 6 {
					if sx < w {
						c.Rect(sx, y, 4, 1, speed)
					}
				}
			}
		}
		drawCar(c, carX, spin)
	}

	steps := 62
	for i := 0; i <= steps; i++ {
		x := lerp(w, -96, i, steps)
		streak := 18
		if i < 8 {
			streak = i * 2
		}
		scene(x, i, streak)
		add(26)
	}
	for i := 0; i < 5; i++ {
		c.Clear()
		c.HSpan(0, ground, w, track)
		for _, y := range []int{axleY - 6, axleY, floorY} {
			for sx := w - 1; sx > i*22; sx -= 7 {
				c.Rect(sx, y, 3, 1, speed)
			}
		}
		add(45)
	}

	return animation.Scene{Name: "f1", Frames: frames}
}

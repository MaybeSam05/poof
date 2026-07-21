// Package rocket renders a rocket blasting off from the lower terminal.
package rocket

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w   = 56
	h   = 38
	pad = 34 // ground / launch-pad line
)

var (
	white  = color.RGBA{236, 240, 246, 255}
	grey   = color.RGBA{168, 174, 190, 255}
	red    = color.RGBA{220, 44, 50, 255}
	winBl  = color.RGBA{92, 182, 255, 255}
	winHi  = color.RGBA{206, 236, 255, 255}
	dark   = color.RGBA{44, 46, 56, 255}
	yellow = color.RGBA{255, 232, 96, 255}
	orange = color.RGBA{255, 150, 40, 255}
	ember  = color.RGBA{224, 82, 22, 255}
	smoke1 = color.RGBA{150, 152, 162, 255}
	smoke2 = color.RGBA{96, 98, 110, 255}
	star   = color.RGBA{228, 233, 255, 255}
	starD  = color.RGBA{116, 124, 148, 255}
)

var stars = [][3]int{
	{6, 3, 1}, {16, 8, 0}, {40, 2, 1}, {49, 7, 0}, {11, 14, 0},
	{45, 15, 1}, {51, 22, 0}, {4, 20, 1}, {38, 20, 0},
}

const rocketH = 22

// drawRocket renders a rocket whose nose tip is at (cx, topY).
func drawRocket(c *renderer.Canvas, cx, topY int) {
	bodyTop := topY + 5
	bodyBot := topY + 18
	// Nose cone (red), widening from the tip.
	for dy := 0; dy < 5; dy++ {
		ww := 2 + dy
		c.Rect(cx-ww/2, topY+dy, ww, 1, red)
	}
	// Body tube.
	c.Rect(cx-3, bodyTop, 7, bodyBot-bodyTop, white)
	c.Rect(cx+2, bodyTop, 1, bodyBot-bodyTop, grey) // right-side shade
	c.Rect(cx-3, bodyTop, 1, bodyBot-bodyTop, white)
	// A grey band + window.
	c.Rect(cx-3, bodyTop+8, 7, 1, grey)
	c.Disc(cx, bodyTop+3, 2, winBl)
	c.Set(cx-1, bodyTop+2, winHi)
	// Fins (red) flaring out near the base.
	for dy := 0; dy < 6; dy++ {
		fw := 1 + dy/2
		c.Rect(cx-3-fw, bodyBot-6+dy, fw, 1, red)
		c.Rect(cx+4, bodyBot-6+dy, fw, 1, red)
	}
	// Nozzle.
	c.Rect(cx-2, bodyBot, 5, 2, dark)
	c.Rect(cx-1, bodyBot+2, 3, 1, dark)
}

// drawFlame renders an exhaust plume of length hanging below y0, centered at cx.
func drawFlame(c *renderer.Canvas, cx, y0, length, flick int) {
	if length <= 0 {
		return
	}
	for dy := 0; dy < length; dy++ {
		hw := (length-dy)*3/length + 1
		// flicker the width a touch per row/frame
		if (dy+flick)%3 == 0 && hw > 1 {
			hw--
		}
		for x := -hw; x <= hw; x++ {
			col := orange
			if abs(x)*2 < hw {
				col = yellow
			}
			if dy > length*3/4 || abs(x) == hw {
				col = ember
			}
			c.Set(cx+x, y0+dy, col)
		}
	}
}

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

// Build assembles the launch: ignition, liftoff climbing off the top, then the
// smoke cloud settling on the pad.
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
		c.HSpan(0, pad, w, smoke2) // pad line
	}
	cx := w / 2
	smoke := func(spread int) {
		for i, o := range [][2]int{{-1, -1}, {2, 0}, {-3, 1}, {4, 1}, {0, 2}} {
			r := 2 + (spread*(i%3+1))/6
			cc := smoke1
			if i%2 == 0 {
				cc = smoke2
			}
			c.Disc(cx+o[0]*spread/4, pad+o[1]-1, r, cc)
		}
	}

	restTop := pad - rocketH // nozzle rests on the pad

	// Phase 1: ignition — flame grows, first smoke.
	for i := 0; i <= 5; i++ {
		sky()
		smoke(lerp(1, 5, i, 5))
		drawRocket(c, cx, restTop)
		drawFlame(c, cx, restTop+21, lerp(1, 7, i, 5), i)
		add(45)
	}
	// Phase 2: liftoff — accelerate up and off the top.
	steps := 20
	for i := 0; i <= steps; i++ {
		sky()
		smoke(6)
		// ease-in: quadratic rise
		top := restTop - (restTop+rocketH+6)*i*i/(steps*steps)
		drawRocket(c, cx, top)
		drawFlame(c, cx, top+21, 7+i/3, i)
		add(34)
	}
	// Phase 3: smoke lingers and thins.
	for i := 0; i <= 5; i++ {
		sky()
		for j, o := range [][2]int{{-4, 0}, {0, 1}, {5, 0}, {-1, 2}} {
			r := 4 - i/2 + j%2
			if r > 0 {
				c.Disc(cx+o[0]*(i+2)/3, pad+o[1]-1, r, smoke2)
			}
		}
		add(60)
	}

	return animation.Scene{Name: "rocket", Frames: frames}
}

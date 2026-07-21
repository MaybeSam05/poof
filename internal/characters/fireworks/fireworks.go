// Package fireworks renders shells launching and bursting over a city skyline.
package fireworks

import (
	"image/color"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/renderer"
)

const (
	w      = 90
	h      = 40
	bottom = 37
)

var (
	white   = color.RGBA{246, 249, 255, 255}
	trail   = color.RGBA{255, 224, 150, 255}
	bldg    = color.RGBA{26, 28, 40, 255}
	bldgLit = color.RGBA{58, 60, 80, 255}
	win     = color.RGBA{255, 214, 120, 255}
	star    = color.RGBA{210, 216, 240, 255}
)

// A palette entry: bright spark color + a dimmer fade color.
type hue struct{ bright, dim color.RGBA }

var (
	red   = hue{color.RGBA{255, 78, 78, 255}, color.RGBA{150, 40, 44, 255}}
	gold  = hue{color.RGBA{255, 206, 74, 255}, color.RGBA{158, 120, 36, 255}}
	green = hue{color.RGBA{96, 240, 128, 255}, color.RGBA{44, 140, 74, 255}}
	blue  = hue{color.RGBA{96, 176, 255, 255}, color.RGBA{44, 96, 160, 255}}
	pink  = hue{color.RGBA{255, 118, 206, 255}, color.RGBA{158, 60, 120, 255}}
)

// 12-spoke burst directions (roughly radius-5 vectors).
var dirs = [][2]int{
	{5, 0}, {4, 3}, {3, 4}, {0, 5}, {-3, 4}, {-4, 3},
	{-5, 0}, {-4, -3}, {-3, -4}, {0, -5}, {3, -4}, {4, -3},
}

// fw is one firework: launches at t0 from column sx and bursts at tb at (bx,by).
type fw struct {
	t0, tb, sx, bx, by int
	col                hue
}

var stars = [][2]int{{8, 3}, {26, 6}, {50, 2}, {72, 5}, {84, 8}, {15, 11}, {60, 9}}

// skyline building tops, indexed loosely across the width.
var buildings = [][3]int{
	{0, 8, 5}, {9, 6, 4}, {16, 11, 6}, {24, 7, 4}, {31, 5, 3},
	{37, 10, 5}, {46, 6, 4}, {53, 13, 7}, {63, 8, 5}, {72, 6, 4}, {79, 11, 6},
}

func drawScene(c *renderer.Canvas) {
	c.Clear()
	for _, s := range stars {
		c.Set(s[0], s[1], star)
	}
	// skyline
	for _, b := range buildings {
		x, ht, ww := b[0], b[1], b[2]
		c.Rect(x, bottom-ht, ww, ht+(h-bottom), bldg)
		// a couple lit windows
		if ht > 6 {
			c.Set(x+1, bottom-ht+2, win)
			c.Set(x+ww-2, bottom-ht+4, win)
		}
	}
	c.HSpan(0, bottom, w, bldgLit)
}

func drawLaunch(c *renderer.Canvas, f fw, frame int) {
	span := f.tb - f.t0
	if span <= 0 {
		return
	}
	p := (frame - f.t0)
	px := f.sx + (f.bx-f.sx)*p/span
	py := bottom - (bottom-f.by)*p/span
	c.Set(px, py, white)
	c.Set(px, py+1, trail)
	c.Set(px, py+2, f.col.dim)
}

func drawBurst(c *renderer.Canvas, f fw, age, maxAge int) {
	if age < 3 {
		c.Disc(f.bx, f.by, 1, white) // initial flash
	}
	r := 2 + age*9/maxAge
	grav := age * age / 9
	col := f.col.bright
	if age > maxAge/2 {
		col = f.col.dim
	}
	for _, d := range dirs {
		px := f.bx + d[0]*r/5
		py := f.by + d[1]*r/5 + grav
		c.Set(px, py, col)
		// a shorter inner spark for a fuller look
		c.Set(f.bx+d[0]*(r-3)/5, f.by+d[1]*(r-3)/5+grav*2/3, f.col.dim)
	}
}

// Build assembles a short show of overlapping fireworks.
func Build() animation.Scene {
	c := renderer.NewCanvas(w, h)
	var frames []animation.Frame
	add := func(ms int) {
		frames = append(frames, animation.Frame{Grid: c.Snapshot(), DurationMs: ms})
	}
	const burstDur = 14
	shows := []fw{
		{0, 9, 22, 24, 12, red},
		{7, 17, 62, 58, 9, gold},
		{18, 27, 44, 46, 6, green},
		{27, 36, 30, 28, 14, blue},
		{31, 40, 66, 68, 10, pink},
		{40, 48, 46, 48, 7, gold},
	}
	last := 0
	for _, s := range shows {
		if s.tb+burstDur > last {
			last = s.tb + burstDur
		}
	}
	for frame := 0; frame <= last; frame++ {
		drawScene(c)
		for _, s := range shows {
			switch {
			case frame >= s.t0 && frame < s.tb:
				drawLaunch(c, s, frame)
			case frame >= s.tb && frame < s.tb+burstDur:
				drawBurst(c, s, frame-s.tb, burstDur)
			}
		}
		add(55)
	}
	return animation.Scene{Name: "fireworks", Frames: frames}
}

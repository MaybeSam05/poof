package renderer

import "image/color"

// Sprite is a static piece of pixel art, parsed once from a palette + index grid.
type Sprite struct {
	W, H int
	Px   [][]color.RGBA // [y][x]; alpha 0 == transparent
}

// ParseSprite converts rows of palette-index characters into a Sprite. Each
// character maps to an index into palette: '.' and ' ' map to 0, '0'-'9' map to
// 0-9, and 'a'-'z' map to 10-35. Palette index 0 should be transparent. Rows of
// differing length are padded with transparent pixels.
func ParseSprite(palette []color.RGBA, rows []string) Sprite {
	h := len(rows)
	w := 0
	for _, r := range rows {
		if len(r) > w {
			w = len(r)
		}
	}
	px := make([][]color.RGBA, h)
	for y, r := range rows {
		px[y] = make([]color.RGBA, w)
		for x := 0; x < len(r); x++ {
			idx := paletteIndex(r[x])
			if idx >= 0 && idx < len(palette) {
				px[y][x] = palette[idx]
			}
		}
	}
	return Sprite{W: w, H: h, Px: px}
}

func paletteIndex(ch byte) int {
	switch {
	case ch >= '0' && ch <= '9':
		return int(ch - '0')
	case ch >= 'a' && ch <= 'z':
		return int(ch-'a') + 10
	default: // '.', ' ', anything else
		return 0
	}
}

// Canvas is a per-frame working buffer. Characters clear it, blit sprites and
// draw rectangles onto it, then Snapshot it into an immutable frame grid.
type Canvas struct {
	W, H int
	Px   [][]color.RGBA
}

// NewCanvas allocates a transparent canvas of the given size.
func NewCanvas(w, h int) *Canvas {
	px := make([][]color.RGBA, h)
	for y := range px {
		px[y] = make([]color.RGBA, w)
	}
	return &Canvas{W: w, H: h, Px: px}
}

// Clear resets every pixel to transparent.
func (c *Canvas) Clear() {
	for y := range c.Px {
		for x := range c.Px[y] {
			c.Px[y][x] = color.RGBA{}
		}
	}
}

// Fill sets every pixel to col.
func (c *Canvas) Fill(col color.RGBA) {
	for y := range c.Px {
		for x := range c.Px[y] {
			c.Px[y][x] = col
		}
	}
}

// Set writes one pixel, ignoring out-of-bounds coordinates.
func (c *Canvas) Set(x, y int, col color.RGBA) {
	if x < 0 || y < 0 || x >= c.W || y >= c.H {
		return
	}
	c.Px[y][x] = col
}

// Blit draws s at (ox, oy), skipping transparent pixels so shapes overlay cleanly.
func (c *Canvas) Blit(s Sprite, ox, oy int) {
	for y := 0; y < s.H; y++ {
		for x := 0; x < s.W; x++ {
			p := s.Px[y][x]
			if p.A == 0 {
				continue
			}
			c.Set(ox+x, oy+y, p)
		}
	}
}

// Rect fills a w×h rectangle with col (used for beams, foam, speed lines, sparks).
func (c *Canvas) Rect(x, y, w, h int, col color.RGBA) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			c.Set(xx, yy, col)
		}
	}
}

// Disc fills a circle of radius r centered at (cx, cy). The slightly padded
// radius test gives a rounder edge, which reads well for wheels and domes.
func (c *Canvas) Disc(cx, cy, r int, col color.RGBA) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r+r {
				c.Set(cx+x, cy+y, col)
			}
		}
	}
}

// HSpan fills a horizontal run from x to x+w-1 at row y (a 1px line).
func (c *Canvas) HSpan(x, y, w int, col color.RGBA) {
	c.Rect(x, y, w, 1, col)
}

// Snapshot returns a deep copy of the current pixels as an immutable Grid, so a
// reused Canvas can safely produce many frames.
func (c *Canvas) Snapshot() Grid {
	g := make(Grid, c.H)
	for y := range c.Px {
		row := make([]color.RGBA, c.W)
		copy(row, c.Px[y])
		g[y] = row
	}
	return g
}

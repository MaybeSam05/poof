package renderer

import (
	"image/color"
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

func TestParseSprite(t *testing.T) {
	pal := []color.RGBA{
		{},                // 0 transparent
		{10, 20, 30, 255}, // 1
		{40, 50, 60, 255}, // 2
	}
	s := ParseSprite(pal, []string{
		".1.",
		"22", // shorter row → padded with transparent
	})
	if s.W != 3 || s.H != 2 {
		t.Fatalf("got %dx%d, want 3x2", s.W, s.H)
	}
	if s.Px[0][0].A != 0 {
		t.Errorf("'.' should be transparent, got %+v", s.Px[0][0])
	}
	if s.Px[0][1] != pal[1] {
		t.Errorf("'1' should map to palette[1], got %+v", s.Px[0][1])
	}
	if s.Px[1][0] != pal[2] {
		t.Errorf("'2' should map to palette[2], got %+v", s.Px[1][0])
	}
	if s.Px[1][2].A != 0 {
		t.Errorf("padded cell should be transparent, got %+v", s.Px[1][2])
	}
}

func TestCanvasBlitClipsAndSkipsTransparent(t *testing.T) {
	c := NewCanvas(3, 3)
	red := color.RGBA{255, 0, 0, 255}
	s := Sprite{W: 2, H: 2, Px: [][]color.RGBA{
		{red, {}},
		{{}, {}},
	}}
	// Off the top-left corner: the opaque pixel maps to (-1,-1) and is clipped.
	c.Blit(s, -1, -1)
	for y := range c.Px {
		for x := range c.Px[y] {
			if c.Px[y][x].A != 0 {
				t.Fatalf("expected canvas untouched at (%d,%d), got %+v", x, y, c.Px[y][x])
			}
		}
	}
	c.Blit(s, 0, 0)
	if c.Px[0][0] != red {
		t.Errorf("expected red at (0,0), got %+v", c.Px[0][0])
	}
	if c.Px[0][1].A != 0 {
		t.Errorf("transparent sprite pixel must not overwrite (1,0)")
	}
}

var (
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	clear = color.RGBA{}
)

func withTrueColor(t *testing.T) {
	t.Helper()
	orig := profile
	profile = termenv.TrueColor
	t.Cleanup(func() { profile = orig })
}

func TestQuadCellFull(t *testing.T) {
	withTrueColor(t)
	out := quadCell(red, red, red, red)
	if !strings.Contains(out, "█") {
		t.Errorf("all-opaque cell should be a full block, got %q", out)
	}
	if !strings.Contains(out, "38;2;255;0;0") {
		t.Errorf("expected 24-bit red, got %q", out)
	}
}

func TestQuadCellCorner(t *testing.T) {
	withTrueColor(t)
	// Only top-left opaque → the ▘ quadrant.
	out := quadCell(red, clear, clear, clear)
	if !strings.Contains(out, "▘") {
		t.Errorf("expected top-left quadrant glyph, got %q", out)
	}
}

func TestQuadCellBlank(t *testing.T) {
	withTrueColor(t)
	if got := quadCell(clear, clear, clear, clear); !strings.HasSuffix(got, " ") {
		t.Errorf("fully transparent cell should be blank, got %q", got)
	}
}

func TestQuadCellTwoColor(t *testing.T) {
	withTrueColor(t)
	// Top row red, bottom row blue, fully opaque → ▀ with red fg, blue bg.
	out := quadCell(red, red, blue, blue)
	if !strings.Contains(out, "▀") {
		t.Errorf("expected upper-half glyph, got %q", out)
	}
	if !strings.Contains(out, "38;2;255;0;0") || !strings.Contains(out, "48;2;0;0;255") {
		t.Errorf("expected red fg + blue bg, got %q", out)
	}
}

func TestRenderPositionsAndClips(t *testing.T) {
	withTrueColor(t)
	// 4x4 grid → 2x2 cells. In an 80x24 terminal it should be centered
	// horizontally and anchored near the bottom.
	g := make(Grid, 4)
	for y := range g {
		g[y] = []color.RGBA{red, red, red, red}
	}
	out := Render(g, 80, 24)
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("expected cursor positioning in output")
	}
	// 2 cell-rows anchored 2 rows above the bottom → rows 21 and 22.
	if !strings.Contains(out, "\x1b[21;") || !strings.Contains(out, "\x1b[22;") {
		t.Errorf("expected rows 21 and 22 near the bottom, got %q", out)
	}
}

// Package renderer converts pixel grids into ANSI block-glyph strings and holds
// the low-level terminal control helpers used by the animation player.
//
// Rendering uses the Unicode quadrant block glyphs (▘▝▀▖▌▞▛▗▚▐▜▄▙▟█), which pack
// a 2×2 grid of sub-pixels into a single terminal cell. That gives four pixels
// per character cell — twice the resolution of half-blocks in each axis — so the
// art reads as fine pixels rather than chunky blocks. Each cell can carry two
// colors (a foreground for the "on" sub-pixels and a background for the rest);
// with the minimalist palettes used here that is plenty.
package renderer

import (
	"fmt"
	"image/color"
	"sort"
	"strings"

	"github.com/muesli/termenv"
)

// Grid is a pixel grid indexed as [y][x]. An alpha of 0 means transparent, and
// the terminal background shows through.
type Grid [][]color.RGBA

// profile is detected once at startup: TrueColor → 24-bit, else termenv
// downsamples, and Ascii → uncolored glyphs.
var profile = termenv.ColorProfile()

// bottomMargin leaves a couple of blank rows below the scene so it rests in the
// lower portion of the terminal rather than dead center.
const bottomMargin = 2

// quadGlyph maps a 4-bit sub-pixel mask (TL=1, TR=2, BL=4, BR=8) to its glyph.
var quadGlyph = [16]rune{
	' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛',
	'▗', '▚', '▐', '▜', '▄', '▙', '▟', '█',
}

// Render converts g into a positioned block of glyphs that overlays the current
// terminal content, centered horizontally and anchored near the bottom. It does
// not clear anything outside its own box: each row is placed with an absolute
// cursor move and every cell is repainted each frame, so moving sprites leave no
// trails and surrounding text stays visible.
func Render(g Grid, cols, rows int) string {
	gh := len(g)
	if gh == 0 || cols <= 0 || rows <= 0 {
		return ""
	}
	gw := 0
	for _, r := range g {
		if len(r) > gw {
			gw = len(r)
		}
	}
	if gw == 0 {
		return ""
	}

	cellCols := (gw + 1) / 2
	cellRows := (gh + 1) / 2

	drawCols := cellCols
	if drawCols > cols {
		drawCols = cols
	}
	drawRows := cellRows
	if drawRows > rows {
		drawRows = rows
	}

	anchorCol := (cols-drawCols)/2 + 1
	if anchorCol < 1 {
		anchorCol = 1
	}
	anchorRow := rows - drawRows - bottomMargin + 1
	if anchorRow < 1 {
		anchorRow = 1
	}

	var b strings.Builder
	for cy := 0; cy < drawRows; cy++ {
		fmt.Fprintf(&b, "\x1b[%d;%dH", anchorRow+cy, anchorCol)
		for cx := 0; cx < drawCols; cx++ {
			tl := at(g, 2*cx, 2*cy)
			tr := at(g, 2*cx+1, 2*cy)
			bl := at(g, 2*cx, 2*cy+1)
			br := at(g, 2*cx+1, 2*cy+1)
			b.WriteString(quadCell(tl, tr, bl, br))
		}
		b.WriteString("\x1b[0m")
	}
	return b.String()
}

func at(g Grid, x, y int) color.RGBA {
	if y < 0 || y >= len(g) {
		return color.RGBA{}
	}
	row := g[y]
	if x < 0 || x >= len(row) {
		return color.RGBA{}
	}
	return row[x]
}

// quadCell renders one terminal cell from its four sub-pixels. It picks up to two
// representative colors: the "on" foreground and, when the cell is fully opaque
// with exactly two colors, a background. Otherwise transparent sub-pixels become
// the (default) background and any extra colors merge into the foreground — a
// negligible loss with clean, limited palettes.
func quadCell(tl, tr, bl, br color.RGBA) string {
	subs := [4]color.RGBA{tl, tr, bl, br}
	counts := map[color.RGBA]int{}
	masks := map[color.RGBA]int{}
	order := make([]color.RGBA, 0, 4)
	anyTransparent := false
	for i, c := range subs {
		if c.A == 0 {
			anyTransparent = true
			continue
		}
		if _, ok := counts[c]; !ok {
			order = append(order, c)
		}
		counts[c]++
		masks[c] |= 1 << i
	}

	if len(order) == 0 {
		return "\x1b[0m " // fully transparent → blank cell
	}
	// Most frequent color first (stable to keep rendering deterministic).
	sort.SliceStable(order, func(a, b int) bool { return counts[order[a]] > counts[order[b]] })

	var fg color.RGBA
	var bg *color.RGBA
	var mask int
	switch {
	case len(order) == 1:
		fg = order[0]
		if anyTransparent {
			mask = masks[fg]
		} else {
			mask = 0b1111
		}
	case len(order) == 2 && !anyTransparent:
		fg = order[0]
		second := order[1]
		bg = &second
		mask = masks[fg]
	default:
		fg = order[0]
		for _, c := range order {
			mask |= masks[c]
		}
	}

	glyph := string(quadGlyph[mask])
	if profile == termenv.Ascii {
		return glyph
	}
	if bg != nil {
		return "\x1b[0m\x1b[" + seq(fg, false) + ";" + seq(*bg, true) + "m" + glyph
	}
	return "\x1b[0m\x1b[" + seq(fg, false) + "m" + glyph
}

// seq returns the SGR parameters (without the leading ESC[ and trailing m) for c.
func seq(c color.RGBA, bg bool) string {
	col := profile.Color(fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B))
	if col == nil {
		return ""
	}
	return col.Sequence(bg)
}

package renderer

import (
	"os"

	"golang.org/x/term"
)

// IsTTY reports whether stdout is connected to a terminal. When it is not (piped
// or redirected), callers should skip the animation entirely.
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Size returns the terminal size in (cols, rows), falling back to 80x24.
func Size() (cols, rows int) {
	c, r, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || c <= 0 || r <= 0 {
		return 80, 24
	}
	return c, r
}

// SaveCursor stores the current cursor position and attributes so they can be
// restored after the animation overlays the screen.
func SaveCursor() { os.Stdout.WriteString("\x1b7") }

// RestoreCursor returns the cursor to the position saved by SaveCursor.
func RestoreCursor() { os.Stdout.WriteString("\x1b8") }

// HideCursor hides the cursor during playback.
func HideCursor() { os.Stdout.WriteString("\x1b[?25l") }

// ShowCursor makes the cursor visible again.
func ShowCursor() { os.Stdout.WriteString("\x1b[?25h") }

// ClearScreen clears the current screen and homes the cursor.
func ClearScreen() { os.Stdout.WriteString("\x1b[2J\x1b[H") }

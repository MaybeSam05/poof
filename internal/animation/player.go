package animation

import (
	"bufio"
	"bytes"
	"image/color"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samarthverma/poof/internal/renderer"
)

// Options controls playback.
type Options struct {
	// Speed multiplies playback rate: 1.0 is normal, <1 slower, >1 faster.
	Speed float64
	// Preview clears the screen first, giving a clean backdrop for sprite dev.
	// Normal playback overlays the animation on top of the current content.
	Preview bool
}

// Play renders every frame of s to stdout with buffered, single-write flushes to
// avoid tearing. The animation overlays the current screen in a contained box
// rather than taking over the terminal; the wrapping `clear` wipes the leftover
// box afterwards. The cursor is always restored on exit, including on Ctrl-C.
func Play(s Scene, opts Options) {
	if opts.Speed <= 0 {
		opts.Speed = 1.0
	}
	cols, rows := renderer.Size()

	// A transparent grid matching the scene's dimensions. Rendering it overprints
	// every cell of the overlay box with a blank space using the same positioning
	// as a real frame, so it erases the box cleanly on exit (no leftover pixels).
	blank := blankGrid(s)

	restored := false
	restore := func() {
		if restored {
			return
		}
		restored = true
		os.Stdout.WriteString(renderer.Render(blank, cols, rows))
		renderer.RestoreCursor()
		renderer.ShowCursor()
	}

	// Guarantee the cursor is restored even if the process is interrupted.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		restore()
		os.Exit(130)
	}()
	defer func() {
		signal.Stop(sigCh)
		restore()
	}()

	if opts.Preview {
		renderer.ClearScreen()
	}
	renderer.SaveCursor()
	renderer.HideCursor()

	w := bufio.NewWriter(os.Stdout)
	var buf bytes.Buffer
	for _, f := range s.Frames {
		buf.Reset()
		buf.WriteString(renderer.Render(f.Grid, cols, rows))
		w.Write(buf.Bytes())
		w.Flush()

		d := time.Duration(float64(f.DurationMs)/opts.Speed) * time.Millisecond
		time.Sleep(d)
	}
}

// blankGrid returns an all-transparent grid the same size as the scene's frames.
func blankGrid(s Scene) renderer.Grid {
	if len(s.Frames) == 0 || len(s.Frames[0].Grid) == 0 {
		return nil
	}
	h := len(s.Frames[0].Grid)
	w := len(s.Frames[0].Grid[0])
	g := make(renderer.Grid, h)
	for y := range g {
		g[y] = make([]color.RGBA, w)
	}
	return g
}

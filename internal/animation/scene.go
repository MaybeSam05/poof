// Package animation defines the frame/scene model and the playback loop.
package animation

import "github.com/samarthverma/poof/internal/renderer"

// Frame is one rendered pixel grid shown for DurationMs milliseconds.
type Frame struct {
	Grid       renderer.Grid
	DurationMs int
}

// Scene is a named sequence of frames produced by a character's Build function.
type Scene struct {
	Name   string
	Frames []Frame
}

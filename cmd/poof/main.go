// Command poof plays a short pixel-art animation, intended to run just before
// the real `clear`. A shell function wraps `clear` to call poof first.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/characters"
	"github.com/samarthverma/poof/internal/renderer"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	var (
		character = flag.String("character", "", "character to play (default: random)")
		list      = flag.Bool("list", false, "list available characters and exit")
		preview   = flag.Bool("preview", false, "clear the screen first, then play on a clean backdrop (for sprite dev)")
		speed     = flag.Float64("speed", 1.0, "speed multiplier (<1 slower, >1 faster)")
		showVer   = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	if *showVer {
		fmt.Println("poof", version)
		return
	}

	if *list {
		fmt.Println("Available characters:")
		for _, n := range characters.Names() {
			fmt.Println("  " + n)
		}
		return
	}

	// If output is not a terminal (piped/redirected) and we're not previewing,
	// do nothing so the wrapping shell function's `command clear` still works.
	if !*preview && !renderer.IsTTY() {
		return
	}

	scene, ok := pickScene(*character)
	if !ok {
		fmt.Fprintf(os.Stderr, "poof: unknown character %q (try --list)\n", *character)
		os.Exit(1)
	}

	animation.Play(scene, animation.Options{
		Speed:   *speed,
		Preview: *preview,
	})
}

func pickScene(name string) (animation.Scene, bool) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return characters.Random(), true
	}
	return characters.Get(name)
}

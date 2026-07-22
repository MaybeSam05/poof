// Command poof plays a short pixel-art animation, intended to run just before
// the real `clear`. A shell function wraps `clear` to call `poof play` first.
//
// Friendly commands (all persist to the config file):
//
//	poof car          set the default character (car|surf|alien|random|rotate)
//	poof car 0.5      set character and speed
//	poof 0.5          set just the speed
//	poof disable      stop animating on clear (poof enable to undo)
//	poof preview car  play once without clearing, to try it out
//	poof status       show current settings
//	poof list         list characters
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/characters"
	"github.com/samarthverma/poof/internal/renderer"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	args := os.Args[1:]
	cfg := loadConfig()

	// Play mode — no args or explicit "play" (this is what the clear hook runs).
	if len(args) == 0 || strings.EqualFold(args[0], "play") {
		if !cfg.Enabled || !renderer.IsTTY() {
			return
		}
		ch := cfg.Character
		if isRotate(ch) {
			ch = advanceRotate(&cfg)
			persist(cfg)
		}
		play(ch, cfg.Speed, false)
		return
	}

	// Legacy flag form (e.g. --preview --character f1 --speed 0.5).
	if strings.HasPrefix(args[0], "-") {
		runFlags(cfg)
		return
	}

	switch cmd := strings.ToLower(args[0]); cmd {
	case "disable", "off", "stop":
		cfg.Enabled = false
		persist(cfg)
		fmt.Println("poof disabled — clear won't animate. Run `poof enable` to turn it back on.")
	case "enable", "on", "start":
		cfg.Enabled = true
		persist(cfg)
		fmt.Println("poof enabled.")
	case "status", "config", "settings":
		printStatus(cfg)
	case "list", "chars", "characters":
		printList()
	case "uninstall", "remove":
		removed := uninstall()
		if len(removed) == 0 {
			fmt.Println("poof doesn't appear to be installed.")
			return
		}
		fmt.Println("poof uninstalled:")
		for _, r := range removed {
			fmt.Println("  removed " + r)
		}
		fmt.Println("Restart your terminal to finish (clear goes back to normal).")
	case "version":
		fmt.Println("poof", version)
	case "help", "h":
		printHelp()
	case "preview", "test", "demo":
		ch, sp := cfg.Character, cfg.Speed
		for _, a := range args[1:] {
			if c, ok := resolveChar(a); ok {
				ch = c
			} else if s, err := strconv.ParseFloat(a, 64); err == nil && s > 0 {
				sp = s
			}
		}
		if isRotate(ch) {
			ch = advanceRotate(&cfg)
			persist(cfg)
		}
		play(ch, sp, true)
	default:
		// A character alias and/or a speed number, in any order.
		changed := false
		for _, a := range args {
			if c, ok := resolveChar(a); ok {
				if cfg.Character != c {
					cfg.RotateIdx = 0
				}
				cfg.Character = c
				changed = true
			} else if s, err := strconv.ParseFloat(a, 64); err == nil && s > 0 {
				cfg.Speed = s
				changed = true
			}
		}
		if !changed {
			fmt.Fprintf(os.Stderr, "poof: didn't understand %q\n\n", strings.Join(args, " "))
			printHelp()
			os.Exit(1)
		}
		persist(cfg)
		fmt.Printf("Set: character=%s, speed=%sx. (Run `poof preview` to see it.)\n",
			cfg.Character, trimFloat(cfg.Speed))
	}
}

// play builds and plays the given character, or a random one for random/unknown.
func play(character string, speed float64, preview bool) {
	scene, ok := characters.Get(character)
	if character == "" || character == "random" || isRotate(character) || !ok {
		scene = characters.Random()
	}
	animation.Play(scene, animation.Options{Speed: speed, Preview: preview})
}

func isRotate(character string) bool {
	return character == "rotate" || character == "cycle"
}

func advanceRotate(cfg *config) string {
	name, next := nextRotateCharacter(cfg.RotateIdx)
	cfg.RotateIdx = next
	return name
}

func nextRotateCharacter(index int) (string, int) {
	names := characters.Names()
	if len(names) == 0 {
		return "random", 0
	}
	name := names[index%len(names)]
	return name, (index + 1) % len(names)
}

// runFlags preserves the original --flag interface.
func runFlags(cfg config) {
	fs := flag.NewFlagSet("poof", flag.ExitOnError)
	character := fs.String("character", "", "character to play")
	list := fs.Bool("list", false, "list characters")
	preview := fs.Bool("preview", false, "play without clearing")
	speed := fs.Float64("speed", cfg.Speed, "speed multiplier")
	showVer := fs.Bool("version", false, "print version")
	_ = fs.Parse(os.Args[1:])

	switch {
	case *showVer:
		fmt.Println("poof", version)
	case *list:
		printList()
	default:
		ch := cfg.Character
		if *character != "" {
			if c, ok := resolveChar(*character); ok {
				ch = c
			} else {
				ch = *character
			}
		}
		if *preview {
			if isRotate(ch) {
				ch = advanceRotate(&cfg)
				persist(cfg)
			}
			play(ch, *speed, true)
			return
		}
		if cfg.Enabled && renderer.IsTTY() {
			if isRotate(ch) {
				ch = advanceRotate(&cfg)
				persist(cfg)
			}
			play(ch, *speed, false)
		}
	}
}

func persist(cfg config) {
	if err := saveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "poof: could not save config: %v\n", err)
	}
}

func printStatus(cfg config) {
	state := "enabled"
	if !cfg.Enabled {
		state = "disabled"
	}
	fmt.Printf("poof is %s\n  character: %s\n  speed:     %sx\n  config:    %s\n",
		state, cfg.Character, trimFloat(cfg.Speed), configPath())
}

func printList() {
	fmt.Println("Characters:")
	for _, n := range characters.Names() {
		fmt.Println("  " + n)
	}
	fmt.Println("Aliases: car=f1, wave=surf, ufo=alien, launch=rocket, trex=dino,")
	fmt.Println("         boom=fireworks, choo=train, heli=helicopter, jaws=shark, random, rotate")
}

func printHelp() {
	fmt.Print(`poof — a pixel-art animation before your terminal clears

Usage:
  poof car            set the default character (car | surf | alien | rocket |
                      dino | fireworks | train | helicopter | shark | random |
                      rotate)
  poof car 0.5        set character and speed (0.5 = half speed, 2 = double)
  poof 0.5            set just the speed
  poof disable        stop animating on clear   (poof enable turns it back on)
  poof preview car    play once now, without clearing
  poof status         show current settings
  poof list           list characters
  poof uninstall      remove poof completely
  poof version

Characters: car, surf, alien, rocket, dino, fireworks, train, helicopter, shark, random, rotate
`)
}

func trimFloat(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

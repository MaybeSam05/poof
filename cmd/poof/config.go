package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// config holds the persisted user preferences that drive the clear hook.
type config struct {
	Enabled   bool
	Character string // "random" or a character name
	Speed     float64
}

func defaultConfig() config {
	return config{Enabled: true, Character: "random", Speed: 1.0}
}

func configPath() string {
	dir, err := os.UserConfigDir()
	if err != nil || dir == "" {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(dir, "poof", "config")
}

func loadConfig() config {
	cfg := defaultConfig()
	f, err := os.Open(configPath())
	if err != nil {
		return cfg
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		switch k {
		case "enabled":
			cfg.Enabled = v == "true" || v == "1" || v == "yes" || v == "on"
		case "character":
			cfg.Character = v
		case "speed":
			if s, err := strconv.ParseFloat(v, 64); err == nil && s > 0 {
				cfg.Speed = s
			}
		}
	}
	return cfg
}

func saveConfig(cfg config) error {
	p := configPath()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	content := fmt.Sprintf("enabled=%t\ncharacter=%s\nspeed=%s\n",
		cfg.Enabled, cfg.Character, strconv.FormatFloat(cfg.Speed, 'g', -1, 64))
	return os.WriteFile(p, []byte(content), 0o644)
}

// stripPoofBlock removes the guarded "# >>> poof >>>" … "# <<< poof <<<" block
// from an rc file. It reports whether a block was found and removed.
func stripPoofBlock(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines))
	skip, changed := false, false
	for _, ln := range lines {
		switch strings.TrimSpace(ln) {
		case "# >>> poof >>>":
			skip, changed = true, true
			continue
		case "# <<< poof <<<":
			skip = false
			continue
		}
		if skip {
			continue
		}
		out = append(out, ln)
	}
	if !changed {
		return false
	}
	// collapse a blank line left where the block was
	joined := strings.ReplaceAll(strings.Join(out, "\n"), "\n\n\n", "\n\n")
	return os.WriteFile(path, []byte(joined), 0o644) == nil
}

// uninstall reverses the installer: strips the shell block from rc files, and
// removes the binary and config. Returns the list of things removed.
func uninstall() []string {
	home, _ := os.UserHomeDir()
	var removed []string
	for _, rc := range []string{".bashrc", ".zshrc", ".bash_profile", ".profile"} {
		if stripPoofBlock(filepath.Join(home, rc)) {
			removed = append(removed, "shell hook in ~/"+rc)
		}
	}
	if os.RemoveAll(filepath.Dir(configPath())) == nil {
		removed = append(removed, "settings ("+filepath.Dir(configPath())+")")
	}
	bin := filepath.Join(home, ".local", "bin", "poof")
	if err := os.Remove(bin); err == nil {
		removed = append(removed, bin)
	}
	return removed
}

// resolveChar maps friendly names/aliases to a canonical character.
func resolveChar(s string) (string, bool) {
	switch strings.ToLower(s) {
	case "f1", "car", "formula1", "formula", "race", "racecar":
		return "f1", true
	case "surf", "surfer", "wave", "surfing":
		return "surf", true
	case "alien", "ufo", "saucer":
		return "alien", true
	case "rocket", "launch", "blastoff", "rocketship", "spaceship":
		return "rocket", true
	case "dino", "dinosaur", "trex", "t-rex", "rex", "raptor":
		return "dino", true
	case "fireworks", "firework", "fw", "boom", "celebrate":
		return "fireworks", true
	case "train", "locomotive", "steam", "choochoo", "choo":
		return "train", true
	case "helicopter", "heli", "chopper", "copter":
		return "helicopter", true
	case "shark", "jaws", "fin":
		return "shark", true
	case "random", "rand", "any", "surprise":
		return "random", true
	}
	return "", false
}

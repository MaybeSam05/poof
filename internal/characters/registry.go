// Package characters maps character names to their scene builders and provides
// lookup, listing, and random selection.
package characters

import (
	"sort"
	"time"

	"github.com/samarthverma/poof/internal/animation"
	"github.com/samarthverma/poof/internal/characters/alien"
	"github.com/samarthverma/poof/internal/characters/dino"
	"github.com/samarthverma/poof/internal/characters/f1"
	"github.com/samarthverma/poof/internal/characters/fireworks"
	"github.com/samarthverma/poof/internal/characters/helicopter"
	"github.com/samarthverma/poof/internal/characters/rocket"
	"github.com/samarthverma/poof/internal/characters/shark"
	"github.com/samarthverma/poof/internal/characters/surf"
	"github.com/samarthverma/poof/internal/characters/train"
)

// Builder produces a fresh Scene each call.
type Builder func() animation.Scene

var registry = map[string]Builder{
	"alien":      alien.Build,
	"surf":       surf.Build,
	"f1":         f1.Build,
	"rocket":     rocket.Build,
	"dino":       dino.Build,
	"fireworks":  fireworks.Build,
	"train":      train.Build,
	"helicopter": helicopter.Build,
	"shark":      shark.Build,
}

// Names returns the available character names, sorted.
func Names() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Get builds the named scene, reporting whether the name exists.
func Get(name string) (animation.Scene, bool) {
	b, ok := registry[name]
	if !ok {
		return animation.Scene{}, false
	}
	return b(), true
}

// Random builds a randomly chosen scene.
func Random() animation.Scene {
	names := Names()
	i := int(time.Now().UnixNano()) % len(names)
	if i < 0 {
		i = -i
	}
	return registry[names[i]]()
}

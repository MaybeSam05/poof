package main

import (
	"testing"

	"github.com/samarthverma/poof/internal/characters"
)

func TestNextRotateCharacterDoesNotRepeatBeforeFullRotation(t *testing.T) {
	names := characters.Names()
	seen := make(map[string]bool, len(names))
	idx := 0

	for range names {
		var name string
		name, idx = nextRotateCharacter(idx)
		if seen[name] {
			t.Fatalf("%s repeated before all characters were seen", name)
		}
		seen[name] = true
	}

	name, _ := nextRotateCharacter(idx)
	if name != names[0] {
		t.Fatalf("rotation did not wrap to %s, got %s", names[0], name)
	}
}

package characters

import "testing"

func TestRegistryScenesAreWellFormed(t *testing.T) {
	names := Names()
	if len(names) == 0 {
		t.Fatal("expected at least one character")
	}
	for _, name := range names {
		s, ok := Get(name)
		if !ok {
			t.Fatalf("Get(%q) reported missing", name)
		}
		if len(s.Frames) == 0 {
			t.Errorf("%s: no frames", name)
			continue
		}
		w := len(s.Frames[0].Grid[0])
		h := len(s.Frames[0].Grid)
		for i, f := range s.Frames {
			if f.DurationMs <= 0 {
				t.Errorf("%s frame %d: non-positive duration %d", name, i, f.DurationMs)
			}
			if len(f.Grid) != h || len(f.Grid[0]) != w {
				t.Errorf("%s frame %d: inconsistent size %dx%d, want %dx%d",
					name, i, len(f.Grid[0]), len(f.Grid), w, h)
			}
		}
	}
}

func TestGetUnknownReturnsFalse(t *testing.T) {
	if _, ok := Get("nope"); ok {
		t.Error("expected unknown character to report ok=false")
	}
}

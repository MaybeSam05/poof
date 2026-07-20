# poof

Make clearing your terminal fun. `poof` plays a short pixel-art animation *on top
of* your current terminal content — a contained box overlaid in place, not a
full-screen takeover — and then the real `clear` wipes the screen.

```
you type:  clear   (or cls, or press Ctrl+L)
              │
   shell function runs poof  ──►  a random animation (alien / surf / f1) poofs
              │                    on top of the screen, surrounding text still visible
        command clear         ──►  the real clear runs afterwards
```

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/samarthverma/poof/main/install.sh | bash
```

This downloads the right binary for your OS/arch to `~/.local/bin/poof` and adds a
`# >>> poof >>>` block to your `~/.bashrc` or `~/.zshrc`. Restart your terminal
(or `source` the rc file), then type `clear`.

To remove it, delete the `# >>> poof >>>` block from your rc file and remove
`~/.local/bin/poof`.

## Triggers

The installer wires the animation into the usual ways you clear the screen:

| trigger        | what it does                                    |
| -------------- | ----------------------------------------------- |
| `clear`        | shell function → poof, then `command clear`     |
| `cls`          | same, for the common `cls` alias                |
| `Ctrl+L`       | the readline/zle clear-screen shortcut          |

All of them call `poof_clear`, so the animation always runs *before* the real
clear. Want more? Any command that ends by clearing the screen is a good fit —
e.g. add `reset() { poof_clear; command reset; }` to the block, or wrap your own
alias. To opt a trigger out, just delete its line from the `# >>> poof >>>` block.

## Characters

| name    | scene                                                    |
| ------- | -------------------------------------------------------- |
| `alien` | a sleek UFO cruises in, pulses a tractor beam, and zips off |
| `surf`  | a surfer rides a swell that rolls in and breaks into foam |
| `f1`    | a Formula 1 car sweeps across trailing motion streaks    |

The animations are minimalist, rendered as fine pixels (Unicode quadrant blocks)
and rest in the lower portion of the terminal. By default a random one plays.

## Usage

```
poof                     # random character (a bare run does NOT clear)
poof --character alien    # play a specific character
poof --list               # list available characters
poof --preview            # clear first, then play on a clean backdrop (sprite dev)
poof --speed 0.5          # speed multiplier: <1 slower, >1 faster (default 1.0)
poof --version
```

`poof` never runs `clear` itself — the shell function does. It overlays the
animation on top of whatever is on screen and leaves it there; the wrapping
`clear` cleans up. `--preview` clears first so you can inspect sprites on a blank
backdrop.

If stdout isn't a terminal (piped or redirected), `poof` does nothing and exits
0, so scripts that call `clear` keep working.

## Build from source

```sh
go build -o poof ./cmd/poof
go test ./...
./poof --preview --character f1
```

## How it adds a new character

Each character lives in `internal/characters/<name>/` and exposes
`func Build() animation.Scene`. Sprites are defined once as a palette + index
grid and composed into frames programmatically (position/state computed per
frame) using the `renderer.Canvas` helpers (`Blit`, `Rect`). Register the
builder in `internal/characters/registry.go` and it shows up everywhere.

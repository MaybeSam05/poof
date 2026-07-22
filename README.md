<div align="center">
  <img src="assets/poof_logo.png" alt="Poof logo" width="236">

  <h1>Poof</h1>

  <p>
    A little pixel-art animation plays in your terminal every time you clear the screen.
  </p>

  <p>
    <a href="#one-line-install"><strong>⬇️ Install Poof</strong></a>
  </p>

  <p><strong>Requires macOS or Linux</strong></p>

  <img src="assets/poof_video.gif" alt="Poof demo" width="720">

    dino - car - fireworks - helicopter - shark - surf - train - alien

</div>

## One-line Install

```sh
curl -fsSL https://raw.githubusercontent.com/MaybeSam05/poof/main/install.sh | bash
```

Then restart your terminal (or run `source ~/.bashrc`). Type `clear` and enjoy.

## Use it

After install, `clear`, `cls`, and `Ctrl+L` play the animation before clearing your terminal.

Set one animation:

```sh
poof car
poof surf
poof alien
poof rocket
poof dino
poof fireworks
poof train
poof helicopter
poof shark
```

Change how animations are picked:

```sh
poof random   # pick one at random each time
poof rotate   # cycle through every animation in order
```

Change speed:

```sh
poof car 0.5  # slower
poof car 2    # faster
```

Manage Poof:

```sh
poof preview  # play the current animation now
poof status   # show your current settings
poof disable  # turn animations off
poof enable   # turn animations back on
```

By default a random one plays. Use `poof rotate` to cycle through every animation in order.
Settings are remembered in `~/.config/poof/config`.

## Characters

| Preview | Command |
| --- | --- |
| <img src="assets/characters/alien.png" alt="Alien UFO character" width="360"> | `poof alien` |
| <img src="assets/characters/dino.png" alt="Dino character" width="360"> | `poof dino` |
| <img src="assets/characters/f1.png" alt="F1 car character" width="360"> | `poof car` |
| <img src="assets/characters/fireworks.png" alt="Fireworks character" width="360"> | `poof fireworks` |
| <img src="assets/characters/helicopter.png" alt="Helicopter character" width="360"> | `poof helicopter` |
| <img src="assets/characters/rocket.png" alt="Rocket character" width="360"> | `poof rocket` |
| <img src="assets/characters/shark.png" alt="Shark character" width="360"> | `poof shark` |
| <img src="assets/characters/surf.png" alt="Surf character" width="360"> | `poof surf` |
| <img src="assets/characters/train.png" alt="Train character" width="360"> | `poof train` |

## Uninstall

```sh
poof uninstall
```

---

Build from source

Requires [Go](https://go.dev/dl/) 1.22+.

```sh
git clone https://github.com/MaybeSam05/poof && cd poof
go build -o poof ./cmd/poof
./poof preview          # try it
```

Add a new character in `internal/characters/<name>/` exposing `func Build() animation.Scene`
(compose sprites with the `renderer.Canvas` helpers), then register it in
`internal/characters/registry.go`.

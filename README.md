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
</div>

## One-line Install

```sh
curl -fsSL https://raw.githubusercontent.com/MaybeSam05/poof/main/install.sh | bash
```

Then restart your terminal (or run `source ~/.bashrc`). Type `clear` and enjoy.

## Use it

```
clear           play an animation, then clear   (also works: cls, or Ctrl+L)

poof car        always play the car             poof surf, poof alien, poof rocket,
poof car 0.5    ...at half speed (2 = faster)    dino, fireworks, train, helicopter, shark
poof disable    turn it off                      poof enable turns it back on
poof preview    watch the current one now
poof status     show your settings
```

By default a random one plays. Settings are remembered in `~/.config/poof/config`.

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

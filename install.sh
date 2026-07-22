#!/usr/bin/env bash
# poof installer — downloads the right binary and wraps `clear` to play an
# animation first. Usage: curl -fsSL https://raw.githubusercontent.com/OWNER/poof/main/install.sh | bash
set -euo pipefail

REPO="${POOF_REPO:-MaybeSam05/poof}"
BIN_DIR="${HOME}/.local/bin"
BIN="${BIN_DIR}/poof"

# --- detect platform ---------------------------------------------------------
os="$(uname -s)"
case "$os" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *) echo "poof: unsupported OS: $os" >&2; exit 1 ;;
esac

arch="$(uname -m)"
case "$arch" in
  x86_64|amd64)   ARCH="amd64" ;;
  arm64|aarch64)  ARCH="arm64" ;;
  *) echo "poof: unsupported architecture: $arch" >&2; exit 1 ;;
esac

ASSET="poof_${OS}_${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

# --- download ----------------------------------------------------------------
echo "Installing poof (${OS}/${ARCH}) from ${REPO}..."
mkdir -p "$BIN_DIR"
if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$URL" -o "$BIN"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$BIN" "$URL"
else
  echo "poof: need curl or wget to download" >&2; exit 1
fi
chmod +x "$BIN"

# --- wire up the shell -------------------------------------------------------
# Pick the rc file for the current shell.
shell_name="$(basename "${SHELL:-bash}")"
case "$shell_name" in
  zsh)  RC="${HOME}/.zshrc" ;;
  bash) RC="${HOME}/.bashrc" ;;
  *)    RC="${HOME}/.bashrc" ;;
esac

MARKER_START="# >>> poof >>>"
MARKER_END="# <<< poof <<<"

# Quoted heredoc: nothing is expanded, so the block is written verbatim.
read -r -d '' BLOCK <<'EOF' || true
# >>> poof >>>
export PATH="${HOME}/.local/bin:${PATH}"
# poof plays a short animation on top of the screen, then runs the real clear.
poof_clear() { command poof play; command clear; }
clear() { poof_clear; }
cls() { poof_clear; }
# Animate on the Ctrl+L clear-screen shortcut too.
if [ -n "${BASH_VERSION:-}" ]; then
  bind -x '"\C-l": poof_clear' 2>/dev/null || true
elif [ -n "${ZSH_VERSION:-}" ]; then
  poof-clear-widget() { command poof; command clear; zle reset-prompt; }
  zle -N poof-clear-widget
  bindkey '^L' poof-clear-widget
fi
# <<< poof <<<
EOF

touch "$RC"
if grep -qF "$MARKER_START" "$RC"; then
  # Replace the existing block so re-running the installer stays idempotent.
  tmp="$(mktemp)"
  awk -v s="$MARKER_START" -v e="$MARKER_END" '
    $0==s {skip=1}
    skip==0 {print}
    $0==e {skip=0}
  ' "$RC" > "$tmp"
  printf '%s\n' "$BLOCK" >> "$tmp"
  mv "$tmp" "$RC"
else
  printf '\n%s\n' "$BLOCK" >> "$RC"
fi

echo "Done. poof installed to ${BIN} and wired into ${RC}."
echo "Restart your terminal or run: source ${RC}"

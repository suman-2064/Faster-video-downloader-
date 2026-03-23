#!/data/data/com.termux/files/usr/bin/bash
# ─────────────────────────────────────────────────────────────────────────────
#  faster — Setup Script for Termux
#  GitHub: https://github.com/YOUR_USERNAME/faster
# ─────────────────────────────────────────────────────────────────────────────

set -e

# ── Colors ────────────────────────────────────────────────────────────────────
R="\e[0m"
BOLD="\e[1m"
DIM="\e[2m"
RED="\e[31m"
GREEN="\e[32m"
YELLOW="\e[33m"
BLUE="\e[34m"
MAGENTA="\e[35m"
CYAN="\e[36m"
WHITE="\e[37m"

ok()   { echo -e "${GREEN}${BOLD}  ✓  $*${R}"; }
err()  { echo -e "${RED}${BOLD}  ✗  $*${R}"; }
info() { echo -e "${CYAN}  ➜  $*${R}"; }
warn() { echo -e "${YELLOW}  ⚠  $*${R}"; }
head() { echo -e "${BOLD}${WHITE}$*${R}"; }
div()  { echo -e "${DIM}  ──────────────────────────────────────────────────────${R}"; }

# ── Banner ────────────────────────────────────────────────────────────────────
echo ""
echo -e "${CYAN}${BOLD}  ┌──────────────────────────────────────────────────┐"
echo -e "  │     ⚡  FASTER — Termux Downloader Setup         │"
echo -e "  │         v1.0.0  •  Powered by yt-dlp + Go        │"
echo -e "  └──────────────────────────────────────────────────┘${R}"
echo ""

# ── Root check ────────────────────────────────────────────────────────────────
if [ "$(id -u)" -eq 0 ]; then
  warn "Running as root is not recommended in Termux."
fi

# ── Termux check ─────────────────────────────────────────────────────────────
if [ -z "$PREFIX" ]; then
  err "This script must be run inside Termux."
  exit 1
fi

BIN="$PREFIX/bin"

# ═════════════════════════════════════════════════════════════════════════════
# STEP 1 — Update & Upgrade
# ═════════════════════════════════════════════════════════════════════════════
div
head "  [1/6]  Updating package lists…"
div
pkg update -y && pkg upgrade -y
ok "Packages updated"
echo ""

# ═════════════════════════════════════════════════════════════════════════════
# STEP 2 — Install System Dependencies
# ═════════════════════════════════════════════════════════════════════════════
div
head "  [2/6]  Installing system packages…"
div
PACKAGES="golang ffmpeg python git curl"
for pkg_name in $PACKAGES; do
  if command -v "$pkg_name" &>/dev/null || pkg list-installed 2>/dev/null | grep -q "^$pkg_name"; then
    ok "$pkg_name already installed"
  else
    info "Installing $pkg_name…"
    pkg install -y "$pkg_name"
    ok "$pkg_name installed"
  fi
done
echo ""

# ═════════════════════════════════════════════════════════════════════════════
# STEP 3 — Install Python Dependencies
# ═════════════════════════════════════════════════════════════════════════════
div
head "  [3/6]  Installing Python packages…"
div

# Upgrade pip first
pip install --upgrade pip --quiet
ok "pip upgraded"

# Install / upgrade yt-dlp
info "Installing yt-dlp (latest)…"
pip install --upgrade yt-dlp --quiet
ok "yt-dlp installed: $(yt-dlp --version)"
echo ""

# ═════════════════════════════════════════════════════════════════════════════
# STEP 4 — Build Go Binary
# ═════════════════════════════════════════════════════════════════════════════
div
head "  [4/6]  Compiling faster binary…"
div

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MAIN_GO="$SCRIPT_DIR/main.go"

if [ ! -f "$MAIN_GO" ]; then
  err "main.go not found in $SCRIPT_DIR"
  err "Please run this script from the project root (where main.go lives)."
  exit 1
fi

info "Running: go build -ldflags='-s -w' -o faster main.go"
cd "$SCRIPT_DIR"

# Initialise go module if needed
if [ ! -f "go.mod" ]; then
  go mod init faster
  ok "go.mod created"
fi

go build -ldflags="-s -w" -o faster main.go
ok "Binary compiled successfully"

# ═════════════════════════════════════════════════════════════════════════════
# STEP 5 — Install Binary
# ═════════════════════════════════════════════════════════════════════════════
div
head "  [5/6]  Installing faster to $BIN…"
div

chmod +x faster
cp faster "$BIN/faster"
ok "faster installed to $BIN/faster"

# ── Shell alias (optional) ────────────────────────────────────────────────────
SHELL_RC="$HOME/.bashrc"
if [ -f "$HOME/.zshrc" ]; then
  SHELL_RC="$HOME/.zshrc"
fi

ALIAS_LINE='alias faster="faster"'
if ! grep -q 'alias faster' "$SHELL_RC" 2>/dev/null; then
  echo "" >> "$SHELL_RC"
  echo "# faster — High-Speed Downloader" >> "$SHELL_RC"
  echo "$ALIAS_LINE" >> "$SHELL_RC"
  ok "Alias added to $SHELL_RC"
fi
echo ""

# ═════════════════════════════════════════════════════════════════════════════
# STEP 6 — termux-url-opener
# ═════════════════════════════════════════════════════════════════════════════
div
head "  [6/6]  Setting up Share / termux-url-opener…"
div

URL_OPENER="$HOME/.shortcuts/url-opener"
URL_OPENER_BIN="$BIN/termux-url-opener"

# Create the url-opener script (called by Termux:API when sharing a URL)
cat > "$URL_OPENER_BIN" << 'URLSCRIPT'
#!/data/data/com.termux/files/usr/bin/bash
# termux-url-opener — called automatically when you Share a URL into Termux
URL="$1"
if [ -z "$URL" ]; then
  echo "No URL provided."
  exit 1
fi
# Open a new Termux session running faster with the shared URL
am start \
  --user 0 \
  -n com.termux/com.termux.app.TermuxActivity \
  --es com.termux.ipc.extra.TERMINAL_COMMAND "faster '$URL'" \
  2>/dev/null || true

# Fallback: just run faster directly if am isn't available
faster "$URL"
URLSCRIPT

chmod +x "$URL_OPENER_BIN"
ok "termux-url-opener installed to $URL_OPENER_BIN"

# Also place it in ~/.shortcuts for Termux:Widget users
mkdir -p "$HOME/.shortcuts"
cp "$URL_OPENER_BIN" "$URL_OPENER"
chmod +x "$URL_OPENER"
ok "Shortcut placed at $URL_OPENER"

echo ""

# ═════════════════════════════════════════════════════════════════════════════
# DONE
# ═════════════════════════════════════════════════════════════════════════════
div
echo ""
echo -e "${GREEN}${BOLD}  ╔══════════════════════════════════════════════════╗"
echo -e "  ║   ✅  Installation Complete!                    ║"
echo -e "  ╚══════════════════════════════════════════════════╝${R}"
echo ""
echo -e "${CYAN}  Quick-start:${R}"
echo -e "    ${BOLD}faster${R}              — interactive mode"
echo -e "    ${BOLD}faster <URL>${R}        — download immediately"
echo -e "    ${BOLD}faster --check${R}      — verify dependencies"
echo -e "    ${BOLD}faster --help${R}       — show all options"
echo ""
echo -e "${CYAN}  Share integration:${R}"
echo -e "    Share any YouTube / Instagram / TikTok link"
echo -e "    from your browser → choose ${BOLD}Termux${R} → faster"
echo -e "    opens automatically and asks for quality."
echo ""
echo -e "${DIM}  Files are saved to: /sdcard/Download/Faster/${R}"
echo ""
info "Reload your shell:  source $SHELL_RC"
echo ""

#!/usr/bin/env bash
set -euo pipefail

DAEMON_URL="http://127.0.0.1:4318"
INSTALL_LOG="$HOME/.traceknot/install.log"
PREFIX="${PREFIX:-$HOME/.local}"
BIN_DIR="$PREFIX/bin"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ $# -gt 0 ]]; then
	echo "Unknown argument: $1 (this installer takes no arguments)" >&2
	exit 2
fi

case "$(uname -s)" in
	Darwin) tk_platform="darwin" ;;
	Linux) tk_platform="linux" ;;
	*)
		echo "Unsupported operating system: $(uname -s)" >&2
		exit 1
		;;
esac

case "$(uname -m)" in
	x86_64 | amd64) arch="x64" ;;
	arm64 | aarch64) arch="arm64" ;;
	*)
		echo "Unsupported architecture: $(uname -m)" >&2
		exit 1
		;;
esac

BINARY="$SCRIPT_DIR/downloads/traceknot-$tk_platform-$arch"

tk_log() {
	mkdir -p "$(dirname "$INSTALL_LOG")"
	printf '%s  %s\n' "$(date -u +%H:%M:%S)" "$*" >>"$INSTALL_LOG"
}

if [[ ! -f "$BINARY" ]]; then
	echo "Error: $BINARY not found next to bootstrap.sh." >&2
	exit 1
fi

mkdir -p "$BIN_DIR"
TMP_BIN="$BIN_DIR/.traceknot-install-$$"
trap 'rm -f "$TMP_BIN"' EXIT

echo "Installing traceknot..."
tk_log "install started"
cp "$BINARY" "$TMP_BIN"
chmod +x "$TMP_BIN"
mv "$TMP_BIN" "$BIN_DIR/traceknot"
tk_log "binary installed at $BIN_DIR/traceknot"

if { [[ ! -t 0 ]] || [[ ! -t 1 ]]; } && [[ -e /dev/tty ]]; then
	"$BIN_DIR/traceknot" < /dev/tty
else
	"$BIN_DIR/traceknot"
fi
tk_log "install: interactive menu completed"
tk_log "install complete"

echo ""
echo "Installation complete!"
echo "Dashboard available at: $DAEMON_URL/"
echo "Reconfigure anytime with: traceknot"
echo "Installation log: $INSTALL_LOG"

#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-https://downloads.traceknot.com}"
DAEMON_URL="http://127.0.0.1:4318"
INSTALL_LOG="$HOME/.traceknot/install.log"
PREFIX="${PREFIX:-$HOME/.local}"
BIN_DIR="$PREFIX/bin"

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

if command -v curl >/dev/null 2>&1; then
	download() {
		curl -fsSL "$1" -o "$2"
	}
elif command -v wget >/dev/null 2>&1; then
	download() {
		wget -q -O "$2" "$1"
	}
else
	echo "Either curl or wget is required but neither is installed" >&2
	exit 1
fi

tk_sha256() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | awk '{print $1}'
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$1" | awk '{print $1}'
	else
		echo "Error: sha256sum or shasum is required to verify the download but neither is installed" >&2
		exit 1
	fi
}

tk_log() {
	mkdir -p "$(dirname "$INSTALL_LOG")"
	printf '%s  %s\n' "$(date -u +%H:%M:%S)" "$*" >>"$INSTALL_LOG"
}

if [[ -n "${TRACEKNOT_VERSION:-}" ]]; then
	RELEASE="$TRACEKNOT_VERSION"
else
	TMP_LATEST="$(mktemp)"
	trap 'rm -f "$TMP_LATEST"' EXIT
	download "$BASE_URL/latest.txt" "$TMP_LATEST"
	RELEASE="$(tr -d '[:space:]' <"$TMP_LATEST")"
	rm -f "$TMP_LATEST"
	trap - EXIT
fi

if [[ -z "$RELEASE" ]]; then
	echo "Error: could not determine which release to install from $BASE_URL/latest.txt" >&2
	exit 1
fi

BINARY_URL="$BASE_URL/releases/$RELEASE/downloads/traceknot-$tk_platform-$arch"
CHECKSUM_URL="$BINARY_URL.sha256"
mkdir -p "$BIN_DIR"
TMP_BIN="$BIN_DIR/.traceknot-install-$$"
TMP_SUM="$BIN_DIR/.traceknot-install-$$.sha256"
trap 'rm -f "$TMP_BIN" "$TMP_SUM"' EXIT

echo "Installing traceknot $RELEASE..."
tk_log "install started ($RELEASE)"
download "$BINARY_URL" "$TMP_BIN"
download "$CHECKSUM_URL" "$TMP_SUM"

expected_sum="$(awk '{print $1}' "$TMP_SUM" 2>/dev/null)"
actual_sum="$(tk_sha256 "$TMP_BIN")"
if [[ -z "$expected_sum" || "$actual_sum" != "$expected_sum" ]]; then
	tk_log "ERROR: checksum verification failed for $BINARY_URL"
	echo "Error: the download from $BINARY_URL failed checksum verification (possible corruption)." >&2
	exit 1
fi
chmod +x "$TMP_BIN"
mv "$TMP_BIN" "$BIN_DIR/traceknot"
tk_log "binary installed at $BIN_DIR/traceknot"

if { [[ ! -t 0 ]] || [[ ! -t 1 ]]; } && [[ -e /dev/tty ]]; then
	"$BIN_DIR/traceknot" post-install < /dev/tty
else
	"$BIN_DIR/traceknot" post-install
fi
tk_log "post-install done"
tk_log "install complete"

echo ""
echo "Installation complete!"
echo "Dashboard available at: $DAEMON_URL/"
echo "Reconfigure anytime with: traceknot"
echo "Installation log: $INSTALL_LOG"

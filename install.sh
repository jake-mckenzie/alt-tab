#!/bin/sh
# Install the terminal version of Alt-Tab from a GitHub release.
set -eu

usage() {
    cat <<'HELP'
Usage: sh install.sh [--version vX.Y.Z] [--dir PATH] [--run]

Installs the latest stable terminal release to ~/.local/bin by default.
  --version TAG  Install a specific release
  --dir PATH     Choose the installation directory
  --run          Launch after installation (requires a terminal)
  --help         Show this help
HELP
}
fail() { printf 'alt-tab: %s\n' "$*" >&2; exit 1; }
version=latest
install_dir=${HOME:?HOME must be set}/.local/bin
run=no
while [ "$#" -gt 0 ]; do
    case "$1" in
        --version|--dir)
            [ "$#" -ge 2 ] || fail "$1 requires a value"
            case "$1" in --version) version=$2 ;; --dir) install_dir=$2 ;; esac
            shift 2 ;;
        --run) run=yes; shift ;;
        --help|-h) usage; exit 0 ;;
        *) fail "unknown option: $1" ;;
    esac
done
case "$version" in
    latest) ;;
    v[0-9]*) case "$version" in *[!a-zA-Z0-9._-]*) fail 'invalid version tag' ;; esac ;;
    *) fail 'version must be a v-prefixed release tag' ;;
esac
[ -n "$install_dir" ] || fail 'installation directory cannot be empty'
case "$(uname -s)" in
    Linux) os=linux ;;
    Darwin) os=darwin ;;
    *) fail 'supported systems: Linux and macOS' ;;
esac
case "$(uname -m)" in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) fail 'supported architectures: x86_64 and ARM64' ;;
esac
command -v curl >/dev/null 2>&1 || fail 'curl is required'
if command -v sha256sum >/dev/null 2>&1; then
    checksum=sha256sum
elif command -v shasum >/dev/null 2>&1; then
    checksum=shasum
else
    fail 'sha256sum or shasum is required'
fi
base=https://github.com/jake-mckenzie/alt-tab/releases
case "$version" in
    latest) base=$base/latest/download ;;
    *) base=$base/download/$version ;;
esac
asset=alt-tab_${os}_${arch}
tmp_dir=$(mktemp -d) || fail 'could not create temporary directory'
staged=
cleanup() {
    rm -rf "$tmp_dir"
    if [ -n "$staged" ]; then rm -f "$staged"; fi
}
trap cleanup EXIT
trap 'exit 1' HUP INT TERM
printf 'Downloading %s (%s)...\n' "$asset" "$version"
curl -fLsS --retry 3 --proto '=https' --tlsv1.2 "$base/$asset" -o "$tmp_dir/$asset"
curl -fLsS --retry 3 --proto '=https' --tlsv1.2 "$base/checksums.txt" -o "$tmp_dir/checksums.txt"
expected=$(awk -v name="$asset" '$2 == name {print $1}' "$tmp_dir/checksums.txt")
[ "${#expected}" -eq 64 ] || fail 'missing or invalid release checksum'
case "$expected" in *[!0-9a-f]*) fail 'invalid release checksum' ;; esac
if [ "$checksum" = sha256sum ]; then
    actual=$(sha256sum "$tmp_dir/$asset")
else
    actual=$(shasum -a 256 "$tmp_dir/$asset")
fi
actual=${actual%% *}
[ "$actual" = "$expected" ] || fail 'checksum mismatch; installation cancelled'
mkdir -p "$install_dir"
staged=$(mktemp "$install_dir/.alt-tab.XXXXXXXX")
cp "$tmp_dir/$asset" "$staged"
chmod 755 "$staged"
mv -f "$staged" "$install_dir/alt-tab"
staged=
printf 'Installed %s/alt-tab\n' "$install_dir"
case ":$PATH:" in
    *":$install_dir:"*) ;;
    *) printf 'Add this directory to your PATH: %s\n' "$install_dir" ;;
esac
if [ "$run" = yes ]; then
    [ -r /dev/tty ] && [ -w /dev/tty ] || fail 'launch requires an interactive terminal'
    "$install_dir/alt-tab" </dev/tty >/dev/tty 2>&1
else
    printf 'Run: %s/alt-tab\n' "$install_dir"
fi

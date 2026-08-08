#!/usr/bin/env sh
set -eu

repository="ly-wang19/AnyWork"
version="${ANYWORK_VERSION:-v0.3.0-alpha.1}"
install_dir="${ANYWORK_INSTALL_DIR:-${HOME}/.local/bin}"
download_base="${ANYWORK_DOWNLOAD_BASE:-}"
run_setup=false
agent="all"
language=""

fail() {
  printf '%s\n' "AnyWork installer: $*" >&2
  exit 1
}

usage() {
  printf '%s\n' "Usage: install.sh [--version vX.Y.Z] [--install-dir path] [--setup] [--agent id|all] [--lang en|zh-CN|ja]"
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      [ "$#" -ge 2 ] || fail "--version requires a value"
      version=$2
      shift 2
      ;;
    --install-dir)
      [ "$#" -ge 2 ] || fail "--install-dir requires a value"
      install_dir=$2
      shift 2
      ;;
    --setup)
      run_setup=true
      shift
      ;;
    --agent)
      [ "$#" -ge 2 ] || fail "--agent requires a value"
      agent=$2
      shift 2
      ;;
    --lang)
      [ "$#" -ge 2 ] || fail "--lang requires a value"
      language=$2
      shift 2
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      fail "unknown option: $1"
      ;;
  esac
done

if [ -n "${ANYWORK_OS:-}" ]; then
  target_os=$ANYWORK_OS
else
  case "$(uname -s)" in
    Darwin) target_os=darwin ;;
    Linux) target_os=linux ;;
    *) fail "unsupported operating system: $(uname -s)" ;;
  esac
fi

if [ -n "${ANYWORK_ARCH:-}" ]; then
  target_arch=$ANYWORK_ARCH
else
  case "$(uname -m)" in
    x86_64|amd64) target_arch=amd64 ;;
    arm64|aarch64) target_arch=arm64 ;;
    *) fail "unsupported architecture: $(uname -m)" ;;
  esac
fi

asset="anywork_${version}_${target_os}_${target_arch}.tar.gz"
if [ -z "$download_base" ]; then
  download_base="https://github.com/${repository}/releases/download/${version}"
fi

temporary=$(mktemp -d)
trap 'rm -rf "$temporary"' EXIT HUP INT TERM

download() {
  source_url=$1
  destination=$2
  case "$source_url" in
    file://*) cp "${source_url#file://}" "$destination" ;;
    *)
      if command -v curl >/dev/null 2>&1; then
        curl --fail --location --silent --show-error "$source_url" --output "$destination"
      elif command -v wget >/dev/null 2>&1; then
        wget -q "$source_url" -O "$destination"
      else
        fail "curl or wget is required"
      fi
      ;;
  esac
}

download "${download_base}/${asset}" "$temporary/$asset"
download "${download_base}/SHA256SUMS" "$temporary/SHA256SUMS"

expected=$(awk -v plain="$asset" -v dotted="./$asset" '$2 == plain || $2 == dotted { print $1; exit }' "$temporary/SHA256SUMS")
[ -n "$expected" ] || fail "checksum entry missing for $asset"
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$temporary/$asset" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  actual=$(shasum -a 256 "$temporary/$asset" | awk '{print $1}')
else
  fail "sha256sum or shasum is required"
fi
[ "$expected" = "$actual" ] || fail "checksum mismatch for $asset"

tar -xzf "$temporary/$asset" -C "$temporary"
[ -f "$temporary/anywork" ] || fail "archive does not contain anywork"
mkdir -p "$install_dir"
if command -v install >/dev/null 2>&1; then
  install -m 0755 "$temporary/anywork" "$install_dir/anywork"
else
  cp "$temporary/anywork" "$install_dir/anywork"
  chmod 0755 "$install_dir/anywork"
fi

"$install_dir/anywork" --version
printf '%s\n' "Installed AnyWork to $install_dir/anywork"

if [ "$run_setup" = true ]; then
  if [ -n "$language" ]; then
    "$install_dir/anywork" setup --agent "$agent" --scope user --lang "$language"
  else
    "$install_dir/anywork" setup --agent "$agent" --scope user
  fi
else
  printf '%s\n' "Activate all capabilities: $install_dir/anywork setup --agent all"
fi

case ":${PATH}:" in
  *":${install_dir}:"*) ;;
  *) printf '%s\n' "Add $install_dir to PATH to run 'anywork' from any shell." ;;
esac

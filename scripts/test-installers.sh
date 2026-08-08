#!/usr/bin/env sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
temporary=$(mktemp -d)
trap 'rm -rf "$temporary"' EXIT HUP INT TERM

case "$(uname -s)" in
  Darwin) target_os=darwin ;;
  Linux) target_os=linux ;;
  *) printf '%s\n' "Skipping Unix installer smoke test on $(uname -s)"; exit 0 ;;
esac
case "$(uname -m)" in
  x86_64|amd64) target_arch=amd64 ;;
  arm64|aarch64) target_arch=arm64 ;;
  *) printf '%s\n' "Skipping unsupported architecture $(uname -m)"; exit 0 ;;
esac

asset="anywork_test_${target_os}_${target_arch}.tar.gz"
payload="$temporary/payload"
release="$temporary/release"
install_dir="$temporary/install"
isolated_home="$temporary/home"
mkdir -p "$payload" "$release" "$isolated_home"
go build -C "$root" -trimpath -ldflags "-X main.nativeVersion=test" -o "$payload/anywork" .
tar -C "$payload" -czf "$release/$asset" anywork
if command -v sha256sum >/dev/null 2>&1; then
  (cd "$release" && sha256sum "$asset" > SHA256SUMS)
else
  (cd "$release" && shasum -a 256 "$asset" > SHA256SUMS)
fi

HOME="$isolated_home" ANYWORK_DOWNLOAD_BASE="file://$release" ANYWORK_OS="$target_os" ANYWORK_ARCH="$target_arch" \
  sh "$root/install.sh" --version test --install-dir "$install_dir" --setup --agent codex --lang en
"$install_dir/anywork" doctor --lang en
test -f "$isolated_home/.agents/skills/orchestrate-work/SKILL.md"
printf '%s\n' "Unix installer smoke test passed."

#!/usr/bin/env bash
set -euo pipefail

repository_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
output_directory=${1:-"$repository_root/dist"}
release_version=${VERSION:-dev}

if [[ "$output_directory" != /* ]]; then
  output_directory="$repository_root/$output_directory"
fi

if [[ -e "$output_directory" ]]; then
  echo "Output path already exists: $output_directory" >&2
  exit 1
fi

mkdir -p "$output_directory"
build_directory=$(mktemp -d)
trap 'rm -rf "$build_directory"' EXIT

targets=(
  "darwin amd64"
  "darwin arm64"
  "linux amd64"
  "linux arm64"
  "windows amd64"
  "windows arm64"
)

for target in "${targets[@]}"; do
  read -r target_os target_arch <<<"$target"
  archive_base="anywork_${release_version}_${target_os}_${target_arch}"
  target_directory="$build_directory/$archive_base"
  mkdir -p "$target_directory"
  binary_name=anywork
  if [[ "$target_os" == windows ]]; then
    binary_name=anywork.exe
  fi
  CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" \
    go build -C "$repository_root" -trimpath \
      -ldflags "-s -w -X main.nativeVersion=${release_version#v}" \
      -o "$target_directory/$binary_name" .
  cp "$repository_root/LICENSE" "$repository_root/NOTICE" "$repository_root/THIRD_PARTY.yml" "$target_directory/"
  if [[ "$target_os" == windows ]]; then
    (cd "$target_directory" && zip -q "$output_directory/$archive_base.zip" "$binary_name" LICENSE NOTICE THIRD_PARTY.yml)
  else
    tar -C "$target_directory" -czf "$output_directory/$archive_base.tar.gz" "$binary_name" LICENSE NOTICE THIRD_PARTY.yml
  fi
done

(cd "$output_directory" && shasum -a 256 ./* > SHA256SUMS)
echo "Built six release archives in $output_directory"

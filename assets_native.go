package main

import "embed"

// officialContent is the immutable content compiled into an AnyWork release.
// Install and update never execute code from this filesystem.
//
//go:embed registry/*.json evidence/*.json evals/*.json adapters skills
var officialContent embed.FS

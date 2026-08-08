//go:build !windows

package main

import (
	"errors"
	"os"
	"syscall"
)

func atomicReplace(source, destination string) error {
	return os.Rename(source, destination)
}

func syncDirectory(directory string) error {
	file, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := file.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
		return err
	}
	return nil
}

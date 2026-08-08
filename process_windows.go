//go:build windows

package main

import "syscall"

func processAlive(pid int) bool {
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err == syscall.ERROR_ACCESS_DENIED {
		return true
	}
	if err != nil {
		return false
	}
	_ = syscall.CloseHandle(handle)
	return true
}

package util

import (
	"runtime"
	"strings"
)

func IsWindows() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

func IsMac() bool {
	return strings.Contains(runtime.GOOS, "darwin")
}

func IsWinOrMac() bool {
	return IsWindows() || IsMac()
}

func IsLinux() bool {
	return strings.Contains(runtime.GOOS, "linux")
}

package xfmt

import "fmt"

func UnitCount(f string, b int64, unit int64, suffix string) string {

	if b < unit {
		return fmt.Sprintf("%d%s", b, suffix)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf(f+" %c%s",
		float64(b)/float64(div), "KMGTPE"[exp], suffix)
}

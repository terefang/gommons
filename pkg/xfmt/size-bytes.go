package xfmt

func ByteCountSI(b int64) string {
	return UnitCount("%0.2f", b, 1024, "B")
}

func ByteFmtCountSI(f string, b int64) string {
	return UnitCount(f, b, 1000, "B")
}

func ByteCountIEC(b int64) string {
	return UnitCount("%0.2f", b, 1024, "iB")
}

func ByteFmtCountIEC(f string, b int64) string {
	return UnitCount(f, b, 1024, "iB")
}

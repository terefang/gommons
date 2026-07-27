package xconvert

func AsArray[V int64 | float64 | string | bool | rune | int | byte | any](args ...V) []V {
	return args
}

func ToArray[V int64 | float64 | string | bool | rune | int | byte | any](args ...V) []V {
	_ret := make([]V, len(args))
	for i, arg := range args {
		_ret[i] = arg
	}
	return _ret
}

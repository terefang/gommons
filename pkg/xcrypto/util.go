package xcrypto

func CalcBytesFromBits(bits int) int {
	_len := (bits / 8)
	_mod := bits % 8
	if _mod != 0 {
		_len++
	}
	return _len
}

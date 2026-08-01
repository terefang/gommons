package xbytes

func Fill(_buf []byte, _i int) []byte {
    for i := 0; i < len(_buf); i++ {
        _buf[i] = byte(_i)
    }
    return _buf
}

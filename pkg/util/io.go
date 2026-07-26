package util

import (
	"encoding/binary"
	"io"
)

func ReadBeInt32(r io.Reader) (val int32) {
	_ = binary.Read(r, binary.BigEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadBeInt16(r io.Reader) (val int16) {
	_ = binary.Read(r, binary.BigEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadBeInt64(r io.Reader) (val int64) {
	_ = binary.Read(r, binary.BigEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadLeInt32(r io.Reader) (val int32) {
	_ = binary.Read(r, binary.LittleEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadLeInt16(r io.Reader) (val int16) {
	_ = binary.Read(r, binary.LittleEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadLeInt64(r io.Reader) (val int64) {
	_ = binary.Read(r, binary.LittleEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadByte(r io.Reader) (val byte) {
	_ = binary.Read(r, binary.BigEndian, &val)
	//if err != nil {
	//		f.err = err
	//	}
	return
}

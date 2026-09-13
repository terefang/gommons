/*
 * Copyright (C) 2024 Andrew Ayer
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * Except as contained in this notice, the name(s) of the above copyright
 * holders shall not be used in advertising or otherwise to promote the
 * sale, use or other dealings in this Software without prior written
 * authorization
 */

package xar

import (
	"fmt"
	"io"
)

type Reader struct {
	r         io.Reader
	fileBytes int64
	padded    bool
}

func NewReader(r io.Reader) (*Reader, error) {
	var header [8]byte
	if _, err := io.ReadFull(r, header[:]); err == io.EOF {
		return nil, io.EOF
	} else if err != nil {
		return nil, fmt.Errorf("error reading archive header: %w", err)
	}
	if string(header[:]) != "!<arch>\n" {
		return nil, fmt.Errorf("does not look like an ar archive")
	}
	return &Reader{r: r}, nil
}

func (r *Reader) Next() (*Header, error) {
	if r.fileBytes > 0 {
		n, err := io.CopyN(io.Discard, r.r, r.fileBytes)
		r.fileBytes -= n
		if err != nil {
			return nil, err
		}
	}
	if r.padded {
		var buf [1]byte
		if _, err := r.r.Read(buf[:]); err != nil {
			return nil, err
		}
		r.padded = false
		if buf[0] != '\n' {
			return nil, fmt.Errorf("archive contains incorrect padding")
		}
	}

	header, err := ReadHeader(r.r)
	if err != nil {
		return nil, err
	}
	r.fileBytes = header.Size
	r.padded = (header.Size%2 == 1)
	return header, nil
}

func (r *Reader) Read(buf []byte) (int, error) {
	if r.fileBytes == 0 {
		return 0, io.EOF
	} else if int64(len(buf)) > r.fileBytes {
		buf = buf[:r.fileBytes]
	}
	n, err := r.r.Read(buf)
	r.fileBytes -= int64(n)
	return n, err
}

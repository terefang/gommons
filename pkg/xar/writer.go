/*
 * Copyright (C) 2025 Andrew Ayer
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

type Writer struct {
	w         io.Writer
	fileBytes int64
	padded    bool
}

func NewWriter(w io.Writer) (*Writer, error) {
	if _, err := io.WriteString(w, "!<arch>\n"); err != nil {
		return nil, fmt.Errorf("error writing archive header: %w", err)
	}
	return &Writer{w: w}, nil
}

func (w *Writer) WriteHeader(header *Header) error {
	if err := w.Flush(); err != nil {
		return err
	}
	if err := header.Write(w.w); err != nil {
		return err
	}
	w.fileBytes = header.Size
	w.padded = (header.Size%2 == 1)
	return nil
}

func (w *Writer) Write(buf []byte) (int, error) {
	if int64(len(buf)) > w.fileBytes {
		return 0, fmt.Errorf("too many bytes written to current file")
	}
	n, err := w.w.Write(buf)
	w.fileBytes -= int64(n)
	return n, err
}

func (w *Writer) Flush() error {
	if w.fileBytes > 0 {
		return fmt.Errorf("not enough bytes written to current file")
	}
	if w.padded {
		if _, err := w.w.Write([]byte{'\n'}); err != nil {
			return fmt.Errorf("error writing file padding: %w", err)
		}
		w.padded = false
	}
	return nil
}

func (w *Writer) Close() error {
	return w.Flush()
}

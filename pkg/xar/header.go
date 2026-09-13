/*
 * Copyright (C) 2024 Andrew Ayer
 * Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
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
	"strconv"
	"strings"
)

type Header struct {
	Name      string
	Timestamp int64
	OwnerID   int64
	GroupID   int64
	FileMode  string
	Size      int64
}

func ReadHeader(r io.Reader) (*Header, error) {
	line := make([]byte, 60)
	if _, err := io.ReadFull(r, line); err == io.EOF {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("error reading ar file header: %w", err)
	}

	// +-------------------------------------------------------
	// | Offset  Length  Name                         Format
	// +-------------------------------------------------------
	// | 0       16      File name                    ASCII
	// | 16      12      File modification timestamp  Decimal
	// | 28      6       Owner ID                     Decimal
	// | 34      6       Group ID                     Decimal
	// | 40      8       File mode                    Octal
	// | 48      10      File size in bytes           Decimal
	// | 58      2       File magic                   0x60 0x0A
	if line[58] != 0x60 || line[59] != 0x0A {
		return nil, fmt.Errorf("error parsing ar file header: malformed line ending")
	}

	hdr := &Header{
		Name:     strings.TrimSuffix(strings.TrimSpace(string(line[0:16])), "/"),
		FileMode: strings.TrimSpace(string(line[40:48])),
	}

	for target, value := range map[*int64][]byte{
		&hdr.Timestamp: line[16:28],
		&hdr.OwnerID:   line[28:34],
		&hdr.GroupID:   line[34:40],
		&hdr.Size:      line[48:58],
	} {
		intValue, err := strconv.ParseInt(strings.TrimSpace(string(value)), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing decimal value in ar file header: %w", err)
		}
		*target = intValue
	}

	return hdr, nil
}

func (hdr *Header) Write(w io.Writer) error {
	if len(hdr.Name) > 16 {
		return fmt.Errorf("file name is too long")
	}
	if hdr.Timestamp < 0 || hdr.Timestamp > 999999999999 {
		return fmt.Errorf("timestamp is out of range")
	}
	if hdr.OwnerID < 0 || hdr.OwnerID > 999999 {
		return fmt.Errorf("owner ID is out of range")
	}
	if hdr.GroupID < 0 || hdr.GroupID > 999999 {
		return fmt.Errorf("group ID is out of range")
	}
	if len(hdr.FileMode) > 8 {
		return fmt.Errorf("file mode is too long")
	}
	if hdr.Size < 0 || hdr.Size > 9999999999 {
		return fmt.Errorf("file size is out of range")
	}
	_, err := fmt.Fprintf(w, "%-16s%-12d%-6d%-6d%-8s%-10d`\n", hdr.Name, hdr.Timestamp, hdr.OwnerID, hdr.GroupID, hdr.FileMode, hdr.Size)
	return err
}

package xfs

import (
    "io/fs"
    "slices"
)

// Copyright ©2024 The overlayfs Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package overlyafs provides a way to overlay multiple fs.FS.

// OverlayFS overlays multiple fs.FS.
type OverlayFS struct {
    stack []fs.FS
}

var (
    _ fs.FS    = (*OverlayFS)(nil)
    _ fs.SubFS = (*OverlayFS)(nil)
)

// OverlayFrom creates a new overlay fs from the provided stack of OverlayFS.
// OverlayFS contents are considered in the reverse order of the provided slice.
// ie: the last stacked layer wins.
func OverlayFrom(stack ...fs.FS) OverlayFS {
    o := OverlayFS{
        stack: make([]fs.FS, len(stack)),
    }
    copy(o.stack, stack)
    slices.Reverse(o.stack)
    return o
}

// Open opens the named file.
//
// Open tries each overlaid filesystem in-turn, in the
// reverse order they were overlaid.
// ie: the last stacked layer wins.
func (fsys OverlayFS) Open(name string) (fs.File, error) {
    for _, v := range fsys.stack {
        f, err := v.Open(name)
        if err == nil {
            return f, nil
        }
    }

    return nil, fs.ErrNotExist
}

// Sub returns an OverlayFS corresponding to the subtree rooted at dir.
//
// Sub tries each overlaid filesystem in-turn, in the
// reverse order they were overlaid.
// ie: the last stacked layer wins.
func (fsys OverlayFS) Sub(dir string) (fs.FS, error) {
    for _, v := range fsys.stack {
        sub, err := fs.Sub(v, dir)
        if err == nil {
            return sub, nil
        }
    }

    return nil, fs.ErrNotExist
}

func (fsys *OverlayFS) AddFs(f fs.FS) {
    fsys.stack = append([]fs.FS{f}, fsys.stack...)
}
